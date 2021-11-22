package movedata

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log"

	"dariobalinzo.com/elastic/v2/elastic"

	"github.com/elastic/go-elasticsearch/v7/esutil"
	_ "github.com/lib/pq"
)

func MoveData(db *sql.DB, query, index string) error {
	lastDate, err := elastic.FindLastUpdateOrEpochDate(index)
	if err != nil {
		return nil
	}
	log.Print("found last date ", lastDate)

	elasticBulk, err := elastic.GetBulkIndexer(index)
	if err != nil {
		return nil
	}
	defer elasticBulk.Close(context.Background())

	log.Print("going to execute query ", query)
	rows, err := db.Query(query, lastDate)
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

		//TODO better parsing data json
		//var jsonData map[string]interface{}
		//byteData := document["data"].([]byte)
		//json.Unmarshal(byteData, &jsonData)
		//document["data"] = jsonData

		documentToSend, err := json.Marshal(document)
		if err != nil {
			return err
		}

		err = elasticBulk.Add(context.Background(), esutil.BulkIndexerItem{
			Action: "index",
			Body:   bytes.NewBuffer(documentToSend),
			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
				if err != nil {
					log.Printf("ERROR: %s", err)
				} else {
					log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
				}
			},
		})
		if err != nil {
			return err
		}
	}
	return nil

}
