package casting

import (
	"reflect"
	"strconv"
	"strings"
)

func CastArrayOfData(inputTypeMap map[string]string, inputNameArray []string, inputDataArray []interface{}) []interface{} {
	castedColumns := make([]interface{}, len(inputDataArray))

	for i := range inputDataArray {
		if column_type, ok := inputTypeMap[inputNameArray[i]]; ok {
			switch strings.ToLower(column_type) {
			case "string":
				castedColumns[i] = castToString(inputDataArray[i])
			case "integer":
				castedColumns[i] = castToInt(inputDataArray[i])
			case "float":
				castedColumns[i] = castToFloat(inputDataArray[i])
			case "boolean":
				castedColumns[i] = castToBool(inputDataArray[i])
			default:
				castedColumns[i] = inputDataArray[i]
			}
		} else {
			castedColumns[i] = inputDataArray[i]
		}
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

func castToInt(inputVar interface{}) int64 {
	switch varType := reflect.TypeOf(inputVar).String(); varType {
	case "string":
		res, err := strconv.ParseInt(inputVar.(string), 10, 64) //have to manage this error
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
		return int64(inputVar.(float64))
	}
	return inputVar.(int64)
}

func castToFloat(inputVar interface{}) float64 {
	switch varType := reflect.TypeOf(inputVar).String(); varType {
	case "string":
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
