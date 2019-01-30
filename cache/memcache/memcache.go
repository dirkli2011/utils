package memcache

import (
	"errors"
	"strings"

	"encoding/json"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/dirkli2011/utils/cache"
)

type Cache struct {
	conn     *memcache.Client
	conninfo []string
}

func init() {
	cache.Register("memcache", Newcache)
}

func Newcache() cache.Cache {
	return &Cache{}
}

func (self *Cache) Start(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["conn"]; !ok {
		return errors.New("config has no conn key")
	}
	self.conninfo = strings.Split(cf["conn"], ";")
	if self.conn == nil {
		self.connect()
	}
	return nil
}

func (self *Cache) connect() {
	if self.conn != nil {
		return
	}
	self.conn = memcache.New(self.conninfo...)
}

func (self *Cache) Get(key string) (string, error) {
	self.connect()
	var err error
	if item, err := self.conn.Get(key); err == nil {
		return string(item.Value), nil
	}
	return "", err
}

func (self *Cache) Set(key string, val string, timeout int32) error {
	self.connect()
	item := memcache.Item{Key: key, Expiration: timeout}
	item.Value = []byte(val)
	return self.conn.Set(&item)
}

func (self *Cache) Delete(key string) error {
	self.connect()
	return self.conn.Delete(key)
}

func (self *Cache) Incr(key string) (uint64, error) {
	self.connect()
	if !self.Exist(key) {
		item := memcache.Item{Key: key, Expiration: 0}
		item.Value = []byte("0")
		self.conn.Set(&item)
	}
	return self.conn.Increment(key, 1)
}

func (self *Cache) Decr(key string) (uint64, error) {
	self.connect()
	if !self.Exist(key) {
		return uint64(0), nil
	}
	return self.conn.Decrement(key, 1)
}

func (self *Cache) Exist(key string) bool {
	self.connect()
	_, err := self.conn.Get(key)
	return err == nil
}

func (self *Cache) GetMulti(keys []string) (map[string]string, error) {
	self.connect()
	items, err := self.conn.GetMulti(keys)
	if err != nil {
		return nil, err
	}

	vals := make(map[string]string)
	for k, v := range items {
		vals[k] = string(v.Value)
	}
	return vals, nil
}

func (self *Cache) ClearAll() error {
	self.connect()
	return self.conn.FlushAll()
}
