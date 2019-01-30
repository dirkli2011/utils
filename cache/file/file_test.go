package file

import (
	"testing"

	"github.com/dirkli2011/utils/cache"
)

func TestCache(t *testing.T) {

	file, err := cache.NewCache("file", `{}`)

	if err != nil {
		t.Error("init err")
	}

	key := "dirkli"
	val := "1111"
	if err = file.Set(key, val, 10); err != nil {
		t.Error("set error", err)

	}

	v, err := file.Get("dirkli")
	if err != nil {
		t.Error("get error", err)
	}
	if v != val {
		t.Error("get val fail")
	}

	if !file.Exist(key) {
		t.Error("exist, but return fail")
	}

	err = file.Delete(key)
	if err != nil {
		t.Error("delete error", err)
	}
	if file.Exist(key) {
		t.Error("not exist, but return true")
	}

	num, err := file.Incr(key)
	if err != nil {
		t.Error("Incr error", err)
	}
	if num != 1 {
		t.Error("Incr fail")
	}

	for i := 0; i < 10; i++ {
		num, err = file.Incr(key)
	}
	if num != 11 {
		t.Error("Incr fail")
	}

	for i := 0; i < 5; i++ {
		num, err = file.Decr(key)
	}
	if num != 6 {
		t.Error("Decr fail")
	}

	valConf := map[string]string{
		"k1": "v1",
		"k2": "v2",
	}

	for k, v := range valConf {
		file.Set(k, v, 0)
	}

	vals, err := file.GetMulti([]string{"k1", "k2"})
	if err != nil {
		t.Error("GetMulti error", err)
	}

	for k, v := range valConf {
		if v != vals[k] {
			t.Error("GetMulti fail")
		}
	}

	err = file.ClearAll()
	if err != nil {
		t.Error("ClearAll error", err)
	}

	if file.Exist(key) {
		t.Error("expect not exist, but return true")
	}

}
