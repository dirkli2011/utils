package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dirkli2011/utils"
	"github.com/dirkli2011/utils/cache"
	"github.com/gomodule/redigo/redis"
)

type Cache struct {
	pool     *redis.Pool //连接池
	conninfo string      //连接配置信息
	dbNum    int         // select 的库序号
	key      string      //key前缀
	password string      //密码
	maxIdle  int         //连接池中的最大空闲连接数
}

var DefaultKey = "redisCache"

func init() {
	cache.Register("redis", Newcache)
}

func Newcache() cache.Cache {
	return &Cache{key: DefaultKey}
}

func (self *Cache) Start(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["conn"]; !ok {
		return errors.New("config has no conn key")
	}

	if _, ok := cf["dbNum"]; !ok {
		cf["dbNum"] = "0"
	}
	if _, ok := cf["password"]; !ok {
		cf["password"] = ""
	}
	if _, ok := cf["maxIdle"]; !ok {
		cf["maxIdle"] = "3"
	}

	if cf["key"] != "" {
		self.key = cf["key"]
	}
	self.conninfo = cf["conn"]
	self.dbNum, _ = strconv.Atoi(cf["dbNum"])
	self.password = cf["password"]
	self.maxIdle, _ = strconv.Atoi(cf["maxIdle"])
	self.connect()

	c := self.pool.Get()
	defer c.Close()
	return c.Err()
}

func (self *Cache) connect() {
	if self.pool != nil {
		return

	}

	dial := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", self.conninfo)
		if err != nil {
			return nil, err
		}

		if self.password != "" {
			if _, err = c.Do("AUTH", self.password); err != nil {
				return nil, err
			}
		}

		_, selecterr := c.Do("SELECT", self.dbNum)
		if selecterr != nil {
			c.Close()
			return nil, selecterr
		}
		return
	}

	self.pool = &redis.Pool{
		MaxIdle:     self.maxIdle,
		IdleTimeout: 180 * time.Second,
		Dial:        dial,
	}
}

func (self *Cache) do(cmd string, args ...interface{}) (reply interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("cmd missing required arguments")
	}

	args[0] = fmt.Sprintf("%s:%s", self.key, args[0])
	c := self.pool.Get()
	defer c.Close()
	return c.Do(cmd, args...)
}

func (self *Cache) Get(key string) (string, error) {
	v, err := self.do("GET", key)

	if err == nil {
		return utils.GetString(v), nil
	}
	return "", err
}

func (self *Cache) Set(key string, val string, timeout int32) error {
	if timeout <= 0 {
		timeout = 3600
	}
	_, err := self.do("SETEX", key, int64(timeout), val)
	return err
}

func (self *Cache) Delete(key string) error {
	_, err := self.do("DEL", key)
	return err
}

func (self *Cache) Incr(key string) (uint64, error) {
	v, err := self.do("INCRBY", key, 1)
	if err != nil {
		return uint64(0), err
	}
	return uint64(utils.GetInt64(v)), nil
}

func (self *Cache) Decr(key string) (uint64, error) {
	v, err := self.do("INCRBY", key, -1)
	if err != nil {
		return uint64(0), err
	}

	return uint64(utils.GetInt64(v)), nil
}

func (self *Cache) Exist(key string) bool {
	v, err := redis.Bool(self.do("EXISTS", key))
	if err != nil {
		return false
	}
	return v
}

func (self *Cache) GetMulti(keys []string) (map[string]string, error) {

	c := self.pool.Get()
	defer c.Close()
	var args []interface{}
	for _, key := range keys {
		args = append(args, fmt.Sprintf("%s:%s", self.key, key))
	}
	values, err := redis.Values(c.Do("MGET", args...))
	if err != nil {
		return nil, err
	}

	vals := make(map[string]string)
	for k, v := range values {
		vals[keys[k]] = utils.GetString(v)
	}
	return vals, nil
}

func (self *Cache) ClearAll() error {
	c := self.pool.Get()
	defer c.Close()
	cachedKeys, err := redis.Strings(c.Do("KEYS", self.key+":*"))
	if err != nil {
		return err
	}
	for _, str := range cachedKeys {
		if _, err = c.Do("DEL", str); err != nil {
			return err
		}
	}
	return err
}
