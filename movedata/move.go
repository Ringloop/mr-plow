package movedata

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log"

	"dariobalinzo.com/elastic/v2/config"
	"dariobalinzo.com/elastic/v2/elastic"

	"github.com/elastic/go-elasticsearch/v7/esutil"
	_ "github.com/lib/pq"
)

func MoveData(db *sql.DB, c config.QueryModel) error {
	lastDate, err := elastic.FindLastUpdateOrEpochDate(c.Index)
	if err != nil {
		return nil
	}
	log.Print("found last date ", lastDate)

	elasticBulk, err := elastic.GetBulkIndexer(c.Index)
	if err != nil {
		return nil
	}
	defer elasticBulk.Close(context.Background())

	log.Print("going to execute query ", c.Query)
	rows, err := db.Query(c.Query, lastDate)
	if err != nil {
		return err
	}
	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	columnsMap := make(map[string]string)
	for _, colConfig := range c.Fields {
		columnsMap[colConfig.Name] = colConfig.Type
	}

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {

			columnPointers[i] = &columns[i]

			//HERE:
			// per ogni ColumnPointers[i] devo fare la validazione del tipo. Se il target è intero, chiama la funzione giusta

		}

		if err := rows.Scan(columnPointers...); err != nil {
			return err
		}

		document := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})

			//Qua va fatta la stessa cosa di sopra agendo su VAL

			document[colName] = *val //qui ci andrà messo il tipo convertito
		}

		for _, jsonfield := range c.JSONFields {
			//TODO: check if jsonField is a json itself

			//TODO: Build the json structure
			var jsonData map[string]interface{}

			//TODO: Throw an error if the fieldName cannot be found
			byteData := document[jsonfield.FieldName].([]byte)
			//TODO consider also the field types in parsing
			json.Unmarshal(byteData, &jsonData)
			document[jsonfield.FieldName] = jsonData
		}

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
