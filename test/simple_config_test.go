package test

import (
	"testing"

	"dariobalinzo.com/elastic/v2/test_util"

	"dariobalinzo.com/elastic/v2/config"
)

type readerTest struct{}

// 'readerTest' implementing the Interface
func (*readerTest) ReadConfig() ([]byte, error) {

	testSimpleConfig := `
database: "databaseValue"
queries:
  - index: indexValue
    query: queryValue
`
	// Prepare data you want to return without reading from the file
	return []byte(testSimpleConfig), nil
}

func TestGetSimpleConfig(t *testing.T) {
	testReader := readerTest{}
	configVal, err := config.ParseConfiguration(&testReader)
	if err != nil {
		t.Errorf("Parsing config, got error %s", err)
	}

	test_util.AssertEqual(t, err, nil)
	test_util.AssertEqual(t, configVal.Queries[0].Index, "indexValue")
	test_util.AssertEqual(t, configVal.Queries[0].Query, "queryValue")
	test_util.AssertEqual(t, configVal.Database, "databaseValue")

}
