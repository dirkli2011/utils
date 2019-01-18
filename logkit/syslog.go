package logkit

import (
	"log/syslog"
)

type LoggerSyslog struct {
	tag    string
	writer *syslog.Writer
}

func (self *LoggerSyslog) init(tag string) {
	self.tag = tag
	writer, _ := syslog.New(syslog.LOG_INFO|syslog.LOG_LOCAL6, logName)
	self.writer = writer
}

func (self *LoggerSyslog) write(str []byte) {
	self.writer.Info(string(str))
}

func (self *LoggerSyslog) flush() error {
	return nil
}

func (self *LoggerSyslog) free() {
	self.writer.Close()
}
