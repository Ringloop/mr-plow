package movedata

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"dariobalinzo.com/elastic/v2/config"
	"dariobalinzo.com/elastic/v2/elastic"

	"github.com/elastic/go-elasticsearch/v7/esutil"
	_ "github.com/lib/pq"
)

func MoveData(db *sql.DB, globalConfig *config.ImportConfig, queryConf config.QueryModel) error {
	repo, err := elastic.NewClient(globalConfig)
	if err != nil {
		return err
	}

	lastDate, err := repo.FindLastUpdateOrEpochDate(queryConf.Index, queryConf.UpdateDate)
	if err != nil {
		return err
	}
	log.Print("found last date ", lastDate)

	elasticBulk, err := repo.GetBulkIndexer(queryConf.Index)
	if err != nil {
		return err
	}
	defer elasticBulk.Close(context.Background())

	log.Print("going to execute query ", queryConf.Query)
	rows, err := db.Query(queryConf.Query, lastDate)
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

		for _, jsonfield := range queryConf.JSONFields {
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
		addDocumentId(&queryConf, document, &bulkItem)

		err = elasticBulk.Add(context.Background(), bulkItem)
		if err != nil {
			return err
		}
	}
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
