package test

import (
	"testing"

	"github.com/Ringloop/mr-plow/config"
)

type invalidConf1 struct{}

// 'readerTest' implementing the Interface
func (*invalidConf1) ReadConfig() ([]byte, error) {

	testComplexConfig := `
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

	// Prepare data you want to return without reading from the file
	return []byte(testComplexConfig), nil
}

func TestInvalidConfig(t *testing.T) {
	testReader := invalidConf1{}
	_, err := config.ParseConfiguration(&testReader)
	if err == nil {
		t.Errorf("Invalid config without polling returned succes")
		t.Fail()
	}

}
