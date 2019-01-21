package logkit

import (
	"os"
	"strings"
	"sync"
	"time"
)

type LoggerAsync struct {
	tag      string
	fileName string
	data     []string
	lock     *sync.RWMutex
	writer   *os.File
	status   int
}

const (
	syncTime  = 1
	bufferLen = 1024
	syncInit  = iota
	syncDoing
)

func (self *LoggerAsync) init(tag string) {
	self.tag = tag
	self.data = make([]string, 0, bufferLen)
	self.lock = new(sync.RWMutex)
	self.status = syncInit
	timer := time.NewTicker(time.Second * syncTime)
	go func() {
		for {
			select {
			case <-timer.C:
				self.lock.RLock()
				if self.status != syncDoing {
					go self.flush()
				}
				self.lock.RUnlock()
			}
		}
	}()
}

func (self *LoggerAsync) write(str []byte) {
	self.lock.Lock()
	self.data = append(self.data, string(str))
	self.lock.Unlock()
}

func (self *LoggerAsync) openfile() {
	fileName := time.Now().Format(format_filerotate)
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

func (self *LoggerAsync) flush() error {
	if len(self.data) == 0 || self.status == syncDoing {
		return nil
	}
	self.status = syncDoing
	defer func() {
		self.status = syncInit
	}()

	self.lock.Lock()
	self.openfile()
	msg := self.data
	self.data = make([]string, 0, bufferLen)
	self.writer.WriteString(strings.Join(msg, ""))
	self.free()
	self.lock.Unlock()
	return nil
}

func (self *LoggerAsync) free() {
	if self.writer != nil {
		self.writer.Close()
		self.writer = nil
	}
}
