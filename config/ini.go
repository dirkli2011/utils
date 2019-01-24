package config

func IniString(key string) string {
	if Config.ini == nil {
		return ""
	}
	val, err := Config.ini.GetString(Config.env, key)
	if err != nil {
		val, _ = Config.ini.GetString("default", key)
	}
	return val
}

func IniInt(key string) int {
	if Config.ini == nil {
		return 0
	}
	val, err := Config.ini.GetInt(Config.env, key)
	if err != nil {
		val, _ = Config.ini.GetInt("default", key)
	}
	return val
}

func IniInt64(key string) int64 {
	if Config.ini == nil {
		return 0
	}
	val, err := Config.ini.GetInt64(Config.env, key)
	if err != nil {
		val, _ = Config.ini.GetInt64("default", key)
	}
	return val
}

func IniBool(key string) bool {
	if Config.ini == nil {
		return false
	}
	val, err := Config.ini.GetBool(Config.env, key)
	if err != nil {
		val, _ = Config.ini.GetBool("default", key)
	}
	return val
}

func IniFloat(key string) float64 {
	if Config.ini == nil {
		return float64(0)
	}
	val, err := Config.ini.GetFloat(Config.env, key)
	if err != nil {
		val, _ = Config.ini.GetFloat("default", key)
	}
	return val
}
