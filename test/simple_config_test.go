package test

import (
	"testing"

	"dariobalinzo.com/elastic/v2/test_util"

	"dariobalinzo.com/elastic/v2/config"
)

type readerTest struct {
	fileName string
}

// 'readerTest' implementing the Interface
func (r *readerTest) ReadConfig() ([]byte, error) {

	testSimpleConfig := `
sql: "sqlValue"
query: "queryValue"
index: "indexValue"
`
	// Prepare data you want to return without reading from the file
	return []byte(testSimpleConfig), nil
}

func TestGetSimpleConfig(t *testing.T) {
	testReader := readerTest{fileName: "Sample File Name"}
	configVal, err := config.ParseConfiguration(&testReader)
	if err != nil {
		t.Errorf("Parsing config, got error %s", err)
	}

	test_util.AssertEqual(t, err, nil)
	test_util.AssertEqual(t, configVal.Index, "indexValue")
	test_util.AssertEqual(t, configVal.Query, "queryValue")
	test_util.AssertEqual(t, configVal.SqlConfig, "sqlValue")

}
