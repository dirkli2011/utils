package yaml

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/dirkli2011/utils/env"
	"github.com/dirkli2011/utils/file"
	yml "gopkg.in/yaml.v3"
)

type ConfigYaml struct {
	data map[string]map[string]interface{}
}

var boolStrings = map[string]bool{
	"0":     false,
	"1":     true,
	"false": false,
	"true":  true,
	"n":     false,
	"y":     true,
	"no":    false,
	"yes":   true,
	"off":   false,
	"on":    true,
}

func ReadConfigFile(f string) (*ConfigYaml, error) {
	c, err := file.GetContent(f)
	if err != nil {
		return nil, err
	}

	cfg := &ConfigYaml{}
	err = yml.Unmarshal(c, &cfg.data)
	if err != nil {
		return nil, err
	}
	for env, section := range cfg.data {
		for option, value := range section {
			switch value.(type) {
			case string:
				cfg.data[env][option] = parseEnv(value.(string), cfg.data[env])
			}
		}
	}
	return cfg, nil
}

func (c *ConfigYaml) String(env, key string) string {
	v := c.data[env][key]
	switch result := v.(type) {
	case string:
		return result

	case []byte:
		return string(result)
	case bool:
		if result {
			return "true"
		}
		return "false"
	default:
		if v != nil {
			return fmt.Sprint(result)
		}
	}
	return ""
}

func (c *ConfigYaml) Int(env, key string) int {
	v := c.data[env][key]
	switch result := v.(type) {
	case int:
		return result
	case int32:
		return int(result)
	case int64:
		return int(result)
	case bool:
		if result {
			return 1
		}
	default:
		if vv := fmt.Sprint(result); vv != "" {
			value, _ := strconv.Atoi(vv)
			return value
		}
	}
	return 0
}

func (c *ConfigYaml) Int64(env, key string) int64 {
	v := c.data[env][key]
	switch result := v.(type) {
	case int:
		return int64(result)
	case int32:
		return int64(result)
	case int64:
		return result
	case bool:
		if result {
			return 1
		}
	default:
		if vv := fmt.Sprint(result); vv != "" {
			value, _ := strconv.ParseInt(vv, 10, 64)
			return value
		}
	}
	return 0
}

func (c *ConfigYaml) Bool(env, key string) bool {
	v := c.data[env][key]
	switch result := v.(type) {
	case bool:
		return result
	default:
		if vv := fmt.Sprint(result); vv != "" {
			bv, ok := boolStrings[strings.ToLower(vv)]
			if !ok {
				return false
			}
			return bv
		}
	}
	return false
}

func (c *ConfigYaml) Float(env, key string) float64 {
	v := c.data[env][key]
	switch result := v.(type) {
	case float64:
		return result
	default:
		if vv := fmt.Sprint(result); vv != "" {
			value, _ := strconv.ParseFloat(vv, 64)
			return value
		}
	}
	return 0
}

func (c *ConfigYaml) Data(env string) map[string]interface{} {
	return c.data[env]
}

func parseEnv(s string, section map[string]interface{}) string {
	rst := regexp.MustCompile(`{ENV\.([_\d\w]+)}`).FindAllStringSubmatch(s, -1)
	for _, v := range rst {
		s = strings.Replace(s, v[0], env.Get(v[1]), 1)
	}

	rst = regexp.MustCompile(`{\$([_\d\w]+)}`).FindAllStringSubmatch(s, -1)
	for _, v := range rst {
		s = strings.Replace(s, v[0], fmt.Sprintf("%s", section[v[1]]), 1)
	}
	return s
}
