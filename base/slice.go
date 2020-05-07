package base

import (
	"errors"
	"reflect"
)

func InSliceWithError(s interface{}, v interface{}) (exist bool, err error) {
	vs := reflect.ValueOf(s)
	if vs.Kind() != reflect.Slice {
		err = errors.New("param s must be a slice")
		return
	}

	if vs.Len() == 0 {
		exist = false
		return
	}

	stype := reflect.TypeOf(s).String()[2:]
	if stype != "interface {}" && stype != reflect.TypeOf(v).String() {
		err = errors.New("type of param v is not match param s")
		return
	}

	for i := 0; i < vs.Len(); i++ {
		if vs.Index(i).Interface() == v {
			exist = true
			return
		}
	}

	return
}

// 判存
func InSlice(s interface{}, v interface{}) bool {
	exist, _ := InSliceWithError(s, v)
	return exist
}

func UniqSliceWithError(s interface{}) (r interface{}, err error) {
	vs := reflect.ValueOf(s)
	if vs.Kind() != reflect.Slice {
		err = errors.New("param s must be a slice")
		return
	}

	length := vs.Len()

	if length == 0 {
		r = s
		return
	}

	seen := make(map[interface{}]bool, length)
	j := 0
	for i := 0; i < length; i++ {
		val := vs.Index(i)
		v := val.Interface()
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = true
		vs.Index(j).Set(val)
		j++
	}

	return vs.Slice(0, j).Interface(), nil
}

// 排重
func UniqSlice(s interface{}) interface{} {
	r, err := UniqSliceWithError(s)
	if err != nil {
		r = s
	}
	return r
}
