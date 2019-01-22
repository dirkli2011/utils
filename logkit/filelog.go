package logkit

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LoggerFile struct {
	fileName string
	lock     *sync.Mutex
	writer   *os.File
}

func (self *LoggerFile) init() {
	self.lock = new(sync.Mutex)
}

func (self *LoggerFile) write(str []byte) {
	self.openfile()
	self.lock.Lock()
	self.writer.Write(str)
	self.lock.Unlock()
}

func (self *LoggerFile) openfile() {
	fileName := time.Now().Format(format_filerotate)
	filePath := logPath + "/" + logName + fileName
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
	self.writer, _ = os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
}

func (self *LoggerFile) flush() error {
	self.free()
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
