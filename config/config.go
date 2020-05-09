package config

import (
	"strings"

	"github.com/dirkli2011/utils/config/ini"
	"github.com/dirkli2011/utils/config/json"
	"github.com/dirkli2011/utils/config/yaml"
	"github.com/dirkli2011/utils/env"
	"github.com/dirkli2011/utils/file"
)

type cfg struct {
	ini  *ini.ConfigFile
	json *json.ConfigJson
	yaml *yaml.ConfigYaml
	env  string
}

var Config = &cfg{}

func init() {
	ini.DefaultSection = "common"
	Reload()
}

func Reload() {
	Config.env = strings.ToLower(env.Get("ENV"))

	// ini格式文件配置读取
	iniPath := env.Get("CONF_PATH_INI")
	if file.Exist(iniPath) {
		Config.ini, _ = ini.ReadConfigFile(iniPath)
	}

	// json格式文件配置读取
	jsonPath := env.Get("CONF_PATH_JSON")
	if file.Exist(jsonPath) {
		Config.json, _ = json.ReadConfigFile(jsonPath)
	}

	// yaml格式文件配置读取
	yamlPath := env.Get("CONF_PATH_YAML")
	if file.Exist(yamlPath) {
		Config.yaml, _ = yaml.ReadConfigFile(yamlPath)
	}
}

func SetEnv(env string) {
	Config.env = strings.ToLower(env)
}

func SetDefaultSection(section string) {
	ini.DefaultSection = strings.ToLower(section)
}

func SetIniData(data string) {
	if data != "" {
		Config.ini, _ = ini.ReadConfigData(data)
	}
}

func SetJsonData(data string) {
	if data != "" {
		Config.json, _ = json.ReadConfigData(data)
	}
}

// 读取所有配置文件
func String(key string, defVal ...string) string {
	switch {
	case Config.ini != nil:
		return IniString(key)
	case Config.yaml != nil:
		return YamlString(key)
	case Config.json != nil:
		return JsonString(key)
	}

	if len(defVal) > 0 {
		return defVal[0]
	}
	return ""
}

func Int(key string, defVal ...int) int {
	switch {
	case Config.ini != nil:
		return IniInt(key)
	case Config.yaml != nil:
		return YamlInt(key)
	case Config.json != nil:
		return JsonInt(key)
	}

	if len(defVal) > 0 {
		return defVal[0]
	}

	return 0
}

func Int64(key string, defVal ...int64) int64 {
	switch {
	case Config.ini != nil:
		return IniInt64(key)
	case Config.yaml != nil:
		return YamlInt64(key)
	case Config.json != nil:
		return JsonInt64(key)
	}

	if len(defVal) > 0 {
		return defVal[0]
	}
	return 0
}

func Bool(key string, defVal ...bool) bool {
	switch {
	case Config.ini != nil:
		return IniBool(key)
	case Config.yaml != nil:
		return YamlBool(key)
	case Config.json != nil:
		return JsonBool(key)
	}

	if len(defVal) > 0 {
		return defVal[0]
	}
	return false
}

func Float(key string, defVal ...float64) float64 {
	switch {
	case Config.ini != nil:
		return IniFloat(key)
	case Config.yaml != nil:
		return YamlFloat(key)
	case Config.json != nil:
		return JsonFloat(key)
	}

	if len(defVal) > 0 {
		return defVal[0]
	}
	return 0
}
