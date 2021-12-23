package test

import (
	"testing"

	"github.com/Ringloop/mr-plow/test_util"

	"github.com/Ringloop/mr-plow/config"
)

type readerCompleteTest struct{}

// 'readerTest' implementing the Interface
func (*readerCompleteTest) ReadConfig() ([]byte, error) {

	testCompleteConfig := `
pollingSeconds: 5
database: databaseValue
queries:
  - index: index_1
    query: select * from table_1
    updateDate: test01
    JSONFields:
      - fieldName: dataField_1
        attributes:
          attributeName: attribute_1_Name
          attributeType: attribute_1_Type
      - fieldName: dataField_2
        attributes:
          attributeName: attribute_2_Name
          attributeType: attribute_2_Type
  - index: index_2
    query: select * from table_2
    updateDate: test02
    JSONFields:
      - fieldName: dataField_2
        attributes:
          attributeName: attribute_1_Name_2
          attributeType: attribute_1_Type_2
elastic:
  url: http://localhost:9200
`

	// Prepare data you want to return without reading from the file
	return []byte(testCompleteConfig), nil
}

func TestGetCompleteConfig(t *testing.T) {
	testReader := readerCompleteTest{}
	configVal, err := config.ParseConfiguration(&testReader)
	if err != nil {
		t.Errorf("Parsing config, got error %s", err)
	}

	test_util.AssertEqual(t, err, nil)
	test_util.AssertEqual(t, configVal.Database, "databaseValue")
	queries := configVal.Queries
	test_util.AssertEqual(t, len(queries), 2)

	//test queries[0]
	test_util.AssertEqual(t, queries[0].Index, "index_1")
	test_util.AssertEqual(t, queries[0].Query, "select * from table_1")
	jsonFields1 := queries[0].JSONFields
	test_util.AssertEqual(t, len(jsonFields1), 2)
	test_util.AssertEqual(t, jsonFields1[0].FieldName, "dataField_1")
	attribute1JsonFields1 := jsonFields1[0].Attributes
	test_util.AssertEqual(t, attribute1JsonFields1.AttributeName, "attribute_1_Name")
	test_util.AssertEqual(t, attribute1JsonFields1.AttributeType, "attribute_1_Type")
	test_util.AssertEqual(t, jsonFields1[1].FieldName, "dataField_2")
	attribute1JsonFields2 := jsonFields1[1].Attributes
	test_util.AssertEqual(t, attribute1JsonFields2.AttributeName, "attribute_2_Name")
	test_util.AssertEqual(t, attribute1JsonFields2.AttributeType, "attribute_2_Type")

	//test queries[1]
	test_util.AssertEqual(t, queries[1].Index, "index_2")
	test_util.AssertEqual(t, queries[1].Query, "select * from table_2")
	jsonFields2 := queries[1].JSONFields
	test_util.AssertEqual(t, len(jsonFields2), 1)
	test_util.AssertEqual(t, jsonFields2[0].FieldName, "dataField_2")
	attribute2JsonFields1 := jsonFields2[0].Attributes
	test_util.AssertEqual(t, attribute2JsonFields1.AttributeName, "attribute_1_Name_2")
	test_util.AssertEqual(t, attribute2JsonFields1.AttributeType, "attribute_1_Type_2")
}
