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
type IndexFunc func(indexed any) (keys []string)

// Indexed 接口, 一个结构实现了Indexed 接口才可以被Indexer 使用
type Indexed interface {
	Indexes() map[string]IndexFunc
	ID() (mainKey string)
}

func IndexGet[T any](i Indexed) (rx T, ok bool) {
	rx, ok = i.(T)
	if !ok {
		return rx, false
	}
	return rx, true
}

// Indexer 是带索引的cache
type Indexer[T Indexed] struct {
	cs   map[string]*Cache[string, *Set[string]] // 索引表
	rw   sync.RWMutex
	main *Cache[string, T] // 主表
	opts []Option[string, T]
}

// NewIndexer 创建一个带索引的cache
func NewIndexer[T Indexed](ops ...Option[string, T]) *Indexer[T] {
	ix := new(Indexer[T])
	ix.cs = make(map[string]*Cache[string, *Set[string]])
	ix.rw = sync.RWMutex{}
	ix.opts = ops
	ix.main = NewCache[string, T](ops...)
	return ix
}

// Set 设置值，v 必须和 Indexer 的type相同
func (ix *Indexer[T]) Set(v T) bool {
	id := v.ID()
	if ix.main == nil {
		ix.main = NewCache[string, T](ix.opts...)
	}
	ix.Del(id)
	ix.main.Set(id, v)
	idxs := v.Indexes()
	for name, idx := range idxs {
		keys := idx(v)
		ix.rw.Lock()
		if _, ok := ix.cs[name]; !ok {
			ix.cs[name] = NewCache[string, *Set[string]]()
		}
		c := ix.cs[name]
		ix.rw.Unlock()
		for _, key := range keys {
			set, ok := c.Get(key)
			// set, ok2 := v.(*Set)
			if !ok {
				set = NewSet[string]()
			}
			set.Add(id)
			c.Set(key, set)
		}
	}
	return true
}

// Len 返回cache 长度
func (ix *Indexer[T]) Len() int {
	i := 0
	ix.Range(func(k string, v T) bool {
		i++
		return true
	})
	return i
}

// Get 根据id 查找Indexed
func (ix *Indexer[T]) Get(id string) (v T, ok bool) {
	rx, ok := ix.main.Get(id)
	if !ok {
		return v, ok
	}
	return rx, true
}

// Del 删除一个Indexed
func (ix *Indexer[T]) Del(v interface{}) {
	sv, ok := v.(string)
	if ok {
		v2, ok := ix.Get(sv)
		if !ok {
			return
		}
		ix.del(v2)
	}
	v2, ok := v.(T)
	if ok {
		ix.del(v2)
	}
}

func (ix *Indexer[T]) del(req Indexed) {
	id := req.ID()
	if ix.main == nil {
		return
	}
	c := ix.main
	c.Del(id)
	idxs := req.Indexes()
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
			set, ok := c.Get(key)
			// set, ok2 := v.(*Set)
			if ok {
				set.Remove(id)
			}
		}
	}
}

// Range 遍历Indexer
func (ix *Indexer[T]) Range(fn func(k string, v T) bool) {
	ix.main.Range(fn)
}

// SetFromIndex 从indexName 创建一个Set
func (ix *Indexer[T]) SetFromIndex(idxName string) (*Set[string], error) {
	ix.rw.RLock()
	c, ok := ix.cs[idxName]
	ix.rw.RUnlock()
	if !ok {
		return nil, fmt.Errorf(`no such index`)
	}
	keys := make([]string, 0)
	c.Range(func(k string, _ *Set[string]) bool {
		keys = append(keys, k)
		return true
	})
	return NewSetInits[string](keys), nil
}

// SearchResult 是Indexer 根据索引函数查找的结果
type SearchResult[T Indexed] struct {
	e   error
	Res []T
}

// Error 查找的错误
func (sr *SearchResult[T]) Error() error {
	if sr == nil {
		return fmt.Errorf(`SearchResult is nil`)
	}
	return sr.e
}

// Failed 查找是否成功
func (sr *SearchResult[T]) Failed() bool {
	if sr.Error() != nil {
		return true
	}
	return false
}

// InvokeOne 拿一个结果就好
func (sr *SearchResult[T]) InvokeOne() (rx T) {
	if sr.e == nil && len(sr.Res) > 0 {
		return sr.Res[0]
	}
	return
}

// InvokeAll 返回所有搜索结果
func (sr *SearchResult[T]) InvokeAll() []T {
	if sr.e == nil {
		return sr.Res
	}
	return nil
}

// Range 遍历所有搜索结果
func (sr *SearchResult[T]) Range(fn func(v T) bool) {
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
func (ix *Indexer[T]) Search(idxName string, key string) *SearchResult[T] {
	vs, e := ix.search(idxName, key)
	return &SearchResult[T]{
		e:   e,
		Res: vs,
	}
}

func (ix *Indexer[T]) search(idxName string, key string) (vs []T, e error) {
	ix.rw.RLock()
	c, ok := ix.cs[idxName]
	ix.rw.RUnlock()
	if !ok {
		return nil, fmt.Errorf(`no such index`)
	}
	idSet, ok := c.Get(key)
	if !ok {
		return nil, fmt.Errorf(`index %v no such key %v`, idxName, key)
	}

	ids := idSet.List()

	vs = make([]T, 0, len(ids))
	for i := range ids {
		res, ok := ix.main.Get(ids[i])
		v, ok := IndexGet[T](res)
		if !ok {
			return nil, fmt.Errorf(`search index id %T can't get value'`, res)
		}
		vs = append(vs, v)
	}

	return vs, nil
}
