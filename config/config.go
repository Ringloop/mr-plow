package config

import (
	"errors"
	"fmt"
	"log"

	"gopkg.in/yaml.v2"
)

type JSONFields []struct {
	FieldName  string `yaml:"fieldName"`
	Attributes struct {
		AttributeName string `yaml:"attributeName"`
		AttributeType string `yaml:"attributeType"`
	} `yaml:"attributes"`
}

type QueryModel struct {
	Index      string     `yaml:"index"`
	Query      string     `yaml:"query"`
	UpdateDate string     `yaml:"updateDate"`
	JSONFields JSONFields `yaml:"JSONFields"`
	Id         string     `yaml:"id"`
}

type ElasticConfig struct {
	Url        string `yaml:"url"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	CaCertPath string `yaml:"caCertPath"`
	NumWorker  int    `yaml:"NumWorker"`
}

type ImportConfig struct {
	Database string        `yaml:"database"`
	Queries  []QueryModel  `yaml:"queries"`
	Elastic  ElasticConfig `yaml:"elastic"`
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

func validateJsonFields(_ JSONFields, _ int) error {
	//TODO
	return nil
}
