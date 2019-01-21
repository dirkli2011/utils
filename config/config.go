package config

import (
	"fmt"
	"os"
	"reflect"
	"time"
)

type Configer interface {
	Set(key, val string) error
	String(key string) string
}

type Config interface {
	Parse(key string) (Configer, error)
	ParseData(data []byte) (Configer, error)
}

var cfgs = make(map[string]Config)

// 注册一个配置项到全局配置中
func Register(name string, cfg Config) {
	if cfg == nil {
		panic("register cfg is nil")
	}

	if _, ok := cfgs[name]; ok {
		panic("cfg can't register twice!!! name is " + name)
	}
	cfgs[name] = cfg
}

// 创建一个具体的配置
func NewConfig(conftype string, filename string) (Configer, error) {
	cfg, ok := cfgs[conftype]
	if !ok {
		return nil, fmt.Errorf("cfg %q is unknow", conftype)
	}
	return cfg.Parse()
}

func NewConfigData(conftype string, data []byte) (Configer, error) {
	cfg, ok := cfgs[conftype]
	if !ok {
		return nil, fmt.Errorf("cfg %q is unknow", conftype)
	}
	return cfg.ParseData()
}
