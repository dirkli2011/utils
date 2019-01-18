package logkit

import (
	"os"
	"path/filepath"
)

// logger配置读取和初始化
func init() {
	logName = os.Getenv("PRJ_NAME")
	if logName == "" {
		logName = "main"
	}

	logType = os.Getenv("logkit.type")
	if logType == "" {
		logType = "file"
	}

	logPath = os.Getenv("logkit.path")
	if logPath == "" {
		logPath = pwd()
	}

	logLevel = os.Getenv("logkit.level")
	if logLevel == "" {
		logLevel = "info"
	}

	logEnv = os.Getenv("ENV")
	if logEnv == "" {
		logEnv = os.Getenv("USER")
		logLevel = "debug"
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

	logger.init(logName)
}

func pwd() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir + "/logs"
}
