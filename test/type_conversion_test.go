package test

import (
	"testing"

	"github.com/Ringloop/mr-plow/casting"
	"github.com/stretchr/testify/require"
)

//func CastSingleElement(inputTypeMap map[string]string, inputName string, inputData interface{}) interface{} {
func prepareInputType() map[string]string {

	inputTypeMap := make(map[string]string)
	inputTypeMap["stringElement"] = "string"
	inputTypeMap["intElement"] = "integer"
	inputTypeMap["boolElement"] = "boolean"
	inputTypeMap["floatElement"] = "float"

	return inputTypeMap
}

func TestStringConvertion(t *testing.T) {
	var intElement int = 5
	var stringElement string = "StringValue"
	var boolElement bool = true
	var floatElement float64 = 4.54

	var ok bool

	inputMap := prepareInputType()

	var convertedString = casting.CastSingleElement(inputMap, "stringElement", intElement)
	_, ok = convertedString.(string)
	require.True(t, ok)
	require.Equal(t, convertedString, "5")
	convertedString = casting.CastSingleElement(inputMap, "stringElement", stringElement)
	_, ok = convertedString.(string)
	require.True(t, ok)
	require.Equal(t, convertedString, "StringValue")
	convertedString = casting.CastSingleElement(inputMap, "stringElement", boolElement)
	_, ok = convertedString.(string)
	require.True(t, ok)
	require.Equal(t, convertedString, "true")
	convertedString = casting.CastSingleElement(inputMap, "stringElement", floatElement)
	_, ok = convertedString.(string)
	require.True(t, ok)
	require.Equal(t, convertedString, "4.54")

}

func TestIntegerConvertion(t *testing.T) {
	var intElement int = 5
	var stringElement string = "5"
	var boolElement bool = true
	var floatElement float64 = 5.

	var ok bool

	inputMap := prepareInputType()

	var convertedInt = casting.CastSingleElement(inputMap, "intElement", intElement)
	_, ok = convertedInt.(int)
	require.True(t, ok)
	require.Equal(t, convertedInt, 5)
	convertedInt = casting.CastSingleElement(inputMap, "intElement", stringElement)
	_, ok = convertedInt.(int)
	require.True(t, ok)
	require.Equal(t, convertedInt, 5)
	convertedInt = casting.CastSingleElement(inputMap, "intElement", boolElement)
	_, ok = convertedInt.(int)
	require.True(t, ok)
	require.Equal(t, convertedInt, 1)
	convertedInt = casting.CastSingleElement(inputMap, "intElement", floatElement)
	_, ok = convertedInt.(int)
	require.True(t, ok)
	require.Equal(t, convertedInt, 5)

}

func TestFloatConvertion(t *testing.T) {
	var intElement int = 5
	var stringElement string = "5"
	var boolElement bool = true
	var floatElement float64 = 5.

	var ok bool

	inputMap := prepareInputType()

	var convertedFloat = casting.CastSingleElement(inputMap, "floatElement", intElement)
	_, ok = convertedFloat.(float64)
	require.True(t, ok)
	require.Equal(t, convertedFloat, 5.)
	convertedFloat = casting.CastSingleElement(inputMap, "floatElement", stringElement)
	_, ok = convertedFloat.(float64)
	require.True(t, ok)
	require.Equal(t, convertedFloat, 5.)
	convertedFloat = casting.CastSingleElement(inputMap, "floatElement", boolElement)
	_, ok = convertedFloat.(float64)
	require.True(t, ok)
	require.Equal(t, convertedFloat, 1.)
	convertedFloat = casting.CastSingleElement(inputMap, "floatElement", floatElement)
	_, ok = convertedFloat.(float64)
	require.True(t, ok)
	require.Equal(t, convertedFloat, 5.)

}

func TestBooleanConvertion(t *testing.T) {
	var intElementFalse int = 0
	var intElementDefault int = 4
	var stringElementTrue string = "true"
	var stringElementDefault string = "AnyOtherValue"
	var boolElement bool = true
	var floatElementFalse float64 = 0.
	var floatElementDefault float64 = 5.

	var ok bool

	inputMap := prepareInputType()

	var convertedBool = casting.CastSingleElement(inputMap, "boolElement", intElementFalse)
	_, ok = convertedBool.(bool)
	require.True(t, ok)
	require.Equal(t, convertedBool, false)
	convertedBool = casting.CastSingleElement(inputMap, "boolElement", intElementDefault)
	_, ok = convertedBool.(bool)
	require.True(t, ok)
	require.Equal(t, convertedBool, true)
	convertedBool = casting.CastSingleElement(inputMap, "boolElement", stringElementTrue)
	_, ok = convertedBool.(bool)
	require.True(t, ok)
	require.Equal(t, convertedBool, true)
	convertedBool = casting.CastSingleElement(inputMap, "boolElement", stringElementDefault)
	_, ok = convertedBool.(bool)
	require.True(t, ok)
	require.Equal(t, convertedBool, false)
	convertedBool = casting.CastSingleElement(inputMap, "boolElement", boolElement)
	_, ok = convertedBool.(bool)
	require.True(t, ok)
	require.Equal(t, convertedBool, true)
	convertedBool = casting.CastSingleElement(inputMap, "boolElement", floatElementFalse)
	_, ok = convertedBool.(bool)
	require.True(t, ok)
	require.Equal(t, convertedBool, false)
	convertedBool = casting.CastSingleElement(inputMap, "boolElement", floatElementDefault)
	_, ok = convertedBool.(bool)
	require.True(t, ok)
	require.Equal(t, convertedBool, true)

}
