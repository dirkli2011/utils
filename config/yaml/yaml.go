package yaml

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/dirkli2011/utils/file"
	yml "gopkg.in/yaml.v3"
)

type ConfigYaml struct {
	data map[string]interface{}
}

func ReadConfigFile(f string) (*ConfigYaml, error) {
	c, err := file.GetContent(f)
	if err != nil {
		return nil, err
	}

	cfg := &ConfigYaml{}
	yml.Unmarshal(c, &cfg.data)
	spew.Dump(cfg)
	return cfg, nil
}

// func parseEnv(s string) string {
// 	if strings.HasPrefix(s, "ENV.") {
// 		s = env.Get(strings.Replace(s, "ENV.", "", 1))
// 	}
// 	return s
// }
