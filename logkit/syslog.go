package logkit

import (
	"log/syslog"
)

type LoggerSyslog struct {
	writer *syslog.Writer
}

func (self *LoggerSyslog) init() {
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
