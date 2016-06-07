package utils

import (
	"reflect"
	"strconv"
)

const (
	StructStringMaxDefault = 1024000
)

// 合法性检测
func StrcutValid(v interface{}) (err error) {
	// 只检测struct
	structVal := reflect.Indirect(reflect.ValueOf(v))
	if structVal.Kind() != reflect.Struct {
		return
	}
	// 检测每个字段
	structType := structVal.Type()
	for i := 0; i < structType.NumField(); i++ {
		// ignore embedded filed
		fieldType := structType.Field(i)
		//if fieldType.Anonymous == true {
		//	println("igigigigig", fieldType.Name)
		//	continue
		//}
		//println("debug structValid", fieldType.Name)
		// init
		var (
			maxLength = StructStringMaxDefault
			minLength = 0
			maxTag    = fieldType.Tag.Get("maxLength")
			minTag    = fieldType.Tag.Get("minLength")
		)
		if maxTag != "" {
			if maxLength, err = strconv.Atoi(maxTag); err != nil {
				break
			}
		}
		if minTag != "" {
			if minLength, err = strconv.Atoi(minTag); err != nil {
				break
			}
		}
		// ptr to struct and check
		val := reflect.Indirect(structVal.Field(i))
		switch val.Kind() {
		case reflect.String:
			// string
			if (maxLength > 0) && (val.Len() > maxLength) {
				// println("debug", fieldType.Name, maxLength)
				err = ErrStringMax
			} else if val.Len() < minLength {
				err = ErrStringMin
			}
			break
		case reflect.Struct:
			//println("can addr", fieldType.Name, val.CanAddr(), val.CanSet())
			if (fieldType.Name[0] >= 65) && (fieldType.Name[0] <= 90) {
				//println("can exported")
				err = StrcutValid(val.Interface())
			}
			break
		}
		// error
		if err != nil {
			break
		}
	}
	return
}
