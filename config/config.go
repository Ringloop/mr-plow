package config

import (
	"gopkg.in/yaml.v2"
)

type QueryModel struct {
	Index      string `yaml:"index"`
	Query      string `yaml:"query"`
	JSONFields []struct {
		FieldName  string `yaml:"fieldName"`
		Attributes struct {
			AttributeName string `yaml:"attributeName"`
			AttributeType string `yaml:"attributeType"`
		} `yaml:"attributes"`
	} `yaml:"JSONFields"`
}

type ImportConfig struct {
	Database string `yaml:"database"`
	Index    string `yaml:"index"`
	Query    string `yaml:"query"`
	Queries  []struct {
		Index      string `yaml:"index"`
		Query      string `yaml:"query"`
		JSONFields []struct {
			FieldName  string `yaml:"fieldName"`
			Attributes struct {
				AttributeName string `yaml:"attributeName"`
				AttributeType string `yaml:"attributeType"`
			} `yaml:"attributes"`
		} `yaml:"JSONFields"`
	} `yaml:"queries"`
}

func ParseConfiguration(reader IReader) (*ImportConfig, error) {

	yamlFile, err := reader.ReadConfig()
	if err != nil {
		return nil, err
	}

	importConfiguration := &ImportConfig{}
	//importConfiguration := make(map[interface{}]interface{})
	err = yaml.Unmarshal(yamlFile, importConfiguration)
	if err != nil {
		return &ImportConfig{}, err
	}

	return importConfiguration, nil
}
