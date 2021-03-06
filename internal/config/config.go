package config

import (
	"errors"
	"fmt"
	"log"

	"gopkg.in/yaml.v2"
)

type JSONField struct {
	FieldName string  `yaml:"fieldName"`
	Fields    []Field `yaml:"fields"`
}

type Field struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type QueryModel struct {
	Index      string      `yaml:"index"`
	Query      string      `yaml:"query"`
	UpdateDate string      `yaml:"updateDate"`
	Fields     []Field     `yaml:"fields"`
	JSONFields []JSONField `yaml:"JSONFields"`
	Id         string      `yaml:"id"`
}

type ElasticConfig struct {
	Url        string `yaml:"url"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	CaCertPath string `yaml:"caCertPath"`
	NumWorker  int    `yaml:"numWorker"`
}

type ImportConfig struct {
	PollingSeconds int           `yaml:"pollingSeconds"`
	Database       string        `yaml:"database"`
	Queries        []QueryModel  `yaml:"queries"`
	Elastic        ElasticConfig `yaml:"elastic"`
}

func ParseConfiguration(reader IReader) (*ImportConfig, error) {
	yamlFile, err := reader.ReadConfig()
	if err != nil {
		return nil, err
	}

	importConfiguration := &ImportConfig{}
	err = yaml.Unmarshal(yamlFile, importConfiguration)
	if err != nil {
		return &ImportConfig{}, err
	}

	if importConfiguration.PollingSeconds == 0 {
		return nil, errors.New("missing polling seconds url (pollingSeconds)")
	}

	if importConfiguration.Database == "" {
		return nil, errors.New("missing database url (database)")
	}

	if len(importConfiguration.Queries) == 0 {
		return nil, errors.New("missing query configuration (queries)")
	}

	err = validateQueriesConfig(importConfiguration)

	if importConfiguration.Elastic.Url == "" {
		return nil, errors.New("missing elastic url (elastic.url)")
	}

	if importConfiguration.Elastic.NumWorker < 1 {
		importConfiguration.Elastic.NumWorker = 1
		log.Println("using default worker = 1 for each elasticsearch indexer")
	}

	if err != nil {
		return nil, err
	} else {
		return importConfiguration, nil
	}
}

func validateQueriesConfig(importConfiguration *ImportConfig) error {
	for i, query := range importConfiguration.Queries {
		if query.Index == "" {
			return fmt.Errorf("missing output index for query %d  (queries.index)", i)
		}

		if query.Query == "" {
			return fmt.Errorf("missing query for query %d  (queries.query)", i)
		}

		if query.UpdateDate == "" {
			return fmt.Errorf("missing update date for query %d  (queries.updateDate)", i)
		}

		if err := validateJsonFields(query.JSONFields, i); err != nil {
			return err
		}
	}
	return nil
}

func validateJsonFields(_ []JSONField, _ int) error {
	//TODO
	return nil
}
