package cache

import (
	"fmt"
)

type Cache interface {
	Get(key string) (string, error)
	Set(key string, val string, timeout int32) error
	Delete(key string) error
	Incr(key string) (uint64, error)
	Decr(key string) (uint64, error)
	Exist(key string) bool
	GetMulti(keys []string) (map[string]string, error)
	ClearAll() error
	Start(config string) error
}

type InstaceFunc func() Cache

var adapters = make(map[string]InstaceFunc)

// cache需要在init中调用该方法注册一个初始化方法
func Register(name string, adapter InstaceFunc) {
	if adapter == nil {
		panic("cache " + name + " register fail, adapter is nil")
	}

	if _, ok := adapters[name]; !ok {
		adapters[name] = adapter
	}
}

func NewCache(name string, config string) (cache Cache, err error) {
	insFunc, ok := adapters[name]
	if !ok {
		err = fmt.Errorf("unknow cache type: %s, may be not import???", name)
		return
	}

	cache = insFunc()
	err = cache.Start(config)
	if err != nil {
		cache = nil
	}
	return
}
