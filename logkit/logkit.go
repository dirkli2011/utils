package logkit

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	format_datetime   = "2006-01-02 15:04:05"
	format_filerotate = "/2006/01/02/15"
)

const (
	LEVEL_DEBUG = iota
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
)

var Level_Config = map[string]int{
	"debug": LEVEL_DEBUG,
	"info":  LEVEL_INFO,
	"warn":  LEVEL_WARN,
	"error": LEVEL_ERROR,
	"fatal": LEVEL_FATAL,
}

var (
	logger   ILogger //日志对象
	logType  string  //log类型
	logName  string  //日志名
	logPath  string  //日志路径
	logLevel string  //日志等级
	level    int     //日志等级，通过logLevel获取
	logEnv   string  //日志环境
)

type ILogger interface {
	init(tag string)
	free()
	flush() error
	write(str []byte)
}

func Free() {
	logger.free()
}

func Flush() error {
	return logger.flush()
}

// 字符串记录日志
func Debug(str string) {
	if level > LEVEL_DEBUG {
		return
	}
	logger.write(formater(str, "debug"))
}

func Info(str string) {
	if level > LEVEL_INFO {
		return
	}
	logger.write(formater(str, "info"))
}

func Warn(str string) {
	if level > LEVEL_WARN {
		return
	}
	logger.write(formater(str, "warn"))
}

func Error(str string) {
	if level > LEVEL_ERROR {
		return
	}
	logger.write(formater(str, "error"))
}

func Fatal(str string) {
	if level > LEVEL_FATAL {
		return
	}
	logger.write(formater(str, "fatal"))
}

// 用于需要报警的日志
func Alarm(str string) {
	logger.write(formater(str, "alarm"))
}

// 格式化字符串记录日志
func Debugf(format string, params ...interface{}) {
	if level > LEVEL_DEBUG {
		return
	}
	logger.write(formater(fmt.Sprintf(format, params...), "debug"))
}

func Infof(format string, params ...interface{}) {
	if level > LEVEL_INFO {
		return
	}
	logger.write(formater(fmt.Sprintf(format, params...), "info"))
}

func Warnf(format string, params ...interface{}) {
	if level > LEVEL_WARN {
		return
	}
	logger.write(formater(fmt.Sprintf(format, params...), "warn"))
}

func Errorf(format string, params ...interface{}) {
	if level > LEVEL_ERROR {
		return
	}
	logger.write(formater(fmt.Sprintf(format, params...), "error"))
}

func Fatalf(format string, params ...interface{}) {
	if level > LEVEL_FATAL {
		return
	}
	logger.write(formater(fmt.Sprintf(format, params...), "fatal"))
}

// 用于需要报警的日志
func Alarmf(format string, params ...interface{}) {
	logger.write(formater(fmt.Sprintf(format, params...), "alarm"))
}

func formater(str string, level string) []byte {
	tag, evt := caller()
	now := time.Now().Format(format_datetime)

	var buffer bytes.Buffer
	buffer.WriteString(now)
	buffer.WriteString(" tag[")
	buffer.WriteString(tag)
	buffer.WriteString("] ")
	buffer.WriteString("caller[")
	buffer.WriteString(evt)
	buffer.WriteString("] [")
	buffer.WriteString(level)
	buffer.WriteString("] ")
	buffer.WriteString(str)
	buffer.WriteString("\n")
	return buffer.Bytes()
}

func caller() (string, string) {
	tag := logEnv + ",&" + logName
	pc, file, line, _ := runtime.Caller(3)
	name := runtime.FuncForPC(pc).Name()

	idx := strings.LastIndex(file, "src")
	if idx > 0 {
		file = file[idx:]
	}
	return tag, file + ":" + strconv.Itoa(line) + " " + name
}
