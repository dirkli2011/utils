package logkit

import (
	"os"
	"path/filepath"

	"github.com/dirkli2011/utils/env"
)

// logger配置读取和初始化
func init() {
	logName = env.Get("logkit.name", "main")
	logType = env.Get("logkit.type", "std")
	logPath = env.Get("logkit.path", pwd())
	logLevel = env.Get("logkit.level", "debug")
	logEnv = env.Get("logkit.env", env.Get("USER"))

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
