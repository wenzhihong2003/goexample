package main

import (
	"sync"
	"sort"
)

// [Go小技巧] 实现常用的KV缓存（有序且并发安全）
// https://my.oschina.net/henrylee2cn/blog/741315

type kv struct {
	count int
	keys  []string
	hash  map[string]interface{}
	lock  sync.RWMutex
}

// 新建kv缓存(preCapacity为预申请内存容量)
func NewKv(preCapacity uint) *kv {
	return &kv{
		keys: make([]string, 0, int(preCapacity)),
		hash: make(map[string]interface{}, int(preCapacity)),
	}
}

// 添加kv键值对
func (this *kv) Set(k string, v interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if _, ok := this.hash[k]; !ok {
		this.keys = append(this.keys, k)
		sort.Strings(this.keys)
		this.count++
	}
	this.hash[k] = v
}

// 获取长度
func (this *kv) Count() int {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.count
}

// 由key检索value
func (this *kv) Get(k string) (interface{}, bool) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	v, ok := this.hash[k]
	return v, ok
}

// 根据key排序, 返回有序的value切片
func (this *kv) Values() []interface{} {
	this.lock.RLock()
	defer this.lock.RUnlock()
	vals := make([]interface{}, this.count)
	for i := 0; i < this.count; i++ {
		vals[i] = this.hash[this.keys[i]]
	}
	return vals
}

func main() {

}