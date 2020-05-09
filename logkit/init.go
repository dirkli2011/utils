package logkit

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/dirkli2011/utils/config"
)

// logger配置读取和初始化
func init() {
	InitConf(nil)
}

func InitConf(conf map[string]string) {
	u, _ := user.Current()
	logName = config.String("logkit.name", "main")
	logType = config.String("logkit.type", "std")
	logPath = config.String("logkit.path", pwd())
	logLevel = config.String("logkit.level", "debug")
	logEnv = config.String("logkit.env", u.Name)

	if len(conf) > 0 {
		logName = conf["logkit.name"]
		logType = conf["logkit.type"]
		logPath = conf["logkit.path"]
		logLevel = conf["logkit.level"]
		logEnv = conf["app_mode"]
	}

	if logEnv == "online" && logLevel == "debug" {
		logLevel = "info"
	}
	level = Level_Config[logLevel]

	switch logType {
	case "std":
		logger = new(LoggerStd)
	case "syslog":
		logger = new(LoggerSyslog)
	case "file":
		logger = new(LoggerFile)
	case "async":
		logger = new(LoggerAsync)
	default:
		logger = new(LoggerStd)
	}
	logger.init()
}

func pwd() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir + "/logs"
}
