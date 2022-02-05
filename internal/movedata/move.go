package movedata

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Ringloop/mr-plow/internal/casting"
	"github.com/Ringloop/mr-plow/internal/config"
	"github.com/Ringloop/mr-plow/internal/elastic"

	"github.com/elastic/go-elasticsearch/v7/esutil"
	_ "github.com/lib/pq"
)

type Mover struct {
	lastDate     *time.Time
	db           *sql.DB
	globalConfig *config.ImportConfig
	queryConf    *config.QueryModel
	columnsMap   map[string]string
	jsonColsMap  map[string]map[string]string
	canExec      chan bool
}

func New(db *sql.DB, globalConfig *config.ImportConfig, queryConf *config.QueryModel) *Mover {
	columnsMap, jsonColsMap := getColumnsConfiguration(queryConf)

	mover := &Mover{
		db:           db,
		globalConfig: globalConfig,
		queryConf:    queryConf,
		columnsMap:   columnsMap,
		jsonColsMap:  jsonColsMap,
		canExec:      make(chan bool, 1)}
	mover.canExec <- true
	return mover
}

func getColumnsConfiguration(queryConf *config.QueryModel) (map[string]string, map[string]map[string]string) {
	//create the map for the native fields of the query
	columnsMap := make(map[string]string)
	for _, colConfig := range queryConf.Fields {
		columnsMap[colConfig.Name] = colConfig.Type
	}

	//create a nested map for the JSON fields (any JSON contains a set of fields)
	var jsonColsMap = make(map[string]map[string]string)
	for _, jsonColConfig := range queryConf.JSONFields {
		for _, colConfig := range jsonColConfig.Fields {
			if jsonColsMap[jsonColConfig.FieldName] == nil {
				jsonColsMap[jsonColConfig.FieldName] = make(map[string]string)
			}
			jsonColsMap[jsonColConfig.FieldName][colConfig.Name] = colConfig.Type
		}
	}
	return columnsMap, jsonColsMap
}

func (mover *Mover) MoveData() error {
	select {
	case <-mover.canExec:
	default:
		log.Printf("Skipping execution of %s, since the previous one is still executing", mover.queryConf.Query)
		return nil
	}

	defer func() {
		mover.canExec <- true
	}()

	repo, err := elastic.NewClient(mover.globalConfig)
	if err != nil {
		return err
	}

	lastDate, err := mover.getLastDate(repo)
	if err != nil {
		return err
	}

	elasticBulk, err := repo.GetBulkIndexer(mover.queryConf.Index)
	if err != nil {
		return err
	}
	defer elasticBulk.Close(context.Background())

	log.Printf("going to execute query %s with param %s", mover.queryConf.Query, lastDate)
	rows, err := mover.db.Query(mover.queryConf.Query, lastDate)
	if err != nil {
		return err
	}
	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	for rows.Next() {
		columns := make([](interface{}), len(cols))
		for i := range columns {
			columns[i] = &columns[i]
		}

		if err := rows.Scan(columns...); err != nil {
			return err
		}
		converter := casting.NewConverter(mover.columnsMap)
		convertedArrayOfData := converter.CastArrayOfData(cols, columns)
		document := make(map[string]interface{})
		for i, colName := range cols {
			document[colName] = convertedArrayOfData[i]
		}
		for _, jsonfield := range mover.queryConf.JSONFields {
			var jsonData map[string]interface{}

			data, ok := document[jsonfield.FieldName]
			if !ok {
				return fmt.Errorf("error getting ....: %s", err)
			}
			byteData := data.([]byte)

			err := json.Unmarshal(byteData, &jsonData)
			if err != nil {
				return err
			}
			for _, field := range jsonfield.Fields {
				jsonConverter := casting.NewConverter(mover.jsonColsMap[jsonfield.FieldName])
				jsonData[field.Name] = jsonConverter.CastSingleElement(field.Name, jsonData[field.Name])
			}
			document[jsonfield.FieldName] = jsonData
		}

		documentToSend, err := json.Marshal(document)
		if err != nil {
			return err
		}

		bulkItem := esutil.BulkIndexerItem{
			Action: "index",
			Body:   bytes.NewBuffer(documentToSend),
			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
				if err != nil {
					log.Printf("ERROR: %s", err)
				} else {
					log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
				}
			},
		}
		addDocumentId(mover.queryConf, document, &bulkItem)
		err = mover.updateLastUpdate(mover.queryConf, document)
		if err != nil {
			return err
		}

		err = elasticBulk.Add(context.Background(), bulkItem)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mover *Mover) getLastDate(repo *elastic.Repository) (*time.Time, error) {
	if mover.lastDate != nil {
		return mover.lastDate, nil
	}

	lastDate, err := repo.FindLastUpdateOrEpochDate(mover.queryConf.Index, mover.queryConf.UpdateDate)
	if err != nil {
		return nil, err
	}
	log.Print("found last date ", lastDate)
	return lastDate, nil
}

func (mover *Mover) updateLastUpdate(conf *config.QueryModel, document map[string]interface{}) error {
	date, ok := document[conf.UpdateDate]
	if !ok {
		return fmt.Errorf("cannot find %s in results set of %s", conf.UpdateDate, conf.Query)
	}
	dateParsed, ok := date.(time.Time)
	if !ok {
		return fmt.Errorf("cannot cast to date %s in results set of %s", date, conf.Query)
	}
	mover.lastDate = &dateParsed
	return nil
}

func addDocumentId(queryConf *config.QueryModel, document map[string]interface{}, bulkItem *esutil.BulkIndexerItem) {
	if queryConf.Id != "" {
		id, present := document[queryConf.Id]
		if present {
			idAsString := fmt.Sprint(id)
			bulkItem.DocumentID = idAsString
		}
	}
}
