package worker

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

//IfToI converts interface value to int
func IfToI(value interface{}) (int, error) {
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

//IfToF converts interface value to float
func IfToF(value interface{}) (float64, error) {
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

//IfToA converts interface value to string
func IfToA(value interface{}) string {
	return fmt.Sprintf("%v", value)
}
