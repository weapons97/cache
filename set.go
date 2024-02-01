// Package cache
/*=============================================================================
#       Author: peng.wei
#        Email: weapons97@gmail.com
#      Version: 0.0.1
#   LastChange: 20211214
#      History:
=============================================================================*/
package cache

// InterfaceSet 是set的接口
type InterfaceSet[K any] interface {
	Add(key ...K)
	Remove(key ...K)
	Pop() K
	Has(items ...K) bool
	Size() int
	Clear()
	IsEmpty() bool
	IsEqual(s *Set[K]) bool
	IsSubset(s *Set[K]) bool
	IsSuperset(s *Set[K]) bool
	Range(fn func(k K) bool)
	List() []K
	Copy() *Set[K]
	Merge(s *Set[K])
	Separate(t *Set[K])
}

// Set 是Set 型的cache 会多几个集合操作
type Set[K any] struct {
	inner *Cache[K, struct{}]
}

var setVal = struct{}{}

func newSet[K any](opts ...Option[K, struct{}]) *Set[K] {
	c := NewCache[K, struct{}](opts...)
	s := new(Set[K])
	s.inner = c
	return s
}

// NewSet 新创建set
func NewSet[K any](opts ...Option[K, struct{}]) *Set[K] {
	return newSet[K](opts...)
}

// NewSetInits 新创建set
func NewSetInits[K any](inits []K, opts ...Option[K, struct{}]) *Set[K] {
	s := newSet[K](opts...)
	for i := range inits {
		s.Add(inits[i])
	}
	return s
}

// Add 添加key
func (c *Set[K]) Add(key ...K) {
	for i := range key {
		c.inner.Set(key[i], setVal)
	}
}

// Has 已查找已传递的项目是否存在。如果未传递任何内容，则返回 false。对于多个项目，仅当所有项目都存在时，它才返回 true。
func (c *Set[K]) Has(items ...K) bool {
	if len(items) == 0 {
		return false
	}

	for _, item := range items {
		if _, ok := c.inner.Get(item); !ok {
			return false
		}
	}
	return true
}

func (c *Set[K]) Pop() (res K) {
	c.inner.Range(func(k K, v struct{}) bool {
		res = k
		c.inner.Del(k)
		return false
	})
	return res
}

// Remove 删除key
func (c *Set[K]) Remove(key ...K) {
	for i := range key {
		c.inner.Del(key[i])
	}
}

// Range 遍历set
func (c *Set[K]) Range(fn func(k K) bool) {
	c.inner.Range(func(k K, v struct{}) bool {
		return fn(k)
	})
}

// Size 返回cache 长度
func (c *Set[K]) Size() int {
	return c.inner.Len()
}

// Get 查找key是否存在
func (c *Set[K]) Get(key K) bool {
	_, ok := c.inner.Get(key)
	return ok
}

// List 列出元素
func (c *Set[K]) List() []K {
	res := make([]K, 0)
	c.Range(func(k K) bool {
		res = append(res, k)
		return true
	})
	return res
}

// ListStrings 以string类型列出元素
//func (c *Set[K]) ListStrings() []string {
//	res := make([]string, 0)
//	c.Range(func(k K) bool {
//		res = append(res, fmt.Sprintf(`%v`, k))
//		return true
//	})
//	return res
//}

// Union 并集
func (c *Set[K]) Union(s *Set[K], opts ...Option[K, struct{}]) *Set[K] {
	n := NewSet[K](opts...)
	c.Range(func(k K) bool {
		n.Add(k)
		return true
	})
	s.Range(func(k K) bool {
		n.Add(k)
		return true
	})
	return n
}

// Difference godoc
func (c *Set[K]) Difference(s *Set[K], opts ...Option[K, struct{}]) *Set[K] {
	n := NewSet[K](opts...)
	c.Range(func(k K) bool {
		n.Add(k)
		return true
	})
	s.Range(func(k K) bool {
		ok := n.Get(k)
		if ok {
			n.Remove(k)
		}
		return true
	})
	return n
}

// Intersection 交集
func (c *Set[K]) Intersection(s *Set[K], opts ...Option[K, struct{}]) *Set[K] {
	n := NewSet[K](opts...)
	c.Range(func(k K) bool {
		n.Add(k)
		return true
	})
	n.Range(func(k K) bool {
		ok := s.Get(k)
		if !ok {
			n.Remove(k)
		}
		return true
	})
	return n
}

// Sub 差集
func (c *Set[K]) Sub(s *Set[K], opts ...Option[K, struct{}]) *Set[K] {
	n := NewSet[K](opts...)
	c.Range(func(k K) bool {
		n.Add(k)
		return true
	})
	s.Range(func(k K) bool {
		ok := n.Get(k)
		if ok {
			n.Remove(k)
		}
		return true
	})
	return n
}

// Separate it's not the opposite of Merge.
// Separate removes the set items containing in t from set s. Please aware that
func (c *Set[K]) Separate(t *Set[K]) {
	c.Remove(t.List()...)
}

// Copy 复制set
func (c *Set[K]) Copy() *Set[K] {
	opts := c.inner.opts
	return NewSetInits(c.List(), opts...)
}

// Merge 合并set
func (c *Set[K]) Merge(s *Set[K]) {
	s.Range(func(k K) bool {
		c.Add(k)
		return true
	})
}

// Clear 清空set
func (c *Set[K]) Clear() {
	c.inner = NewCache[K, struct{}]()
}

// IsEmpty 判断set是否为空
func (c *Set[K]) IsEmpty() bool {
	return c.Size() == 0
}

// IsEqual 判断set是否相等
func (c *Set[K]) IsEqual(s *Set[K]) bool {
	if c.Size() != s.Size() {
		return false
	}
	return c.Difference(s).IsEmpty()
}

// IsSubset 判断set是否为子集
func (c *Set[K]) IsSubset(s *Set[K]) (subset bool) {
	subset = true
	s.Range(func(k K) bool {
		subset = c.Get(k)
		return subset
	})
	return
}

// IsSuperset 判断set是否为父集
func (c *Set[K]) IsSuperset(s *Set[K]) bool {
	return s.IsSubset(c)
}
