package ssdb

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/dirkli2011/utils"
	"github.com/dirkli2011/utils/cache"
	"github.com/ssdb/gossdb/ssdb"
)

type Cache struct {
	conn *ssdb.Client
	ip   string
	port int
}

func init() {
	cache.Register("ssdb", Newcache)
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
	info := strings.Split(cf["conn"], ":")
	self.ip = info[0]
	self.port, _ = strconv.Atoi(info[1])
	if self.ip != "" && self.port != 0 {
		return self.connect()
	}
	return nil
}

func (self *Cache) connect() error {
	if self.conn != nil {
		return nil
	}

	conn, err := ssdb.Connect(self.ip, self.port)
	if err != nil {
		return err
	}
	self.conn = conn
	return nil
}

func (self *Cache) Get(key string) (string, error) {
	self.connect()
	var err error
	if item, err := self.conn.Get(key); err == nil {
		return utils.GetString(item), nil
	}
	return "", err
}

func (self *Cache) Set(key string, val string, timeout int32) error {
	self.connect()
	var resp []string
	var err error
	if timeout == 0 {
		resp, err = self.conn.Do("set", key, val)
	} else {
		resp, err = self.conn.Do("setx", key, val, int(timeout))
	}

	if len(resp) == 2 && resp[0] == "ok" {
		return nil
	}
	return err
}

func (self *Cache) Delete(key string) error {
	self.connect()
	_, err := self.conn.Del(key)
	return err
}

func (self *Cache) Incr(key string) (uint64, error) {
	self.connect()
	v, err := self.conn.Do("incr", key, 1)
	if err != nil {
		return uint64(0), err
	}

	if v[0] == "ok" {
		return uint64(utils.GetInt64(v[1])), nil
	}
	return uint64(0), nil

}

func (self *Cache) Decr(key string) (uint64, error) {
	self.connect()
	v, err := self.conn.Do("incr", key, -1)
	if err != nil {
		return uint64(0), err
	}

	if v[0] == "ok" {
		return uint64(utils.GetInt64(v[1])), nil
	}
	return uint64(0), nil
}

func (self *Cache) Exist(key string) bool {
	self.connect()
	resp, err := self.conn.Do("exists", key)
	if err != nil {
		return false
	}

	if len(resp) == 2 && resp[0] == "ok" && resp[1] == "1" {
		return true
	}
	return false
}

func (self *Cache) GetMulti(keys []string) (map[string]string, error) {
	self.connect()
	res, err := self.conn.Do("multi_get", keys)
	if err != nil {
		return nil, err
	}

	vals := make(map[string]string)
	if res[0] == "ok" {
		for i := 1; i < len(res); i += 2 {
			vals[res[i]] = res[i+1]
		}
	}
	return vals, nil
}

func (self *Cache) ClearAll() error {
	self.connect()
	start, end, limit := "", "", 50
	resp, err := self.Scan(start, end, limit)
	for err == nil {
		size := len(resp)
		if size == 1 {
			return nil
		}
		keys := []string{}
		for i := 1; i < size; i += 2 {
			keys = append(keys, resp[i])
		}
		_, e := self.conn.Do("multi_del", keys)
		if e != nil {
			return e
		}
		start = resp[size-2]
		resp, err = self.Scan(start, end, limit)
	}
	return err
}

func (self *Cache) Scan(start string, end string, limit int) ([]string, error) {
	self.connect()
	resp, err := self.conn.Do("scan", start, end, limit)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
