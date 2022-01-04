package test

import (
	"testing"

	"github.com/Ringloop/mr-plow/config"
)

type invalidConf1 struct{}

func (*invalidConf1) ReadConfig() ([]byte, error) {
	yml := `
#pollingSeconds: 5 missing polling...
database: "databaseValue"
queries:
  - query: "query_0_Value"
    index: "index_0_Value"
    updateDate: "test0"
  - query: "query_1_Value"
    index: "index_1_Value"
    updateDate: "test1"
elastic:
  url: http://localhost:9200
`
	return []byte(yml), nil
}

type invalidConf2 struct{}

func (*invalidConf2) ReadConfig() ([]byte, error) {
	yml := `
pollingSeconds: 5
#database: "databaseValue" missing db...
queries:
  - query: "query_0_Value"
    index: "index_0_Value"
    updateDate: "test0"
  - query: "query_1_Value"
    index: "index_1_Value"
    updateDate: "test1"
elastic:
  url: http://localhost:9200
`
	return []byte(yml), nil
}

type invalidConf3 struct{}

func (*invalidConf3) ReadConfig() ([]byte, error) {
	yml := `
pollingSeconds: 5
database: "databaseValue"
elastic:
  url: http://localhost:9200
`
	return []byte(yml), nil
}

type invalidConf4 struct{}

func (*invalidConf4) ReadConfig() ([]byte, error) {
	yml := `
pollingSeconds: 5
#database: "databaseValue" missing db...
queries:
  - query: "query_0_Value"
    index: "index_0_Value"
    updateDate: "test0"
  - query: "query_1_Value"
    index: "index_1_Value"
    updateDate: "test1"
`
	return []byte(yml), nil
}

func TestInvalidConfig(t *testing.T) {
	_, err := config.ParseConfiguration(&invalidConf1{})
	if err == nil {
		t.Errorf("Invalid config without polling returned succes")
		t.Fail()
	}

	_, err = config.ParseConfiguration(&invalidConf2{})
	if err == nil {
		t.Errorf("Invalid config without database returned succes")
		t.Fail()
	}

	_, err = config.ParseConfiguration(&invalidConf3{})
	if err == nil {
		t.Errorf("Invalid config without queries returned succes")
		t.Fail()
	}

	_, err = config.ParseConfiguration(&invalidConf4{})
	if err == nil {
		t.Errorf("Invalid config without elastic url returned succes")
		t.Fail()
	}

}
