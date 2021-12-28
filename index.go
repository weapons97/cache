// Package cache
/*=============================================================================
#       Author: peng.wei
#        Email: weapons97@gmail.com
#      Version: 0.0.1
#   LastChange: 20211214
#      History:
=============================================================================*/
package cache

import (
	"fmt"
	"sync"
)

// IndexFunc 是Indexed元素的索引函数
type IndexFunc func(indexed Indexed) (keys []string)

// Indexed 接口, 一个结构实现了Indexed 接口才可以被Indexer 使用
type Indexed interface {
	Indexs() map[string]IndexFunc
	Id() (mainKey string)
	Set(v interface{}) (Indexed, bool)
	Get(v Indexed) (interface{}, bool)
}

// Indexer 是带索引的cache
type Indexer struct {
	cs    map[string]*Cache // 索引表
	rw    sync.RWMutex
	main  *Cache  // 主表
	typed Indexed // Indexer 元素类型
	opts  []Option
}

// NewIndexer 创建一个带索引的cache
func NewIndexer(typed Indexed, ops ...Option) *Indexer {
	ix := new(Indexer)
	ix.cs = make(map[string]*Cache)
	ix.rw = sync.RWMutex{}
	ix.opts = ops
	ix.typed = typed
	ix.main = NewCache(ops...)
	return ix
}

// Type 返回Indexer 索引的类型 Indexed
func (ix *Indexer) Type() Indexed {
	return ix.typed
}

// Set 设置值，v 必须和 Indexer 的type相同
func (ix *Indexer) Set(v interface{}) bool {

	req, ok := ix.Type().Set(v)
	if !ok {
		return false
	}
	id := req.Id()
	if ix.main == nil {
		ix.main = NewCache(ix.opts...)
	}
	_, ok = ix.Get(id)
	if ok {
		ix.Del(req)
	}
	c := ix.main
	c.Set(id, req)
	idxs := req.Indexs()
	for name, idx := range idxs {
		keys := idx(req)
		ix.rw.Lock()
		if _, ok := ix.cs[name]; !ok {
			ix.cs[name] = NewCache()
		}
		c := ix.cs[name]
		ix.rw.Unlock()
		for _, key := range keys {
			v, ok := c.Get(key)
			set, ok2 := v.(*Set)
			if !ok || !ok2 {
				set = NewSet()
			}
			set.Add(id)
			c.Set(key, set)
		}
	}
	return true
}

// Get 根据id 查找Indexed
func (ix *Indexer) Get(id string) (v interface{}, ok bool) {
	rx, ok := ix.main.Get(id)
	if !ok {
		return nil, false
	}
	var res Indexed
	res, ok = rx.(Indexed)
	if !ok {
		return nil, false
	}
	return ix.typed.Get(res)
}

// Del 删除一个Indexed
func (ix *Indexer) Del(req Indexed) {
	id := req.Id()
	if ix.main == nil {
		return
	}
	c := ix.main
	c.Del(id)
	idxs := req.Indexs()
	for name, idx := range idxs {
		keys := idx(req)
		ix.rw.RLock()
		if _, ok := ix.cs[name]; !ok {
			ix.rw.RUnlock()
			continue
		}
		ix.rw.RUnlock()
		c := ix.cs[name]
		for _, key := range keys {
			v, ok := c.Get(key)
			set, ok2 := v.(*Set)
			if ok && ok2 {
				set.Del(id)
			}
		}
	}
}

// Range 遍历Indexer
func (ix *Indexer) Range(fn func(k, v interface{}) bool) {
	ix.main.smap.Range(func(k, v interface{}) bool {
		xv, ok := ix.main.unWrapTTL(v)
		if !ok {
			return true
		}
		iv, ok := xv.(Indexed)
		if !ok {
			return true
		}
		v, ok = ix.typed.Get(iv)
		if !ok {
			return true
		}
		fn(k, v)
		return true
	})
}

// SearchResult 是Indexer 根据索引函数查找的结果
type SearchResult struct {
	e   error
	Res []interface{}
}

// Error 查找的错误
func (sr *SearchResult) Error() error {
	if sr == nil {
		return fmt.Errorf(`SearchResult is nil`)
	}
	return sr.e
}

// Failed 查找是否成功
func (sr *SearchResult) Failed() bool {
	if sr.Error() != nil {
		return true
	}
	return false
}

// InvokeOne 拿一个结果就好
func (sr *SearchResult) InvokeOne() interface{} {
	if sr.e == nil && len(sr.Res) > 0 {
		return sr.Res[0]
	}
	return nil
}

// InvokeAll 返回所有搜索结果
func (sr *SearchResult) InvokeAll() []interface{} {
	if sr.e == nil {
		return sr.Res
	}
	return nil
}

// Range 遍历所有搜索结果
func (sr *SearchResult) Range(fn func(v interface{}) bool) {
	if sr == nil || sr.e != nil {
		return
	}
	for i := range sr.Res {
		conti := fn(sr.Res[i])
		if !conti {
			return
		}
	}
}

// Search 根据索引函数查找Indexer
func (ix *Indexer) Search(idxName string, key string) *SearchResult {
	vs, e := ix.search(idxName, key)
	return &SearchResult{
		e:   e,
		Res: vs,
	}
}

func (ix *Indexer) search(idxName string, key string) (vs []interface{}, e error) {
	ix.rw.RLock()
	c, ok := ix.cs[idxName]
	ix.rw.RUnlock()
	if !ok {
		return nil, fmt.Errorf(`no such index`)
	}
	rx, ok := c.Get(key)
	if !ok {
		return nil, fmt.Errorf(`index %v no such key %v`, idxName, key)
	}

	idSet, ok := rx.(*Set) // id table must set
	if !ok {
		return nil, fmt.Errorf(`index %v key %v not idSet`, idxName, key)
	}
	ids := idSet.ListStrings()

	vs = make([]interface{}, 0, len(ids))
	for i := range ids {
		rx, ok = ix.main.Get(ids[i])
		var res Indexed
		res, ok = rx.(Indexed)
		if !ok {
			return nil, fmt.Errorf(`search index id %v not main value`, ids[i])
		}
		v, ok := ix.typed.Get(res)
		if !ok {
			return nil, fmt.Errorf(`search index id %T can't get value'`, ix.typed)
		}
		vs = append(vs, v)
	}

	return vs, nil
}
