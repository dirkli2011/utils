package config

func JsonString(path string) string {
	if Config.json == nil {
		return ""
	}
	return Config.json.String(path)
}

func JsonInt(path string) int {
	if Config.json == nil {
		return 0
	}
	return Config.json.Int(path)
}

func JsonInt64(path string) int64 {
	if Config.json == nil {
		return 0
	}
	return Config.json.Int64(path)
}

func JsonBool(path string) bool {
	if Config.json == nil {
		return false
	}
	return Config.json.Bool(path)
}

func JsonFloat(path string) float64 {
	if Config.json == nil {
		return float64(0)
	}
	return Config.json.Float(path)
}

func JsonAll() string {
	if Config.json == nil {
		return ""
	}
	return Config.json.Data()
}
