package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"dariobalinzo.com/elastic/v2/config"
	"dariobalinzo.com/elastic/v2/elastic"

	"github.com/elastic/go-elasticsearch/v7/esutil"
	_ "github.com/lib/pq"
)

func main() {
	ymlConfReader := config.Reader{FileName: "config.yml"} //TODO parse commandline yml path, now assuming is in current dir
	conf, err := config.ParseConfiguration(&ymlConfReader)
	if err != nil {
		log.Fatal("Cannot parse config file", err)
	}
	ConnectAndStart(conf)
}

func ConnectAndStart(conf *config.ImportConfig) {
	db, err := sql.Open("postgres", conf.SqlConfig)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()

	fmt.Println("Connected to postgres", err)

	lastDate, err := elastic.FindLastUpdate(conf.Index)
	if err != nil {
		log.Fatal("error fetching last execution query")
	}

	if lastDate == nil {
		var defaultDate time.Time
		defaultDate, err = time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
		lastDate = &defaultDate
	}
	if err != nil {
		log.Fatal("error in default last update date setting", err)
	}

	err = moveData(db, conf.Query, lastDate, conf.Index)
	if err != nil {
		log.Fatal("error execurting query", err)
	}

}

func moveData(db *sql.DB, query string, last_update *time.Time, index string) error {
	elasticBulk, err := elastic.GetBulkIndexer(index)
	if err != nil {
		return nil
	}
	defer elasticBulk.Close(context.Background())
	rows, err := db.Query(query, last_update)
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

		//parsing data json
		var jsonData map[string]interface{}
		byteData := document["data"].([]byte)
		json.Unmarshal(byteData, &jsonData)
		document["data"] = jsonData

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
