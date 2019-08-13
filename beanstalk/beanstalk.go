package beanstalk

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dirkli2011/utils/logkit"
	"github.com/prep/beanstalk"
)

var TTR_TIMEOUT = 5 * time.Second //处理超时时间
var DELAY_TIME = 3 * time.Second  // 处理失败则3秒后重新扔回队列，优先级置为最高

type JobExecuteor interface {
	Exec([]byte) bool
}

type beanstalkd struct {
	conn              string
	options           *beanstalk.Options
	Producer          *beanstalk.ProducerPool
	Consumer          *beanstalk.ConsumerPool
	ConsumerExecuteor JobExecuteor
}

func New(config string) (*beanstalkd, error) {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["conn"]; !ok {
		return nil, errors.New("config has no conn key")
	}
	options := &beanstalk.Options{
		ReserveTimeout:   time.Second,     // 等待超时
		ReconnectTimeout: time.Second * 3, // 重连超时
		ReadWriteTimeout: time.Second * 5, // 读写超时
		InfoLog:          log.New(os.Stdout, "Info", log.LstdFlags),
		ErrorLog:         log.New(os.Stderr, "Error", log.LstdFlags),
	}
	p, err := beanstalk.NewProducerPool([]string{cf["conn"]}, options)
	if err != nil {
		return nil, err
	}

	return &beanstalkd{
		conn:     cf["conn"],
		options:  options,
		Producer: p,
	}, nil
}

func (self *beanstalkd) Put(tube string, data []byte, priority uint32, args ...int) error {
	delay := 0
	if len(args) > 0 {
		delay = args[0]
	}
	putParams := &beanstalk.PutParams{priority, time.Duration(delay) * time.Second, TTR_TIMEOUT}
	id, err := self.Producer.Put(tube, data, putParams)
	if id > 0 {
		logkit.Info(fmt.Sprintf("Created job success with id: %d", id))
	} else {
		logkit.Info(fmt.Sprintf("Created job fail, errmsg=%s data=%s", err.Error(), string(data)))
	}
	return err
}

func (self *beanstalkd) Len(tube string) (int, error) {
	stats, err := beanstalk.TubeStats([]string{self.conn}, nil, tube)
	if err != nil {
		return 0, err
	}
	for _, stat := range stats {
		if stat.Name == tube {
			return stat.ReadyJobs + stat.DelayedJobs, nil
		}
	}
	return 0, nil
}

func (self *beanstalkd) Subscribe(tube string, obj JobExecuteor) *beanstalkd {
	self.Consumer, _ = beanstalk.NewConsumerPool([]string{self.conn}, []string{tube}, self.options)
	if obj == nil {
		panic("consumerExecuteor is nil")
	}
	self.ConsumerExecuteor = obj
	return self
}

func (self *beanstalkd) Wait() {
	defer self.Consumer.Stop()
	self.Consumer.Play()
	for {
		select {
		case job := <-self.Consumer.C:
			logkit.Info(fmt.Sprintf("Received job with id: %d", job.ID))
			if self.ConsumerExecuteor.Exec(job.Body) {
				logkit.Info(fmt.Sprintf("Finished job with id %d, data:%s", job.ID, string(job.Body)))
				job.Delete()
			} else {
				logkit.Warn(fmt.Sprintf("Failed job with id %d, data:%s", job.ID, string(job.Body)))
				job.ReleaseWithParams(0, DELAY_TIME)
			}
		}
	}
}
