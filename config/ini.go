package config

import (
	"fmt"

	"github.com/dirkli2011/utils/config/ini"
)

func section() string {
	return fmt.Sprintf("%s:%s", Config.env, ini.DefaultSection)
}

func IniString(key string) string {
	if Config.ini == nil {
		return ""
	}
	val, err := Config.ini.GetString(section(), key)
	if err != nil {
		val, _ = Config.ini.GetString(ini.DefaultSection, key)
	}
	return val
}

func IniInt(key string) int {
	if Config.ini == nil {
		return 0
	}
	val, err := Config.ini.GetInt(section(), key)
	if err != nil {
		val, _ = Config.ini.GetInt(ini.DefaultSection, key)
	}
	return val
}

func IniInt64(key string) int64 {
	if Config.ini == nil {
		return 0
	}
	val, err := Config.ini.GetInt64(section(), key)
	if err != nil {
		val, _ = Config.ini.GetInt64(ini.DefaultSection, key)
	}
	return val
}

func IniBool(key string) bool {
	if Config.ini == nil {
		return false
	}
	val, err := Config.ini.GetBool(section(), key)
	if err != nil {
		val, _ = Config.ini.GetBool(ini.DefaultSection, key)
	}
	return val
}

func IniFloat(key string) float64 {
	if Config.ini == nil {
		return float64(0)
	}
	val, err := Config.ini.GetFloat(section(), key)
	if err != nil {
		val, _ = Config.ini.GetFloat(ini.DefaultSection, key)
	}
	return val
}
