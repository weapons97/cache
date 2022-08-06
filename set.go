// Package cache
/*=============================================================================
#       Author: peng.wei
#        Email: weapons97@gmail.com
#      Version: 0.0.1
#   LastChange: 20211214
#      History:
=============================================================================*/
package cache

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
func (c *Set[K]) Add(key K) {
	c.inner.Set(key, setVal)
}

// Del 删除key
func (c *Set[K]) Del(key K) {
	c.inner.Del(key)
}

// Range 遍历set
func (c *Set[K]) Range(fn func(k K) bool) {
	c.inner.Range(func(k K, v struct{}) bool {
		return fn(k)
	})
}

// Len 返回cache 长度
func (c *Set[K]) Len() int {
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

// JoinLeft 左交集
func (c *Set[K]) JoinLeft(s *Set[K], opts ...Option[K, struct{}]) *Set[K] {
	n := NewSet[K](opts...)
	c.Range(func(k K) bool {
		n.Add(k)
		return true
	})
	s.Range(func(k K) bool {
		ok := n.Get(k)
		if !ok {
			n.Del(k)
		}
		return true
	})
	return n
}

// JoinRight 右交集
func (c *Set[K]) JoinRight(s *Set[K], opts ...Option[K, struct{}]) *Set[K] {
	n := NewSet[K](opts...)
	s.Range(func(k K) bool {
		n.Add(k)
		return true
	})
	c.Range(func(k K) bool {
		ok := n.Get(k)
		if !ok {
			n.Del(k)
		}
		return true
	})
	return n
}

// Join 交集
func (c *Set[K]) Join(s *Set[K], opts ...Option[K, struct{}]) *Set[K] {
	n := NewSet[K](opts...)
	c.Range(func(k K) bool {
		n.Add(k)
		return true
	})
	n.Range(func(k K) bool {
		ok := s.Get(k)
		if !ok {
			n.Del(k)
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
			n.Del(k)
		}
		return true
	})
	return n
}
