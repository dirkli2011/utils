package logkit

import (
	"fmt"
	"io"
	"os"
)

type LoggerStd struct {
	tag    string
	writer io.WriteCloser
}

func (self *LoggerStd) init(tag string) {
	self.tag = tag
	self.writer = os.Stderr
}

func (self *LoggerStd) write(str []byte) {
	fmt.Fprint(self.writer, string(str))
}

func (self *LoggerStd) flush() error {
	return nil
}

func (self *LoggerStd) free() {
	self.writer.Close()
}
