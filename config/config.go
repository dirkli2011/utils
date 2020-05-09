package config

import (
	"strings"

	"github.com/dirkli2011/utils/config/ini"
	"github.com/dirkli2011/utils/config/json"
	"github.com/dirkli2011/utils/env"
	"github.com/dirkli2011/utils/file"
)

type cfg struct {
	ini  *ini.ConfigFile
	json *json.ConfigJson
	env  string
}

var Config = &cfg{}

func init() {
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
