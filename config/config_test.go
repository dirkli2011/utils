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
	env.Set("CONF_PATH_YAML", "./data/demo.yml")
	Reload()

	u, _ := user.Current()
	testData = map[string]interface{}{
		"boolval":    true,
		"mysql_host": "dev.com",
		"mysql_port": 3306,
		"app_port":   8888,
		"domain":     u.Name + ".xxx.com",
		"api_domain": u.Name + ".api.xxx.com",
	}
}

func TestConfig(t *testing.T) {
	Init()
	for key, expect := range testData {
		switch expect.(type) {
		case bool:
			assert.Equal(t, expect, IniBool(key))
			assert.Equal(t, expect, JsonBool(key))
			assert.Equal(t, expect, YamlBool(key))
			assert.Equal(t, expect, Bool(key))
		case int, int32:
			assert.Equal(t, expect, IniInt(key))
			assert.Equal(t, expect, JsonInt(key))
			assert.Equal(t, expect, YamlInt(key))
			assert.Equal(t, expect, Int(key))
		case int64:
			assert.Equal(t, expect, IniInt64(key))
			assert.Equal(t, expect, JsonInt64(key))
			assert.Equal(t, expect, YamlInt64(key))
			assert.Equal(t, expect, Int64(key))
		case float32, float64:
			assert.Equal(t, expect, IniFloat(key))
			assert.Equal(t, expect, JsonFloat(key))
			assert.Equal(t, expect, YamlFloat(key))
			assert.Equal(t, expect, Float(key))
		case string, []byte:
			assert.Equal(t, expect, IniString(key))
			assert.Equal(t, expect, JsonString(key))
			assert.Equal(t, expect, YamlString(key))
			assert.Equal(t, expect, String(key))
		}
	}
}
