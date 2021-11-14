package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
)

var es *elasticsearch.Client

func init() {
	var err error
	es, err = elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
}

func Index(index string, document map[string]interface{}) error {
	jsonBytes, err := json.Marshal(document)
	if err != nil {
		return nil
	}
	req := esapi.IndexRequest{
		Index: index,
		//DocumentID: strconv.Itoa(i + 1),
		Body: bytes.NewReader(jsonBytes),
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error indexing", res.Status())
	}
	return nil
}

func GetBulkIndexer(index string) (esutil.BulkIndexer, error) {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         index,
		Client:        es,
		NumWorkers:    10,        //todo config
		FlushBytes:    100000,    //todo config
		FlushInterval: time.Hour, // Disable automatic flushing
	})
	if err != nil {
		return nil, fmt.Errorf("error getting bulkIndexer: %s", err)
	}
	return bi, nil
}

func FindLastUpdate(index string) (*time.Time, error) {
	var query = `
	{
		"sort": [
		  {
			"last_update": {
			  "order": "desc"
			}
		  }
		],
		"size": 1,
		"_source": [
		  "last_update"
		  ]
	}
	`

	var mapResp map[string]interface{}

	// Pass the JSON query to the Golang client's Search() method
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(index),
		es.Search.WithBody(strings.NewReader(query)),
	)

	if err != nil {
		return nil, err
	} else {
		defer res.Body.Close()
		err := json.NewDecoder(res.Body).Decode(&mapResp)
		if err != nil {
			return nil, err
		}

		if mapResp["hits"] == nil {
			return nil, nil //index non existing
		}

		hits := mapResp["hits"].(map[string]interface{})
		hitsList := hits["hits"].([]interface{})
		if len(hitsList) == 0 {
			return nil, nil //no data in the index
		}

		data := hitsList[0].(map[string]interface{})["_source"].(map[string]interface{})
		last_update := data["last_update"].(string)
		fmt.Println("data:", last_update)

		parsed_last_date, err := time.Parse(time.RFC3339, last_update)
		return &parsed_last_date, err

	}

}
