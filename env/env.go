package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/dirkli2011/utils/cmap"
)

var envMap cmap.ConcurrentMap

func init() {
	envMap = cmap.New()
	for _, e := range os.Environ() {
		splits := strings.Split(e, "=")
		envMap.Set(splits[0], splits[1])
	}
}

func Get(key string, defVal ...string) string {
	if val, ok := envMap.Get(key); ok {
		return val.(string)
	}

	if len(defVal) > 0 {
		return defVal[0]
	}

	return ""
}

func MustGet(key string) (string, error) {
	if val, ok := envMap.Get(key); ok {
		return val.(string), nil
	}
	return "", fmt.Errorf("no env variable with %s", key)
}

// 不会影响其他进程的环境变量
func Set(key string, value string) {
	envMap.Set(key, value)
}

// 会影响其他进程的环境变量
func MustSet(key string, value string) error {
	err := os.Setenv(key, value)
	if err != nil {
		return err
	}
	envMap.Set(key, value)
	return nil
}

func GetAll() map[string]string {
	envs := make(map[string]string, envMap.Len())
	fn := func(key string, value interface{}) {
		envs[key] = value.(string)
	}
	envMap.IteratorBy(fn)
	return envs
}

func Has(key string) bool {
	return envMap.Exist(key)
}
