package test

import (
	"github.com/Ringloop/mr-plow/test_util"
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
database: "databaseValue"
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

type invalidConf5 struct{}

func (*invalidConf5) ReadConfig() ([]byte, error) {
	yml := `
pollingSeconds: 5
database: "databaseValue"
queries:
  - query: "query_0_Value"
    updateDate: "test0"
  - query: "query_1_Value"
    index: "index_1_Value"
    updateDate: "test1"
elastic:
  url: http://localhost:9200
`
	return []byte(yml), nil
}

type invalidConf6 struct{}

func (*invalidConf6) ReadConfig() ([]byte, error) {
	yml := `
pollingSeconds: 5
database: "databaseValue"
queries:
  - index: "index_0_Value"
    updateDate: "test0"
  - query: "query_1_Value"
    index: "index_1_Value"
    updateDate: "test1"
elastic:
  url: http://localhost:9200
`
	return []byte(yml), nil
}

type invalidConf7 struct{}

func (*invalidConf7) ReadConfig() ([]byte, error) {
	yml := `
pollingSeconds: 5
database: "databaseValue"
queries:
  - query: "query_0_Value"
    index: "index_0_Value"
  - query: "query_1_Value"
    index: "index_1_Value"
    updateDate: "test1"
elastic:
  url: http://localhost:9200
`
	return []byte(yml), nil
}

func TestInvalidConfig(t *testing.T) {
	_, err := config.ParseConfiguration(&invalidConf1{})
	test_util.AssertEqual(t, err.Error(), "missing polling seconds url (pollingSeconds)")

	_, err = config.ParseConfiguration(&invalidConf2{})
	test_util.AssertEqual(t, err.Error(), "missing database url (database)")

	_, err = config.ParseConfiguration(&invalidConf3{})
	test_util.AssertEqual(t, err.Error(), "missing query configuration (queries)")

	_, err = config.ParseConfiguration(&invalidConf4{})
	test_util.AssertEqual(t, err.Error(), "missing elastic url (elastic.url)")

	_, err = config.ParseConfiguration(&invalidConf5{})
	test_util.AssertNotNull(t, err)
	_, err = config.ParseConfiguration(&invalidConf6{})
	test_util.AssertNotNull(t, err)
	_, err = config.ParseConfiguration(&invalidConf7{})
	test_util.AssertNotNull(t, err)

}
