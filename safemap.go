// 多线程安全的map，对性能有损耗
package utils

import (
	"sync"
)

type SafeMap struct {
	lock *sync.RWMutex
	m    map[interface{}]interface{}
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		lock: new(sync.RWMutex),
		m:    make(map[interface{}]interface{}),
	}
}

func (this *SafeMap) Get(k interface{}) interface{} {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if val, ok := this.m[k]; ok {
		return val
	}
	return nil
}

func (this *SafeMap) Set(k interface{}, v interface{}) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	if val, ok := this.m[k]; !ok {
		this.m[k] = v
	} else if val != v {
		this.m[k] = v
	} else {
		return false
	}
	return true
}

func (this *SafeMap) Check(k interface{}) bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	_, ok := this.m[k]
	return ok
}

func (this *SafeMap) Delete(k interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.m, k)
}

func (this *SafeMap) Count() int {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return len(this.m)
}
