package config

func YamlString(key string) string {
	if Config.yaml == nil {
		return ""
	}
	return Config.yaml.String(Config.env, key)
}

func YamlInt(key string) int {
	if Config.yaml == nil {
		return 0
	}
	return Config.yaml.Int(Config.env, key)
}

func YamlInt64(key string) int64 {
	if Config.yaml == nil {
		return 0
	}
	return Config.yaml.Int64(Config.env, key)
}

func YamlFloat(key string) float64 {
	if Config.yaml == nil {
		return 0
	}
	return Config.yaml.Float(Config.env, key)
}

func YamlBool(key string) bool {
	if Config.yaml == nil {
		return false
	}
	return Config.yaml.Bool(Config.env, key)
}

func YamlAll() map[string]interface{} {
	if Config.yaml == nil {
		return nil
	}
	return Config.yaml.Data(Config.env)
}
