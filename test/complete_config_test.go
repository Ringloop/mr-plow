package test

import (
	"testing"

	"dariobalinzo.com/elastic/v2/test_util"

	"dariobalinzo.com/elastic/v2/config"
)

type readerCompleteTest struct{}

// 'readerTest' implementing the Interface
func (*readerCompleteTest) ReadConfig() ([]byte, error) {

	testCompleteConfig := `
database: databaseValue
queries:
  - index: index_1
    query: select * from table_1
    fields:
      - name: name
        type: String
      - name: last_update
        type: Date
    JSONFields:
      - fieldName: dataField_1
        attributes:
          - attributeName: attribute_1_Name
            attributeType: attribute_1_Type
      - fieldName: dataField_2
        attributes:
          - attributeName: attribute_2_Name
            attributeType: attribute_2_Type
          - attributeName: attribute_2_1_Name
            attributeType: attribute_2_1_Type
  - index: index_2
    query: select * from table_2

    JSONFields:
      - fieldName: dataField_2
        attributes:
          - attributeName: attribute_1_Name_2
            attributeType: attribute_1_Type_2
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
	queryFields := queries[0].Fields
	test_util.AssertEqual(t, queryFields[0].Name, "name")
	test_util.AssertEqual(t, queryFields[0].Type, "String")
	test_util.AssertEqual(t, queryFields[1].Name, "last_update")
	test_util.AssertEqual(t, queryFields[1].Type, "Date")
	jsonFields1 := queries[0].JSONFields
	test_util.AssertEqual(t, len(jsonFields1), 2)
	test_util.AssertEqual(t, jsonFields1[0].FieldName, "dataField_1")
	attribute1JsonFields1 := jsonFields1[0].Attributes
	test_util.AssertEqual(t, attribute1JsonFields1[0].AttributeName, "attribute_1_Name")
	test_util.AssertEqual(t, attribute1JsonFields1[0].AttributeType, "attribute_1_Type")
	test_util.AssertEqual(t, jsonFields1[1].FieldName, "dataField_2")
	attribute1JsonFields2 := jsonFields1[1].Attributes
	test_util.AssertEqual(t, attribute1JsonFields2[0].AttributeName, "attribute_2_Name")
	test_util.AssertEqual(t, attribute1JsonFields2[0].AttributeType, "attribute_2_Type")
	test_util.AssertEqual(t, attribute1JsonFields2[1].AttributeName, "attribute_2_1_Name")
	test_util.AssertEqual(t, attribute1JsonFields2[1].AttributeType, "attribute_2_1_Type")

	//test queries[1]
	test_util.AssertEqual(t, queries[1].Index, "index_2")
	test_util.AssertEqual(t, queries[1].Query, "select * from table_2")
	jsonFields2 := queries[1].JSONFields
	test_util.AssertEqual(t, len(jsonFields2), 1)
	test_util.AssertEqual(t, jsonFields2[0].FieldName, "dataField_2")
	attribute2JsonFields1 := jsonFields2[0].Attributes
	test_util.AssertEqual(t, attribute2JsonFields1[0].AttributeName, "attribute_1_Name_2")
	test_util.AssertEqual(t, attribute2JsonFields1[0].AttributeType, "attribute_1_Type_2")
}
