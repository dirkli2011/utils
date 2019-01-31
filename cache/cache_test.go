package cache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	memory, err := NewCache("memory", `{"interval": 1}`)
	if err != nil {
		t.Error("init err")
	}

	key := "dirkli"
	val := "1111"
	if err = memory.Set(key, val, 2); err != nil {
		t.Error("set error", err)
	}

	v, err := memory.Get("dirkli")
	if err != nil {
		t.Error("get error", err)
	}
	if v != val {
		t.Error("get val fail")
	}

	if !memory.Exist(key) {
		t.Error("exist, but return fail")
	}

	err = memory.Delete(key)
	if err != nil {
		t.Error("delete error", err)
	}

	if memory.Exist(key) {
		t.Error("not exist, but return true")
	}

	num, err := memory.Incr(key)
	if err != nil {
		t.Error("Incr error", err)
	}

	if num != 1 {
		t.Error("Incr fail")
	}

	for i := 0; i < 10; i++ {
		num, err = memory.Incr(key)
	}
	if num != 11 {
		t.Error("Incr fail")
	}

	for i := 0; i < 5; i++ {
		num, err = memory.Decr(key)
	}
	if num != 6 {
		t.Error("Decr fail")
	}

	valConf := map[string]string{
		"k1": "v1",
		"k2": "v2",
	}

	for k, v := range valConf {
		memory.Set(k, v, 3)
	}

	vals, err := memory.GetMulti([]string{"k1", "k2"})
	if err != nil {
		t.Error("GetMulti error", err)
	}

	for k, v := range valConf {
		if v != vals[k] {
			t.Error("GetMulti fail")
		}
	}
	time.Sleep(4 * time.Second)
	if memory.Exist("k1") || memory.Exist("k2") {
		t.Error("interval delete expired key fail")
	}

	memory.ClearAll()
	if memory.Exist(key) {
		t.Error("expect not exist, but return true")
	}

}
