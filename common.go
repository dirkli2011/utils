package utils

import (
	"crypto/rand"
	"fmt"
	r "math/rand"
	"reflect"
	"time"
	"unsafe"
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

// 直接转换指针类型，数据不会被复制
func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
