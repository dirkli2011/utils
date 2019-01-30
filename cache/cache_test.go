package cache

import (
	"testing"
)

func TestMemcache(t *testing.T) {

	cache, err := NewCache("memcache", `{"conn": "127.0.0.1:11211"}`)
	if err != nil {
		t.Error("init err")
	}

	if err = cache.Set("dirkli", "1", 10); err != nil {
		t.Error("set error", err)
	}

	if _, err := cache.Get("dirkli"); err != nil {
		t.Error("get error", err)
	}
}
