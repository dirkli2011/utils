package logkit

import (
	"bytes"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

const (
	format_datetime   = "2006-01-02 15:04:05"
	format_filerotate = "/2006/01/02/15"
	default_tag       = "_default"
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
	init()
	write(str []byte)
	flush() error
	free()
}

func Free() {
	logger.free()
}

func Flush() error {
	return logger.flush()
}

func Debug(str string, tag ...string) {
	if level > LEVEL_DEBUG {
		return
	}
	logger.write(formater(str2bytes(str), "debug", tag...))
}

func Info(str string, tag ...string) {
	if level > LEVEL_INFO {
		return
	}
	logger.write(formater(str2bytes(str), "info", tag...))
}

func Warn(str string, tag ...string) {
	if level > LEVEL_WARN {
		return
	}
	logger.write(formater(str2bytes(str), "warn", tag...))
}

func Error(str string, tag ...string) {
	if level > LEVEL_ERROR {
		return
	}
	logger.write(formater(str2bytes(str), "error", tag...))
}

func Fatal(str string, tag ...string) {
	if level > LEVEL_FATAL {
		return
	}
	logger.write(formater(str2bytes(str), "fatal", tag...))
}

// 用于需要报警的日志
func Alarm(str string, tag ...string) {
	logger.write(formater(str2bytes(str), "alarm", tag...))
}

func formater(str []byte, level string, args ...string) []byte {
	tag := default_tag
	if len(args) > 0 {
		tag = args[0]
	}
	tags, who := caller(tag)
	now := time.Now().Format(format_datetime)

	var buffer bytes.Buffer
	buffer.WriteString(now)
	buffer.WriteString(" tag[")
	buffer.WriteString(tags)
	buffer.WriteString("] ")
	buffer.WriteString("caller[")
	buffer.WriteString(who)
	buffer.WriteString("] [")
	buffer.WriteString(level)
	buffer.WriteString("] ")
	buffer.Write(str)
	buffer.WriteString("\n")
	return buffer.Bytes()
}

func caller(tag string) (string, string) {
	tags := logEnv + ",&" + tag
	pc, file, line, _ := runtime.Caller(3)
	name := runtime.FuncForPC(pc).Name()

	idx := strings.LastIndex(file, "src")
	if idx > 0 {
		file = file[idx:]
	}
	return tags, file + ":" + strconv.Itoa(line) + " " + name
}

// 直接转换指针类型，数据不会被复制
func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
