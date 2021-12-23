package casting

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func CastStringWithPrintf(input_variable interface{}) string {
	return fmt.Sprintf("%s", input_variable)
}

func CastToString(inputVar interface{}) string {
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

func CastToInt(inputVar interface{}) int {
	switch varType := reflect.TypeOf(inputVar).String(); varType {
	case "string":
		res, error := strconv.Atoi(inputVar.(string)) //have to manage this error
		if error == nil {
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
	}
	return inputVar.(int)
}

func CastToFloat(inputVar interface{}) float64 {
	switch varType := reflect.TypeOf(inputVar).String(); varType {
	case "string":
		res, error := strconv.ParseFloat(inputVar.(string), 64) //have to manage this error
		if error == nil {
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

func CastToBool(inputVar interface{}) bool {
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
