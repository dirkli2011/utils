package ssdb

import (
	"testing"

	"github.com/dirkli2011/utils/cache"
)

func TestSsdb(t *testing.T) {

	ssdb, err := cache.NewCache("ssdb", `{"conn": "127.0.0.1:8888"}`)
	if err != nil {
		t.Error("init err")
	}

	key := "dirkli"
	val := "1111"
	if err = ssdb.Set(key, val, 10); err != nil {
		t.Error("set error", err)
	}

	v, err := ssdb.Get("dirkli")
	if err != nil {
		t.Error("get error", err)
	}
	if v != val {
		t.Error("get val fail")
	}

	if !ssdb.Exist(key) {
		t.Error("exist, but return fail")
	}

	err = ssdb.Delete(key)
	if err != nil {
		t.Error("delete error", err)
	}
	if ssdb.Exist(key) {
		t.Error("not exist, but return true")
	}

	num, err := ssdb.Incr(key)
	if err != nil {
		t.Error("Incr error", err)
	}
	if num != 1 {
		t.Error("Incr fail")
	}

	for i := 0; i < 10; i++ {
		num, err = ssdb.Incr(key)
	}
	if num != 11 {
		t.Error("Incr fail")
	}

	for i := 0; i < 5; i++ {
		num, err = ssdb.Decr(key)
	}
	if num != 6 {
		t.Error("Decr fail")
	}

	valConf := map[string]string{
		"k1": "v1",
		"k2": "v2",
	}

	for k, v := range valConf {
		ssdb.Set(k, v, 0)
	}

	vals, err := ssdb.GetMulti([]string{"k1", "k2"})
	if err != nil {
		t.Error("GetMulti error", err)
	}

	for k, v := range valConf {
		if v != vals[k] {
			t.Error("GetMulti fail")
		}
	}

	err = ssdb.ClearAll()
	if err != nil {
		t.Error("ClearAll error", err)
	}

	if ssdb.Exist(key) {
		t.Error("expect not exist, but return true")
	}

}
