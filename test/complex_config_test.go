package test

import (
	"testing"

	"dariobalinzo.com/elastic/v2/test_util"

	"dariobalinzo.com/elastic/v2/config"
)

type readerComplexTest struct {
	fileName string
}

// 'readerTest' implementing the Interface
func (r *readerComplexTest) ReadConfig() ([]byte, error) {

	testComplexConfig := `
sql: "sqlValue"
query: "queryValue"
index: "indexValue"
queries:
  - query: "query_0_Value"
    index: "index_0_Value"
  - query: "query_1_Value"
    index: "index_1_Value"
`

	// Prepare data you want to return without reading from the file
	return []byte(testComplexConfig), nil
}

func TestGetComplexConfig(t *testing.T) {
	testReader := readerComplexTest{fileName: "Sample File Name"}
	configVal, err := config.ParseConfiguration(&testReader)
	if err != nil {
		t.Errorf("Parsing config, got error %s", err)
	}

	test_util.AssertEqual(t, err, nil)
	test_util.AssertEqual(t, configVal.Index, "indexValue")
	test_util.AssertEqual(t, configVal.Query, "queryValue")
	test_util.AssertEqual(t, configVal.SqlConfig, "sqlValue")
	test_util.AssertEqual(t, configVal.Queries[0].Query, "query_0_Value")
	test_util.AssertEqual(t, configVal.Queries[0].Index, "index_0_Value")
	test_util.AssertEqual(t, configVal.Queries[1].Query, "query_1_Value")
	test_util.AssertEqual(t, configVal.Queries[1].Index, "index_1_Value")
}
