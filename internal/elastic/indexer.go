package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/Ringloop/mr-plow/internal/config"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
)

type Repository struct {
	es            *elasticsearch.Client
	numWorkers    int
	flushBytes    int
	flushInterval time.Duration
}

func NewDefaultClient() (*Repository, error) {
	if es, err := elasticsearch.NewDefaultClient(); err != nil {
		return &Repository{}, err
	} else {
		return &Repository{
			es:            es,
			numWorkers:    1,
			flushBytes:    100000,
			flushInterval: 30 * time.Second}, nil
	}
}

func NewClient(config *config.ImportConfig) (*Repository, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			config.Elastic.Url,
		},
	}

	if config.Elastic.User != "" {
		cfg.Username = config.Elastic.User
		cfg.Password = config.Elastic.Password
	}

	if config.Elastic.CaCertPath != "" {
		cert, err := ioutil.ReadFile(config.Elastic.CaCertPath)
		if err != nil {
			return nil, err
		}
		cfg.CACert = cert
	}

	if es, err := elasticsearch.NewClient(cfg); err != nil {
		return &Repository{}, err
	} else {
		return &Repository{
			es:            es,
			numWorkers:    config.Elastic.NumWorker,
			flushBytes:    100000,
			flushInterval: 30 * time.Second}, nil
	}
}

func (repo *Repository) GetBulkIndexer(index string) (esutil.BulkIndexer, error) {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         index,
		Client:        repo.es,
		NumWorkers:    repo.numWorkers,
		FlushBytes:    repo.flushBytes,
		FlushInterval: repo.flushInterval,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting bulkIndexer: %s", err)
	}
	return bi, nil
}

func (repo *Repository) FindLastUpdateOrEpochDate(index, sortingField string) (*time.Time, error) {
	lastDate, err := repo.FindLastUpdate(index, sortingField)
	if err != nil {
		return nil, err
	}

	if lastDate == nil {
		log.Printf("cannot find old values for %s", sortingField)
		var defaultDate time.Time
		defaultDate, err = time.Parse(time.RFC3339, "1970-01-01T00:00:00+00:00")
		lastDate = &defaultDate
	}

	return lastDate, err
}

func (repo *Repository) FindLastUpdate(index, sortingField string) (*time.Time, error) {
	err := repo.Refresh(index)
	if err != nil {
		return nil, err
	}
	var query = `
	{
		"sort": [
		  {
			"$order": {
			  "order": "desc"
			}
		  }
		],
		"size": 1,
		"_source": [
		  "$order"
		  ]
	}
	`
	query = replaceOrderByField(query, sortingField)

	var mapResp map[string]interface{}

	res, err := repo.es.Search(
		repo.es.Search.WithContext(context.Background()),
		repo.es.Search.WithIndex(index),
		repo.es.Search.WithBody(strings.NewReader(query)),
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
		last_update := data[sortingField].(string)
		log.Println("found old data:", last_update)

		parsed_last_date, err := time.Parse(time.RFC3339, last_update)
		return &parsed_last_date, err

	}

}

func (repo *Repository) FindIndexContent(index, sortingField string) (*io.ReadCloser, error) {
	err := repo.Refresh(index)
	if err != nil {
		return nil, err
	}
	var query = `
	{
		"sort": [
		  {
			"$order": {
			  "order": "desc"
			}
		  }
		],
		"size": 1000
	}
	`
	query = replaceOrderByField(query, sortingField)

	res, err := repo.es.Search(
		repo.es.Search.WithContext(context.Background()),
		repo.es.Search.WithIndex(index),
		repo.es.Search.WithBody(strings.NewReader(query)),
	)

	if err != nil {
		return nil, err
	} else {
		return &res.Body, nil
	}
}

func replaceOrderByField(query, sortingField string) string {
	query = strings.Replace(query, "$order", sortingField, 2)
	return query
}

func (repo *Repository) Refresh(index string) error {
	r := esapi.IndicesRefreshRequest{
		Index: []string{index},
	}
	_, err := r.Do(context.Background(), repo.es)
	return err
}

func (repo *Repository) Delete(index string) error {
	_, err := repo.es.Indices.Delete([]string{index})
	return err
}
