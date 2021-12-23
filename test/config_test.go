package test

import (
	"testing"

	"github.com/Ringloop/Mr-Plow/test_util"

	"github.com/Ringloop/Mr-Plow/config"
)

type readerComplexTest struct{}

// 'readerTest' implementing the Interface
func (*readerComplexTest) ReadConfig() ([]byte, error) {

	testComplexConfig := `
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

	// Prepare data you want to return without reading from the file
	return []byte(testComplexConfig), nil
}

func TestGetComplexConfig(t *testing.T) {
	testReader := readerComplexTest{}
	configVal, err := config.ParseConfiguration(&testReader)
	if err != nil {
		t.Errorf("Parsing config, got error %s", err)
	}

	test_util.AssertEqual(t, err, nil)
	test_util.AssertEqual(t, configVal.Database, "databaseValue")
	test_util.AssertEqual(t, configVal.Queries[0].Query, "query_0_Value")
	test_util.AssertEqual(t, configVal.Queries[0].Index, "index_0_Value")
	test_util.AssertEqual(t, configVal.Queries[0].UpdateDate, "test0")
	test_util.AssertEqual(t, configVal.Queries[1].Query, "query_1_Value")
	test_util.AssertEqual(t, configVal.Queries[1].Index, "index_1_Value")
	test_util.AssertEqual(t, configVal.Queries[1].UpdateDate, "test1")
}
