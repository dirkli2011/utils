package utils

import (
	"github.com/cornelk/hashmap"
)

// 返回一个hashmap
func NewHashMap() *hashmap.HashMap {
	return &hashmap.HashMap{}
}

// 返回一个双向链表
func NewList() *hashmap.List {
	return hashmap.NewList()
}
