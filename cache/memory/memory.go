package memory

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/dirkli2011/utils/cache"
)

var DefaultInterval = 60 //每分钟清理一次过期缓存

type MemoryCacheItem struct {
	val     string
	ctime   time.Time
	expired time.Duration
}

func (self *MemoryCacheItem) isExpired() bool {
	if self.expired == 0 {
		return false // 永不过期
	}
	return time.Now().Sub(self.ctime) > self.expired
}

type MemoryCache struct {
	sync.RWMutex
	interval time.Duration
	items    map[string]*MemoryCacheItem
}

func init() {
	cache.Register("memory", Newcache)
}

func Newcache() cache.Cache {
	return &MemoryCache{items: make(map[string]*MemoryCacheItem)}
}

func (self *MemoryCache) Start(config string) error {
	var cf map[string]int
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["interval"]; !ok {
		cf = make(map[string]int)
		cf["interval"] = DefaultInterval
	}

	self.interval = time.Duration(cf["interval"]) * time.Second
	go self.check()
	return nil
}

// 定期处理过期数据
func (self *MemoryCache) check() {

	if self.interval < 1 {
		return
	}

	for {
		<-time.After(self.interval)
		if self.items == nil {
			continue
		}

		for key, item := range self.items {
			if item.isExpired() {
				self.Lock()
				delete(self.items, key)
				self.Unlock()
			}
		}
	}
}

func (self *MemoryCache) Get(key string) (string, error) {
	self.RLock()
	defer self.RUnlock()
	if v, ok := self.items[key]; ok {
		if !v.isExpired() {
			return v.val, nil
		}
	}
	return "", nil
}

func (self *MemoryCache) Set(key string, val string, timeout int32) error {
	self.Lock()
	defer self.Unlock()
	self.items[key] = &MemoryCacheItem{
		val:     val,
		ctime:   time.Now(),
		expired: time.Duration(timeout) * time.Second,
	}
	return nil
}

func (self *MemoryCache) Delete(key string) error {
	self.Lock()
	defer self.Unlock()
	if _, ok := self.items[key]; ok {
		delete(self.items, key)
	}
	return nil
}

func (self *MemoryCache) Incr(key string) (uint64, error) {
	self.Lock()
	defer self.Unlock()

	v, ok := self.items[key]
	if !ok {
		self.items[key] = &MemoryCacheItem{
			val:     "1",
			ctime:   time.Now(),
			expired: 0,
		}
		return uint64(1), nil
	}

	val, _ := strconv.Atoi(v.val)
	val++
	self.items[key] = &MemoryCacheItem{
		val:     strconv.Itoa(val),
		ctime:   time.Now(),
		expired: 0,
	}
	return uint64(val), nil
}

func (self *MemoryCache) Decr(key string) (uint64, error) {
	self.RLock()
	defer self.RUnlock()

	v, ok := self.items[key]
	if !ok {
		self.items[key] = &MemoryCacheItem{
			val:     "0",
			ctime:   time.Now(),
			expired: 0,
		}
		return uint64(0), nil
	}

	val, _ := strconv.Atoi(v.val)
	val--
	if val < 0 {
		val = 0
	}
	self.items[key] = &MemoryCacheItem{
		val:     strconv.Itoa(val),
		ctime:   time.Now(),
		expired: 0,
	}

	return uint64(val), nil
}

func (self *MemoryCache) Exist(key string) bool {
	self.RLock()
	defer self.RUnlock()
	if v, ok := self.items[key]; ok {
		return !v.isExpired()
	}
	return false
}

func (self *MemoryCache) GetMulti(keys []string) (map[string]string, error) {
	vals := make(map[string]string)
	self.RLock()
	defer self.RUnlock()
	for _, key := range keys {
		val, err := self.Get(key)
		if err != nil {
			return nil, err
		}
		vals[key] = val
	}
	return vals, nil
}

func (self *MemoryCache) ClearAll() error {
	self.Lock()
	defer self.Unlock()
	self.items = make(map[string]*MemoryCacheItem)
	return nil
}
