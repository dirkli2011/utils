package config

import (
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dirkli2011/utils/env"
)

func Init() {
	// load ini conf
	env.Set("ENV", "dev")
	env.Set("CONF_PATH_INI", "./data/demo.ini")
	SetDefaultSection("common")
	Reload()
}

func TestConfig(t *testing.T) {
	Init()

	u, _ := user.Current()
	testData := map[string]interface{}{
		"mysql.host": "dev.com",
		"mysql.port": 3306,
		"app_port":   8888,
		"user":       u.Name,
	}

	for key, expect := range testData {
		switch expect.(type) {
		case bool:
			assert.Equal(t, expect, IniBool(key))
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
