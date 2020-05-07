package cmap

import (
	"encoding/json"
	"sync"
)

const BUCKET_COUNT = 32

type ConcurrentMap []*bucketMap

type bucketMap struct {
	items map[string]interface{}
	sync.RWMutex
}

func New() ConcurrentMap {
	m := make(ConcurrentMap, BUCKET_COUNT)
	for i := 0; i < BUCKET_COUNT; i++ {
		m[i] = &bucketMap{items: make(map[string]interface{})}
	}
	return m
}

func (m ConcurrentMap) getBucketMap(key string) *bucketMap {
	return m[bucketNum(key)]
}

func (m ConcurrentMap) Set(key string, value interface{}) {
	b := m.getBucketMap(key)
	b.Lock()
	b.items[key] = value
	b.Unlock()
}

func (m ConcurrentMap) MultiSet(data map[string]interface{}) {
	for key, value := range data {
		b := m.getBucketMap(key)
		b.Lock()
		b.items[key] = value
		b.Unlock()
	}
}

// 不存在时才设置, 返回false表示key存在
func (m ConcurrentMap) SetIfNot(key string, value interface{}) bool {
	b := m.getBucketMap(key)
	b.Lock()
	_, ok := b.items[key]
	if !ok {
		b.items[key] = value
	}
	b.Unlock()
	return !ok
}

func (m ConcurrentMap) Get(key string) (interface{}, bool) {
	b := m.getBucketMap(key)
	b.Lock()
	value, ok := b.items[key]
	b.Unlock()
	return value, ok
}

func (m ConcurrentMap) Len() (count int) {
	for i := 0; i < BUCKET_COUNT; i++ {
		b := m[i]
		b.RLock()
		count += len(b.items)
		b.RUnlock()
	}
	return
}

func (m ConcurrentMap) Exist(key string) bool {
	b := m.getBucketMap(key)
	b.RLock()
	_, ok := b.items[key]
	b.RUnlock()
	return ok
}

func (m ConcurrentMap) Delete(key string) {
	b := m.getBucketMap(key)
	b.Lock()
	delete(b.items, key)
	b.Unlock()
}

// 获取并删除
func (m ConcurrentMap) Pop(key string) (interface{}, bool) {
	b := m.getBucketMap(key)
	b.Lock()
	value, ok := b.items[key]
	delete(b.items, key)
	b.Unlock()
	return value, ok
}

func (m ConcurrentMap) Empty() bool {
	return m.Len() == 0
}

type Item struct {
	Key   string
	Value interface{}
}

func (m ConcurrentMap) Iterator() <-chan Item {
	ch := make(chan Item, m.Len())
	var wg sync.WaitGroup
	wg.Add(BUCKET_COUNT)
	for i := 0; i < BUCKET_COUNT; i++ {
		go func(b *bucketMap) {
			defer wg.Done()
			b.RLock()
			for k, v := range b.items {
				ch <- Item{k, v}
			}
			b.RUnlock()
		}(m[i])
	}
	wg.Wait()
	close(ch)
	return ch
}

type IterCallback func(key string, v interface{})

// 对每个item执行IterCallback操作
func (m ConcurrentMap) IteratorBy(fn IterCallback) {
	for i := 0; i < BUCKET_COUNT; i++ {
		b := m[i]
		b.RLock()
		for key, value := range b.items {
			fn(key, value)
		}
		b.RUnlock()
	}
}

type FilterCallback func(key string, v interface{}) bool

// FilterCallback为false的item会被过滤掉
func (m ConcurrentMap) FilterBy(fn FilterCallback) map[string]interface{} {
	rst := make(map[string]interface{})
	for item := range m.Iterator() {
		if fn(item.Key, item.Value) {
			rst[item.Key] = item.Value
		}
	}
	return rst
}

func (m ConcurrentMap) Items() map[string]interface{} {
	rst := make(map[string]interface{})
	for item := range m.Iterator() {
		rst[item.Key] = item.Value
	}
	return rst
}

func (m ConcurrentMap) Keys() []string {
	keys := make([]string, 0, m.Len())
	for item := range m.Iterator() {
		keys = append(keys, item.Key)
	}
	return keys
}

func (m ConcurrentMap) JsonString() string {
	b, err := json.Marshal(m.Items())
	if err == nil {
		return string(b)
	}
	return "{}"
}

func (m ConcurrentMap) Json() ([]byte, error) {
	return json.Marshal(m.Items())
}

func bucketNum(key string) uint {
	return uint(fnv0(key)) % BUCKET_COUNT
}

func fnv0(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}
