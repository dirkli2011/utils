package logkit

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LoggerFile struct {
	tag      string
	fileName string
	lock     *sync.Mutex
	writer   *os.File
}

func (self *LoggerFile) init(tag string) {
	self.tag = tag
	self.fileName = time.Now().Format("/2006/01/02/15")
	self.lock = new(sync.Mutex)
}

func (self *LoggerFile) write(str []byte) {
	self.openfile()
	self.lock.Lock()
	self.writer.Write(str)
	self.lock.Unlock()
}

func (self *LoggerFile) openfile() {
	fileName := time.Now().Format("/2006/01/02/15")
	filePath := logPath + "/" + self.tag + fileName
	if self.writer != nil && self.fileName == fileName && exist(filePath) {
		return
	}

	self.fileName = fileName
	self.free()

	path := getdir(filePath)
	if !exist(path) {
		if err := mkdirAll(path); err != nil {
			panic("logkit use file, but create logPath fail!!!")
		}
	}
	writer, _ := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	self.writer = writer
}

func (self *LoggerFile) flush() error {
	return nil
}

func (self *LoggerFile) free() {
	if self.writer != nil {
		self.writer.Close()
		self.writer = nil
	}
}

func exist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func getdir(path string) string {
	return filepath.Dir(path)
}

func mkdirAll(path string) error {
	if !exist(path) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
