package logkit

import (
	"fmt"
	"io"
	"os"
)

type LoggerAsync struct {
	tag    string
	writer io.WriteCloser
}

func (self *LoggerAsync) init(tag string) {
	self.tag = tag
	self.writer = os.Stderr
}

func (self *LoggerAsync) write(str []byte) {
	fmt.Fprint(self.writer, string(str))
}

func (self *LoggerAsync) flush() error {
	return nil
}

func (self *LoggerAsync) free() {
	self.writer.Close()
}
