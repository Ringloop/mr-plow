package casting

import (
	"reflect"
	"strconv"
	"strings"
)

type Converter struct {
	inputTypeMap map[string]string
}

func NewConverter(inputTypeMap map[string]string) *Converter {
	converter := &Converter{
		inputTypeMap: inputTypeMap}
	return converter
}

func (converter *Converter) CastSingleElement(inputName string, inputData interface{}) interface{} {
	if columnType, ok := converter.inputTypeMap[inputName]; ok {
		switch strings.ToLower(columnType) {
		case "string":
			return castToString(inputData)
		case "integer":
			return castToInt(inputData)
		case "float":
			return castToFloat(inputData)
		case "boolean":
			return castToBool(inputData)
		default:
			return inputData
		}
	} else {
		return inputData
	}
}

func (converter *Converter) CastArrayOfData(inputNameArray []string, inputDataArray []interface{}) []interface{} {
	castedColumns := make([]interface{}, len(inputDataArray))

	for i := range inputDataArray {
		castedColumns[i] = converter.CastSingleElement(inputNameArray[i], inputDataArray[i])
	}

	return castedColumns
}

func castToString(inputVar interface{}) string {
	switch varType := reflect.TypeOf(inputVar).String(); varType {
	case "bool":
		return strconv.FormatBool(inputVar.(bool))
	case "float64":
		return strconv.FormatFloat(inputVar.(float64), 'f', -1, 64)
	case "int":
		return strconv.Itoa(inputVar.(int))
	}
	return inputVar.(string)
}

func castToInt(inputVar interface{}) int {
	switch varType := reflect.TypeOf(inputVar).String(); varType {
	case "string":
		res, err := strconv.Atoi(inputVar.(string)) //have to manage this error
		if err == nil {
			return res
		}
	case "bool":
		if inputVar.(bool) {
			return 1
		} else {
			return 0
		}
	case "float64":
		return int(inputVar.(float64))
	case "int64":
		return int(inputVar.(int64))
	}

	return inputVar.(int)
}

func castToFloat(inputVar interface{}) float64 {
	switch varType := reflect.TypeOf(inputVar).String(); varType {
	case "string":
		if inputVar == "" {
			return (0.)
		}
		inputVar = strings.Replace(inputVar.(string), ",", ".", -1)
		res, err := strconv.ParseFloat(inputVar.(string), 64) //have to manage this error
		if err == nil {
			return res
		}
	case "bool":
		if inputVar.(bool) {
			return 1.
		} else {
			return 0.
		}
	case "int":
		return float64(inputVar.(int))
	}
	return inputVar.(float64)
}

func castToBool(inputVar interface{}) bool {
	switch varType := reflect.TypeOf(inputVar).String(); varType {
	case "string":
		if strings.EqualFold("true", inputVar.(string)) {
			return true
		} else {
			return false
		}
	case "float64":
		if inputVar.(float64) == 0. {
			return false
		} else {
			return true
		}
	case "int":
		if inputVar.(int) == 0 {
			return false
		} else {
			return true
		}
	}
	return inputVar.(bool)
}
