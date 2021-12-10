package config

import (
	"gopkg.in/yaml.v2"
)

type QueryModel struct {
	Index  string `yaml:"index"`
	Query  string `yaml:"query"`
	Fields []struct {
		Name string `yaml:"name"`
		Type string `yaml:"type"`
	} `yaml:"fields"`
	JSONFields []struct {
		FieldName  string `yaml:"fieldName"`
		Attributes []struct {
			AttributeName string `yaml:"attributeName"`
			AttributeType string `yaml:"attributeType"`
		} `yaml:"attributes"`
	} `yaml:"JSONFields"`
}

type ImportConfig struct {
	Database string       `yaml:"database"`
	Queries  []QueryModel `yaml:"queries"`
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
