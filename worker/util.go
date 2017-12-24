package worker

import (
	"errors"
	"reflect"
	"strconv"
)

//InterfaceToInt return int value of given interface
func InterfaceToInt(value interface{}) (int, error) {
	var (
		intVal int
		err    error
	)
	objKind := reflect.TypeOf(value).Kind()
	switch objKind {
	case reflect.Int:
		intVal = value.(int)
	case reflect.Int64:
		intVal = int(value.(int64))
	case reflect.Float64:
		intVal = int(value.(float64))
	case reflect.Float32:
		intVal = int(value.(float32))
	case reflect.String:
		intVal, err = strconv.Atoi(value.(string))
	default:
		err = errors.New("Not Able to typecast " + objKind.String() + " to int")
	}
	return intVal, err
}

//InterfaceToFloat return float value of given interface
func InterfaceToFloat(value interface{}) (float64, error) {
	var (
		floatVal float64
		err      error
	)
	objKind := reflect.TypeOf(value).Kind()
	switch objKind {
	case reflect.Int:
		floatVal = float64(value.(int))
	case reflect.Int64:
		floatVal = float64(value.(int64))
	case reflect.Float64:
		floatVal = value.(float64)
	case reflect.Float32:
		floatVal = float64(value.(float32))
	case reflect.String:
		floatVal, err = strconv.ParseFloat(value.(string), 64)
	default:
		err = errors.New("Not Able to typecast " + objKind.String() + " to float64")
	}
	return floatVal, err
}
