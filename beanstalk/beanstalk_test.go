package beanstalk

import (
	"testing"
)

type JobExec struct {
}

// 这里注意处理超时,及时返回
func (self JobExec) Exec(data []byte) bool {
	return true
}

func TestBeanstalk(t *testing.T) {
	conn, err := New(`{"conn": "beanstalk://127.0.0.1:11380"}`)
	if err != nil {
		t.Error("new beanstalk err", err)
	}

	err = conn.Put("test", []byte("Hello World"), 1024)
	if err != nil {
		t.Error("Put err", err)
	}

	err = conn.Put("test", []byte("Hello World"), 1024, 1)
	if err != nil {
		t.Error("PutDelay err", err)
	}

	n, err := conn.Len("test")
	if err != nil {
		t.Error("stat err", err)
	}
	if n != 2 {
		t.Error("Len err", n)
	}

	conn.Subscribe("test", &JobExec{}).Wait()

	n, err = conn.Len("test")
	if n != 0 {

		t.Error("Len err", n)
	}
}
