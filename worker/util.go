package worker

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
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

//IfToISlice converts slice of interface value to int slice
func IfToISlice(value []interface{}) ([]int, error) {
	var (
		intVal = make([]int, len(value))
		err    error
	)
	for index, curValue := range value {
		if intVal[index], err = IfToI(curValue); err != nil {
			break
		}
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

//IfToFSlice converts slice of interface value to float slice
func IfToFSlice(value []interface{}) ([]float64, error) {
	var (
		floatVal = make([]float64, len(value))
		err      error
	)
	for index, curValue := range value {
		if floatVal[index], err = IfToF(curValue); err != nil {
			break
		}
	}
	return floatVal, err
}

//IfToA converts interface value to string
func IfToA(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

//IfToASlice converts slice of interface value to string slice
func IfToASlice(value []interface{}) []string {
	var floatVal = make([]string, len(value))
	for index, curValue := range value {
		floatVal[index] = IfToA(curValue)
	}
	return floatVal
}

func getSlowDuration(jobPool *JobPool) (duration time.Duration) {
	duration = jobPool.slowDuration
	if duration == 0 {
		duration = 90 * time.Second
	}
	return
}
