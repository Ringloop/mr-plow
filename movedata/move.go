package movedata

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Ringloop/Mr-Plow/config"
	"github.com/Ringloop/Mr-Plow/elastic"

	"github.com/elastic/go-elasticsearch/v7/esutil"
	_ "github.com/lib/pq"
)

type Mover struct {
	lastDate     *time.Time
	db           *sql.DB
	globalConfig *config.ImportConfig
	queryConf    *config.QueryModel
	canExec      chan bool
}

func New(db *sql.DB, globalConfig *config.ImportConfig, queryConf *config.QueryModel) *Mover {
	mover := &Mover{
		db:           db,
		globalConfig: globalConfig,
		queryConf:    queryConf,
		canExec:      make(chan bool, 1)}
	mover.canExec <- true
	return mover
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

	log.Print("going to execute query ", mover.queryConf.Query)
	rows, err := mover.db.Query(mover.queryConf.Query, lastDate)
	if err != nil {
		return err
	}
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return err
		}

		document := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			document[colName] = *val
		}

		for _, jsonfield := range mover.queryConf.JSONFields {
			var jsonData map[string]interface{}
			byteData := document[jsonfield.FieldName].([]byte)
			//TODO consider also the field types in parsing
			err := json.Unmarshal(byteData, &jsonData)
			if err != nil {
				return err
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
		return fmt.Errorf("cannot found %s in results set of %s", conf.UpdateDate, conf.Query)
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
