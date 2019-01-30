package redis

import (
	"testing"

	"github.com/dirkli2011/utils/cache"
)

func TestRedis(t *testing.T) {

	redis, err := cache.NewCache("redis", `{"conn": "127.0.0.1:6379", "password":"111111"}`)
	if err != nil {
		t.Error("init err")
	}

	key := "dirkli"
	val := "1111"
	if err = redis.Set(key, val, 10); err != nil {
		t.Error("set error", err)

	}

	v, err := redis.Get("dirkli")
	if err != nil {
		t.Error("get error", err)
	}
	if v != val {
		t.Error("get val fail")
	}

	if !redis.Exist(key) {
		t.Error("exist, but return fail")
	}

	err = redis.Delete(key)
	if err != nil {
		t.Error("delete error", err)
	}
	if redis.Exist(key) {
		t.Error("not exist, but return true")
	}

	num, err := redis.Incr(key)
	if err != nil {
		t.Error("Incr error", err)
	}
	if num != 1 {
		t.Error("Incr fail")
	}

	for i := 0; i < 10; i++ {
		num, err = redis.Incr(key)
	}
	if num != 11 {
		t.Error("Incr fail")
	}

	for i := 0; i < 5; i++ {
		num, err = redis.Decr(key)
	}
	if num != 6 {
		t.Error("Decr fail")
	}

	valConf := map[string]string{
		"k1": "v1",
		"k2": "v2",
	}

	for k, v := range valConf {
		redis.Set(k, v, 0)
	}

	vals, err := redis.GetMulti([]string{"k1", "k2"})
	if err != nil {
		t.Error("GetMulti error", err)
	}

	for k, v := range valConf {
		if v != vals[k] {
			t.Error("GetMulti fail")
		}
	}

	err = redis.ClearAll()
	if err != nil {
		t.Error("ClearAll error", err)
	}

	if redis.Exist(key) {
		t.Error("expect not exist, but return true")
	}

}
