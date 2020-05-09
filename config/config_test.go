package config

import (
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dirkli2011/utils/env"
)

var testData map[string]interface{}

func Init() {
	// load conf
	env.Set("ENV", "dev")
	env.Set("CONF_PATH_INI", "./data/demo.ini")
	env.Set("CONF_PATH_JSON", "./data/demo.json")
	env.Set("CONF_PATH_YAML", "./data/demo.yaml")
	Reload()

	u, _ := user.Current()
	testData = map[string]interface{}{
		"boolval":    true,
		"mysql.host": "dev.com",
		"mysql.port": 3306,
		"app_port":   8888,
		"user":       u.Name,
	}
}

func TestConfig(t *testing.T) {
	Init()
	for key, expect := range testData {
		switch expect.(type) {
		case bool:
			assert.Equal(t, expect, IniBool(key))
			assert.Equal(t, expect, JsonBool(key))
		case int, int32:
			assert.Equal(t, expect, IniInt(key))
		case int64:
			assert.Equal(t, expect, IniInt64(key))
		case float32, float64:
			assert.Equal(t, expect, IniFloat(key))
		case string, []byte:
			assert.Equal(t, expect, IniString(key))
		}
	}
}
