// need use safe map
package env

import (
	"fmt"
	"os"
	"strings"
)

var env = make(map[string]string)

func init() {
	for _, e := range os.Environ() {
		splits := strings.Split(e, "=")
		Set(splits[0], splits[1])
	}
}

func Get(key string, defVal ...string) string {
	if val, ok := env[key]; ok {
		return val
	}
	if defVal != nil {
		return defVal[0]
	}
	return ""
}

func MustGet(key string) (string, error) {
	if val, ok := env[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("no env variable with %s", key)
}

// 不会影响其他进程的环境变量
func Set(key string, value string) {
	env[key] = value
}

// 会影响其他进程的环境变量
func MustSet(key string, value string) error {
	err := os.Setenv(key, value)
	if err != nil {
		return err
	}
	env[key] = value
	return nil
}

func GetAll() map[string]string {
	return env
}
