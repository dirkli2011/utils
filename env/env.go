package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/dirkli2011/utils"
)

var env *utils.SafeMap

func init() {
	env = utils.NewSafeMap()
	for _, e := range os.Environ() {
		splits := strings.Split(e, "=")
		env.Set(splits[0], splits[1])
	}
}

func Get(key string, defVal ...string) string {
	if val := env.Get(key); val != nil {
		return val.(string)
	}

	if len(defVal) > 0 {
		return defVal[0]
	}

	return ""
}

func MustGet(key string) (string, error) {
	if val := env.Get(key); val != nil {
		return val.(string), nil
	}
	return "", fmt.Errorf("no env variable with %s", key)
}

// 不会影响其他进程的环境变量
func Set(key string, value string) {
	env.Set(key, value)
}

// 会影响其他进程的环境变量
func MustSet(key string, value string) error {
	err := os.Setenv(key, value)
	if err != nil {
		return err
	}
	env.Set(key, value)
	return nil
}

func GetAll() map[string]string {
	items := env.Items()
	envs := make(map[string]string, env.Count())

	for k, v := range items {
		switch k := k.(type) {
		case string:
			switch v := v.(type) {
			case string:
				envs[k] = v
			}
		}
	}
	return envs
}
