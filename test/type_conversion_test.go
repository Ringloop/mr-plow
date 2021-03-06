package test

import (
	"testing"

	"github.com/Ringloop/mr-plow/internal/casting"
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
	var stringEmptyElement string = ""
	var boolElement bool = true
	var floatElement float64 = 4.54

	var ok bool

	inputMap := prepareInputType()

	converter := casting.NewConverter(inputMap)
	var convertedString = converter.CastSingleElement("stringElement", intElement)
	_, ok = convertedString.(string)
	require.True(t, ok)
	require.Equal(t, convertedString, "5")
	convertedString = converter.CastSingleElement("stringElement", stringEmptyElement)
	_, ok = convertedString.(string)
	require.True(t, ok)
	require.Equal(t, convertedString, "")
	convertedString = converter.CastSingleElement("stringElement", stringElement)
	_, ok = convertedString.(string)
	require.True(t, ok)
	require.Equal(t, convertedString, "StringValue")
	convertedString = converter.CastSingleElement("stringElement", boolElement)
	_, ok = convertedString.(string)
	require.True(t, ok)
	require.Equal(t, convertedString, "true")
	convertedString = converter.CastSingleElement("stringElement", floatElement)
	_, ok = convertedString.(string)
	require.True(t, ok)
	require.Equal(t, convertedString, "4.54")

}

func TestIntegerConvertion(t *testing.T) {
	var intElement int = 5
	var stringElement string = "5"
	var stringEmptyElement string = ""
	var stringSpacesElement string = "    "
	var stringAlphanumericElement string = "sglhsdg8478"
	var boolElement bool = true
	var floatElement float64 = 5.

	var ok bool

	inputMap := prepareInputType()

	converter := casting.NewConverter(inputMap)
	var convertedInt = converter.CastSingleElement("intElement", intElement)
	_, ok = convertedInt.(int)
	require.True(t, ok)
	require.Equal(t, convertedInt, 5)
	convertedInt = converter.CastSingleElement("intElement", stringEmptyElement)
	_, ok = convertedInt.(int)
	require.True(t, ok)
	require.Equal(t, convertedInt, 0)
	convertedInt = converter.CastSingleElement("intElement", stringSpacesElement)
	_, ok = convertedInt.(int)
	require.True(t, ok)
	require.Equal(t, convertedInt, 0)
	convertedInt = converter.CastSingleElement("intElement", stringAlphanumericElement)
	_, ok = convertedInt.(int)
	require.True(t, ok)
	require.Equal(t, convertedInt, 0)
	convertedInt = converter.CastSingleElement("intElement", stringElement)
	_, ok = convertedInt.(int)
	require.True(t, ok)
	require.Equal(t, convertedInt, 5)
	convertedInt = converter.CastSingleElement("intElement", boolElement)
	_, ok = convertedInt.(int)
	require.True(t, ok)
	require.Equal(t, convertedInt, 1)
	convertedInt = converter.CastSingleElement("intElement", floatElement)
	_, ok = convertedInt.(int)
	require.True(t, ok)
	require.Equal(t, convertedInt, 5)

}

func TestFloatConvertion(t *testing.T) {
	var intElement int = 5
	var stringElement string = "5"
	var stringElementComma string = "1,024,543.22"
	var stringEmptyElement string = ""
	var stringSpacesElement string = ""
	var stringAlphanumericElement string = "sdgsogdiso94,g09"
	var boolElement bool = true
	var floatElement float64 = 5.

	var ok bool

	inputMap := prepareInputType()
	converter := casting.NewConverter(inputMap)

	var convertedFloat = converter.CastSingleElement("floatElement", intElement)
	_, ok = convertedFloat.(float64)
	require.True(t, ok)
	require.Equal(t, convertedFloat, 5.)
	convertedFloat = converter.CastSingleElement("floatElement", stringEmptyElement)
	_, ok = convertedFloat.(float64)
	require.True(t, ok)
	require.Equal(t, convertedFloat, 0.)
	convertedFloat = converter.CastSingleElement("floatElement", stringSpacesElement)
	_, ok = convertedFloat.(float64)
	require.True(t, ok)
	require.Equal(t, convertedFloat, 0.)
	convertedFloat = converter.CastSingleElement("floatElement", stringAlphanumericElement)
	_, ok = convertedFloat.(float64)
	require.True(t, ok)
	require.Equal(t, convertedFloat, 0.)
	convertedFloat = converter.CastSingleElement("floatElement", stringElement)
	_, ok = convertedFloat.(float64)
	require.True(t, ok)
	require.Equal(t, convertedFloat, 5.)
	convertedFloat = converter.CastSingleElement("floatElement", stringElementComma)
	_, ok = convertedFloat.(float64)
	require.True(t, ok)
	require.Equal(t, convertedFloat, 1024543.22)
	convertedFloat = converter.CastSingleElement("floatElement", boolElement)
	_, ok = convertedFloat.(float64)
	require.True(t, ok)
	require.Equal(t, convertedFloat, 1.)
	convertedFloat = converter.CastSingleElement("floatElement", floatElement)
	_, ok = convertedFloat.(float64)
	require.True(t, ok)
	require.Equal(t, convertedFloat, 5.)

}

func TestBooleanConvertion(t *testing.T) {
	var intElementFalse int = 0
	var intElementDefault int = 4
	var stringEmptyElement string = ""
	var stringSpacesElement string = "    "
	var stringAlphanumericElement string = "sglhsdg8478"
	var stringElementTrue string = "true"
	var stringElementDefault string = "AnyOtherValue"
	var boolElement bool = true
	var floatElementFalse float64 = 0.
	var floatElementDefault float64 = 5.

	var ok bool

	inputMap := prepareInputType()
	converter := casting.NewConverter(inputMap)

	var convertedBool = converter.CastSingleElement("boolElement", intElementFalse)
	_, ok = convertedBool.(bool)
	require.True(t, ok)
	require.Equal(t, convertedBool, false)
	convertedBool = converter.CastSingleElement("boolElement", intElementDefault)
	_, ok = convertedBool.(bool)
	require.True(t, ok)
	require.Equal(t, convertedBool, true)
	convertedBool = converter.CastSingleElement("boolElement", stringEmptyElement)
	_, ok = convertedBool.(bool)
	require.True(t, ok)
	require.Equal(t, convertedBool, false)
	convertedBool = converter.CastSingleElement("boolElement", stringSpacesElement)
	_, ok = convertedBool.(bool)
	require.True(t, ok)
	require.Equal(t, convertedBool, false)
	convertedBool = converter.CastSingleElement("boolElement", stringAlphanumericElement)
	_, ok = convertedBool.(bool)
	require.True(t, ok)
	require.Equal(t, convertedBool, false)
	convertedBool = converter.CastSingleElement("boolElement", stringElementTrue)
	_, ok = convertedBool.(bool)
	require.True(t, ok)
	require.Equal(t, convertedBool, true)
	convertedBool = converter.CastSingleElement("boolElement", stringElementDefault)
	_, ok = convertedBool.(bool)
	require.True(t, ok)
	require.Equal(t, convertedBool, false)
	convertedBool = converter.CastSingleElement("boolElement", boolElement)
	_, ok = convertedBool.(bool)
	require.True(t, ok)
	require.Equal(t, convertedBool, true)
	convertedBool = converter.CastSingleElement("boolElement", floatElementFalse)
	_, ok = convertedBool.(bool)
	require.True(t, ok)
	require.Equal(t, convertedBool, false)
	convertedBool = converter.CastSingleElement("boolElement", floatElementDefault)
	_, ok = convertedBool.(bool)
	require.True(t, ok)
	require.Equal(t, convertedBool, true)

}
