package file

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/dirkli2011/utils"
	"github.com/dirkli2011/utils/cache"
	"github.com/dirkli2011/utils/file"
)

type CacheItem struct {
	Data    string
	Expired time.Time
}

var (
	FileCachePath = "./cache"
	FileSuffix    = ".bin"
)

type Cache struct {
	CachePath string
}

func init() {
	cache.Register("file", Newcache)
}

func Newcache() cache.Cache {
	return &Cache{}
}

func (self *Cache) Start(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["path"]; !ok {
		cf["path"] = FileCachePath
	}
	self.CachePath = cf["path"]
	if !file.Exist(self.CachePath) {
		err := file.MkdirAll(self.CachePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *Cache) getFilename(key string) string {
	keyMd5 := utils.Md5(key)
	cachePath := filepath.Join(self.CachePath, keyMd5[0:2], keyMd5[2:4], keyMd5[4:6])
	return filepath.Join(cachePath, fmt.Sprintf("%s%s", keyMd5, FileSuffix))
}

func (self *Cache) Get(key string) (string, error) {
	filename := self.getFilename(key)
	if !file.Exist(filename) {
		return "", nil
	}
	var err error
	fileData, err := file.GetContent(filename)
	if err != nil {
		return "", err
	}

	var to CacheItem
	GobDecode(fileData, &to)
	if to.Expired.Before(time.Now()) {
		return "", nil
	}
	return utils.GetString(to.Data), nil
}

func (self *Cache) GetMulti(keys []string) (map[string]string, error) {
	vals := make(map[string]string)
	for _, key := range keys {
		v, err := self.Get(key)
		if err != nil {
			return nil, err
		}
		vals[key] = v
	}
	return vals, nil
}

func (self *Cache) Set(key string, val string, timeout int32) error {
	gob.Register(val)

	item := CacheItem{Data: val}
	if timeout == 0 {
		item.Expired = time.Now().Add((86400 * 365 * 10) * time.Second)
	} else {

		item.Expired = time.Now().Add(time.Duration(timeout) * time.Second)
	}
	data, err := GobEncode(item)
	if err != nil {
		return err
	}
	_, err = file.PutContent(self.getFilename(key), data)
	return err
}

func (self *Cache) Delete(key string) error {
	filename := self.getFilename(key)
	if !file.Remove(filename) {
		return fmt.Errorf("remove file fail:%s", filename)
	}
	return nil
}
func (self *Cache) Incr(key string) (uint64, error) {
	data, err := self.Get(key)
	if err != nil {
		return uint64(0), err
	}
	incr, _ := strconv.Atoi(data)
	incr++
	fmt.Println(key, incr)

	err = self.Set(key, utils.GetString(incr), 0)
	if err != nil {
		return uint64(0), err
	}
	return uint64(incr), nil
}

func (self *Cache) Decr(key string) (uint64, error) {
	data, err := self.Get(key)
	fmt.Println(key, data)
	if err != nil {
		return uint64(0), err
	}
	incr, _ := strconv.Atoi(data)
	incr--
	if incr < 0 {
		incr = 0
	}
	err = self.Set(key, utils.GetString(incr), 0)

	if err != nil {
		return uint64(0), err
	}
	return uint64(incr), nil
}

func (self *Cache) Exist(key string) bool {
	return file.Exist(self.getFilename(key))
}

func (self *Cache) ClearAll() error {
	if file.Exist(self.CachePath) {
		return os.RemoveAll(self.CachePath)
	}
	return nil
}

func GobEncode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

func GobDecode(data []byte, to *CacheItem) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(&to)
}
