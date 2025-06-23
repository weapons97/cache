// Package cache
/*=============================================================================
#       Author: peng.wei
#        Email: weapons97@gmail.com
#      Version: 0.0.1
#   LastChange: 20211214
#      History:
=============================================================================*/
package cache

// PowerSet 求幂集
func PowerSet[T any](s []T) [][]T {
	n := len(s)
	powerset := [][]T{}
	for mask := 1; mask < (1 << n); mask++ {
		subSet := []T{}
		for j := 0; j < n; j++ {
			if (mask>>j)&1 == 1 {
				subSet = append(subSet, s[j])
			}
		}
		powerset = append(powerset, subSet)
	}
	return powerset
}

// InterfaceSet 是set的接口
type InterfaceSet[K comparable] interface {
	Add(key ...K)
	Remove(key ...K)
	Pop() K
	Has(items ...K) bool
	HasAny(k ...K) bool
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
type Set[K comparable] struct {
	inner *Cache[K, struct{}]
}

var setVal = struct{}{}

func newSet[K comparable](opts ...Option[K, struct{}]) *Set[K] {
	c := NewCache[K, struct{}](opts...)
	s := new(Set[K])
	s.inner = c
	return s
}

// NewSet 新创建set
func NewSet[K comparable](opts ...Option[K, struct{}]) *Set[K] {
	return newSet[K](opts...)
}

// NewSetInits 新创建set
func NewSetInits[K comparable](inits []K, opts ...Option[K, struct{}]) *Set[K] {
	s := newSet[K](opts...)
	for i := range inits {
		s.Add(inits[i])
	}
	return s
}

// Add 添加key
func (s *Set[K]) Add(key ...K) {
	for i := range key {
		s.inner.Set(key[i], setVal)
	}
}

// Has 已查找已传递的项目是否存在。如果未传递任何内容，则返回 false。对于多个项目，仅当所有项目都存在时，它才返回 true。
func (s *Set[K]) Has(items ...K) bool {
	if len(items) == 0 {
		return false
	}

	for _, item := range items {
		if _, ok := s.inner.Get(item); !ok {
			return false
		}
	}
	return true
}

// HasAny 检查是否存在任何一个传递的项目。如果未传递任何内容，则返回 false。对于多个项目，只要有一个存在就返回 true。
func (s *Set[K]) HasAny(items ...K) bool {
	if len(items) == 0 {
		return false
	}

	for _, item := range items {
		if _, ok := s.inner.Get(item); ok {
			return true
		}
	}
	return false
}

func (s *Set[K]) Pop() (res K) {
	s.inner.Range(func(k K, v struct{}) bool {
		res = k
		s.inner.Del(k)
		return false
	})
	return res
}

// Remove 删除key
func (s *Set[K]) Remove(key ...K) {
	for i := range key {
		s.inner.Del(key[i])
	}
}

// Range 遍历set
func (s *Set[K]) Range(fn func(k K) bool) {
	s.inner.Range(func(k K, v struct{}) bool {
		return fn(k)
	})
}

// Size 返回cache 长度
func (s *Set[K]) Size() int {
	return s.inner.Len()
}

// Get 查找key是否存在
func (s *Set[K]) Get(key K) bool {
	_, ok := s.inner.Get(key)
	return ok
}

// List 列出元素
func (s *Set[K]) List() []K {
	res := make([]K, 0)
	s.Range(func(k K) bool {
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
func (s *Set[K]) Union(o *Set[K], opts ...Option[K, struct{}]) *Set[K] {
	n := NewSet[K](opts...)
	s.Range(func(k K) bool {
		n.Add(k)
		return true
	})
	o.Range(func(k K) bool {
		n.Add(k)
		return true
	})
	return n
}

// Difference godoc
func (s *Set[K]) Difference(o *Set[K], opts ...Option[K, struct{}]) *Set[K] {
	n := NewSet[K](opts...)
	s.Range(func(k K) bool {
		n.Add(k)
		return true
	})
	o.Range(func(k K) bool {
		ok := n.Get(k)
		if ok {
			n.Remove(k)
		}
		return true
	})
	return n
}

// Intersection 交集
func (s *Set[K]) Intersection(o *Set[K], opts ...Option[K, struct{}]) *Set[K] {
	n := NewSet[K](opts...)
	s.Range(func(k K) bool {
		n.Add(k)
		return true
	})
	n.Range(func(k K) bool {
		ok := o.Get(k)
		if !ok {
			n.Remove(k)
		}
		return true
	})
	return n
}

// Sub 差集
func (s *Set[K]) Sub(o *Set[K], opts ...Option[K, struct{}]) *Set[K] {
	n := NewSet[K](opts...)
	s.Range(func(k K) bool {
		n.Add(k)
		return true
	})
	o.Range(func(k K) bool {
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
func (s *Set[K]) Separate(t *Set[K]) {
	s.Remove(t.List()...)
}

// Copy 复制set
func (s *Set[K]) Copy() *Set[K] {
	opts := s.inner.opts
	return NewSetInits(s.List(), opts...)
}

// Merge 合并set
func (s *Set[K]) Merge(o *Set[K]) {
	o.Range(func(k K) bool {
		s.Add(k)
		return true
	})
}

// Clear 清空set
func (s *Set[K]) Clear() {
	s.inner = NewCache[K, struct{}]()
}

// IsEmpty 判断set是否为空
func (s *Set[K]) IsEmpty() bool {
	return s.Size() == 0
}

// IsEqual 判断set是否相等
func (s *Set[K]) IsEqual(o *Set[K]) bool {
	if s.Size() != o.Size() {
		return false
	}
	return s.Difference(o).IsEmpty()
}

// IsSubset 判断set是否为子集
func (s *Set[K]) IsSubset(o *Set[K]) (subset bool) {
	subset = true
	o.Range(func(k K) bool {
		subset = s.Get(k)
		return subset
	})
	return
}

// IsSuperset 判断set是否为父集
func (s *Set[K]) IsSuperset(o *Set[K]) bool {
	return o.IsSubset(s)
}

// PowerSet 取幂集
func (s *Set[K]) PowerSet(key ...K) (powerSet []*Set[K]) {
	n := s.Size()
	l := s.List()
	emptySet := NewSet[K]()
	powerSet = []*Set[K]{
		emptySet,
	}
	for mask := 1; mask < (1 << n); mask++ {
		subSet := NewSet[K]()
		for j := 0; j < n; j++ {
			if (mask>>j)&1 == 1 {
				subSet.Add(l[j])
			}
		}
		powerSet = append(powerSet, subSet)
	}
	return powerSet
}
