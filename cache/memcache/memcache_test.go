package memcache

import (
	"testing"

	"github.com/dirkli2011/utils/cache"
)

func TestCache(t *testing.T) {

	memcache, err := cache.NewCache("memcache", `{"conn": "127.0.0.1:11211"}`)
	if err != nil {
		t.Error("init err")
	}

	key := "dirkli"
	val := "1111"
	if err = memcache.Set(key, val, 10); err != nil {
		t.Error("set error", err)

	}

	v, err := memcache.Get("dirkli")
	if err != nil {
		t.Error("get error", err)
	}
	if v != val {
		t.Error("get val fail")
	}

	if !memcache.Exist(key) {
		t.Error("exist, but return fail")
	}

	err = memcache.Delete(key)
	if err != nil {
		t.Error("delete error", err)
	}
	if memcache.Exist(key) {
		t.Error("not exist, but return true")
	}

	num, err := memcache.Incr(key)
	if err != nil {
		t.Error("Incr error", err)
	}
	if num != 1 {
		t.Error("Incr fail")
	}

	for i := 0; i < 10; i++ {
		num, err = memcache.Incr(key)
	}
	if num != 11 {
		t.Error("Incr fail")
	}

	for i := 0; i < 5; i++ {
		num, err = memcache.Decr(key)
	}
	if num != 6 {
		t.Error("Decr fail")
	}

	valConf := map[string]string{
		"k1": "v1",
		"k2": "v2",
	}

	for k, v := range valConf {
		memcache.Set(k, v, 0)
	}

	vals, err := memcache.GetMulti([]string{"k1", "k2"})
	if err != nil {
		t.Error("GetMulti error", err)
	}

	for k, v := range valConf {
		if v != vals[k] {
			t.Error("GetMulti fail")
		}
	}

	err = memcache.ClearAll()
	if err != nil {
		t.Error("ClearAll error", err)
	}

	if memcache.Exist(key) {
		t.Error("expect not exist, but return true")
	}

}
