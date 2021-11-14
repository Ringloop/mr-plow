package config

import (
	"gopkg.in/yaml.v2"
)

type QueryModel struct {
	Query string `yaml:"query"`
	Index string `yaml:"index"`
}

type ImportConfig struct {
	SqlConfig string `yaml:"sql"`
	Query     string `yaml:"query"`
	Index     string `yaml:"index"`
	Queries   []struct {
		Query string `yaml:"query"`
		Index string `yaml:"index"`
	}
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

	return importConfiguration, nil
}
