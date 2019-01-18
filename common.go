package utils

import (
	"fmt"
	"reflect"
)

// 任意类型转换为string
func ToString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

// 判断变量是否为空
func IsEmpty(v interface{}) bool {
	val := reflect.ValueOf(v)
	valType := val.Kind()
	switch valType {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() == 0
	case reflect.Float32, reflect.Float64:
		return val.Float() == 0
	case reflect.String:
		return val.String() == ""
	case reflect.Interface, reflect.Slice, reflect.Ptr, reflect.Map, reflect.Chan, reflect.Func:
		if val.IsNil() {
			return true
		} else if valType == reflect.Slice || valType == reflect.Map {
			return val.Len() == 0
		}
	case reflect.Struct:
		fieldCount := val.NumField()
		for i := 0; i < fieldCount; i++ {
			field := val.Field(i)
			if field.IsValid() && !IsEmpty(field) {
				return false
			}
		}
		return true
	default:
		return v == nil
	}
	return false
}
