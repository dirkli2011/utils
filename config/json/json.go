package json

import (
	"github.com/dirkli2011/utils/file"

	"github.com/tidwall/gjson"
)

type ConfigJson struct {
	data string
}

func ReadConfigData(data string) (*ConfigJson, error) {
	cfg := &ConfigJson{}
	cfg.data = data
	return cfg, nil
}

func ReadConfigFile(f string) (*ConfigJson, error) {
	c, err := file.GetContent(f)
	if err != nil {
		return nil, err
	}

	cfg := &ConfigJson{}
	cfg.data = string(c)
	return cfg, nil
}

func (c *ConfigJson) String(path string) string {
	return gjson.Get(c.data, path).String()
}

func (c *ConfigJson) Int(path string) int {
	return int(gjson.Get(c.data, path).Int())
}

func (c *ConfigJson) Int64(path string) int64 {
	return gjson.Get(c.data, path).Int()
}

func (c *ConfigJson) Bool(path string) bool {
	return gjson.Get(c.data, path).Bool()
}

func (c *ConfigJson) Float(path string) float64 {
	return gjson.Get(c.data, path).Float()
}

func (c *ConfigJson) Data() string {
	return c.data
}
