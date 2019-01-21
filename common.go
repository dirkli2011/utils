package utils

import (
	"crypto/rand"
	"fmt"
	r "math/rand"
	"reflect"
	"time"
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

// 随机返回一个指定长度的字符串
func RandString(n int, chars ...byte) string {
	if len(chars) == 0 {
		chars = []byte(`0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`)
	}

	var bytes = make([]byte, n)
	var randBy bool
	if num, err := rand.Read(bytes); num != n || err != nil {
		r.Seed(time.Now().UnixNano())
		randBy = true
	}

	for i, b := range bytes {
		if randBy {
			bytes[i] = chars[r.Intn(len(chars))]
		} else {
			bytes[i] = chars[b%byte(len(chars))]
		}
	}
	return string(bytes)
}
