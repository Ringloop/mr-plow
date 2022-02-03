package test

import (
	"testing"

	"github.com/Ringloop/mr-plow/config"
	"github.com/stretchr/testify/require"
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
    fields:
      - name: name
        type: String
      - name: last_update
        type: Date
    JSONFields:
      - fieldName: dataField_1
        fields:
          - name: attribute_1_Name
            type: attribute_1_Type
      - fieldName: dataField_2
        fields:
          - name: attribute_2_Name
            type: attribute_2_Type
          - name: attribute_2_1_Name
            type: attribute_2_1_Type
    id: MyId_1
  - index: index_2
    query: select * from table_2
    updateDate: test02
    JSONFields:
      - fieldName: dataField_2
        fields:
          - name: attribute_1_Name_2
            type: attribute_1_Type_2
    id: MyId_2
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

	require.NoError(t, err)
	require.Equal(t, configVal.Database, "databaseValue")
	queries := configVal.Queries
	require.Equal(t, len(queries), 2)
	//test queries[0]
	require.Equal(t, queries[0].Index, "index_1")
	require.Equal(t, queries[0].Query, "select * from table_1")
	jsonFields1 := queries[0].JSONFields
	require.Equal(t, len(jsonFields1), 2)
	require.Equal(t, jsonFields1[0].FieldName, "dataField_1")
	require.Equal(t, queries[0].UpdateDate, "test01")
	require.Equal(t, queries[0].Id, "MyId_1")
	queryFields := queries[0].Fields
	require.Equal(t, queryFields[0].Name, "name")
	require.Equal(t, queryFields[0].Type, "String")
	require.Equal(t, queryFields[1].Name, "last_update")
	require.Equal(t, queryFields[1].Type, "Date")

	jsonFields1 = queries[0].JSONFields
	require.Equal(t, len(jsonFields1), 2)
	require.Equal(t, jsonFields1[0].FieldName, "dataField_1")
	attribute1JsonFields1 := jsonFields1[0]
	require.Equal(t, attribute1JsonFields1.Fields[0].Name, "attribute_1_Name")
	require.Equal(t, attribute1JsonFields1.Fields[0].Type, "attribute_1_Type")
	require.Equal(t, jsonFields1[1].FieldName, "dataField_2")
	attribute1JsonFields2 := jsonFields1[1]
	require.Equal(t, attribute1JsonFields2.Fields[0].Type, "attribute_2_Type")
	require.Equal(t, attribute1JsonFields2.Fields[0].Name, "attribute_2_Name")
	require.Equal(t, attribute1JsonFields2.Fields[1].Name, "attribute_2_1_Name")
	require.Equal(t, attribute1JsonFields2.Fields[1].Type, "attribute_2_1_Type")

	//test queries[1]
	require.Equal(t, queries[1].Index, "index_2")
	require.Equal(t, queries[1].Query, "select * from table_2")
	require.Equal(t, queries[1].UpdateDate, "test02")
	require.Equal(t, queries[1].Id, "MyId_2")
	jsonFields2 := queries[1].JSONFields
	require.Equal(t, len(jsonFields2), 1)
	require.Equal(t, jsonFields2[0].FieldName, "dataField_2")
	attribute2JsonFields1 := jsonFields2[0].Fields
	require.Equal(t, attribute2JsonFields1[0].Name, "attribute_1_Name_2")
	require.Equal(t, attribute2JsonFields1[0].Type, "attribute_1_Type_2")
}
