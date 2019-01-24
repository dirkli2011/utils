package config

import (
	"utils/config/ini"
	"utils/config/json"
	"utils/env"
	"utils/file"
)

type cfg struct {
	ini  *ini.ConfigFile
	json *json.ConfigJson
	env  string
}

var Config = &cfg{}

func init() {
	Config.env = env.Get("ENV")

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
