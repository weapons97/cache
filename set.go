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
type Set struct {
	inner *Cache
}

var setVal = struct{}{}

func newSet(opts ...Option) *Set {
	c := NewCache(opts...)
	s := new(Set)
	s.inner = c
	return s
}

// NewSet 新创建set
func NewSet(opts ...Option) *Set {
	return newSet(opts...)
}

// Add 添加key
func (c *Set) Add(key interface{}) {
	c.inner.Set(key, setVal)
}

// Del 删除key
func (c *Set) Del(key interface{}) {
	c.inner.Del(key)
}

// Range 遍历set
func (c *Set) Range(fn func(k interface{}) bool) {
	c.inner.Range(func(k, v interface{}) bool {
		return fn(k)
	})
}

// Get 查找key是否存在
func (c *Set) Get(key interface{}) bool {
	_, ok := c.inner.Get(key)
	return ok
}

// List 列出元素
func (c *Set) List() []interface{} {
	res := make([]interface{}, 0)
	c.Range(func(k interface{}) bool {
		res = append(res, k)
		return true
	})
	return res
}

// ListStrings 以string类型列出元素
func (c *Set) ListStrings() []string {
	res := make([]string, 0)
	c.Range(func(k interface{}) bool {
		s, ok := k.(string)
		if ok {
			res = append(res, s)
		}
		return true
	})
	return res
}

// Union 并集
func (c *Set) Union(s *Set, opts ...Option) *Set {
	n := NewSet(opts...)
	c.Range(func(k interface{}) bool {
		n.Add(k)
		return true
	})
	s.Range(func(k interface{}) bool {
		n.Add(k)
		return true
	})
	return n
}

// JoinLeft 左交集
func (c *Set) JoinLeft(s *Set, opts ...Option) *Set {
	n := NewSet(opts...)
	c.Range(func(k interface{}) bool {
		n.Add(k)
		return true
	})
	s.Range(func(k interface{}) bool {
		ok := n.Get(k)
		if !ok {
			n.Del(k)
		}
		return true
	})
	return n
}

// JoinRight 右交集
func (c *Set) JoinRight(s *Set, opts ...Option) *Set {
	n := NewSet(opts...)
	s.Range(func(k interface{}) bool {
		n.Add(k)
		return true
	})
	c.Range(func(k interface{}) bool {
		ok := n.Get(k)
		if !ok {
			n.Del(k)
		}
		return true
	})
	return n
}

// Join 交集
func (c *Set) Join(s *Set, opts ...Option) *Set {
	n := NewSet(opts...)
	c.Range(func(k interface{}) bool {
		n.Add(k)
		return true
	})
	n.Range(func(k interface{}) bool {
		ok := s.Get(k)
		if !ok {
			n.Del(k)
		}
		return true
	})
	return n
}

// Sub 差集
func (c *Set) Sub(s *Set, opts ...Option) *Set {
	n := NewSet(opts...)
	c.Range(func(k interface{}) bool {
		n.Add(k)
		return true
	})
	s.Range(func(k interface{}) bool {
		ok := n.Get(k)
		if ok {
			n.Del(k)
		}
		return true
	})
	return n
}
