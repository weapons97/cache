package cache

import (
	"sort"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/stretchr/testify/require"
)

func TestNewCache(t *testing.T) {
	c := NewCache(WithTTL[string, int](time.Second))
	b := 1
	c.Set(`a`, b)
	d, ok := c.Get(`a`)
	require.True(t, ok)
	require.Equal(t, b, d)
	time.Sleep(time.Second)
	d, ok = c.Get(`a`)
	require.False(t, ok)
	// 超时返回0值
	require.Equal(t, d, 0)
}

func TestNewCacheInits(t *testing.T) {
	inits := map[string]int{
		`a`: 1,
		`b`: 2,
	}
	c := NewCacheInits(inits)
	d, ok := c.Get(`a`)
	require.True(t, ok)
	require.Equal(t, 1, d)
	d, ok = c.Get(`b`)
	require.True(t, ok)
	require.Equal(t, 2, d)
}

type timeVal[T any] struct {
	t time.Time
	v T
}

func (tv *timeVal[T]) Time() time.Time {
	return tv.t
}

func TestNewCacheTimeout(t *testing.T) {
	c := NewCache(WithTTL[string, *timeVal[int]](time.Second))
	now := time.Now()
	b := &timeVal[int]{now.Add(-time.Second), 1}
	c.Set(`a`, b)
	d, ok := c.Get(`a`)
	require.False(t, ok)
	spew.Dump(d)
}

func TestCacheList(t *testing.T) {
	c := NewCache[string, int]()
	c.Set(`a`, 1)
	c.Set(`b`, 2)
	c.Set(`c`, 3)
	ks, vs := c.List()
	wantK := []string{`a`, `b`, `c`}
	wantV := []int{1, 2, 3}
	sort.Ints(vs)
	sort.Strings(ks)

	require.Equal(t, ks, wantK)
	require.Equal(t, vs, wantV)
	spew.Dump(ks, vs)
}

func TestCacheRange(t *testing.T) {
	c := NewCache[string, int]()
	c.Set(`a`, 1)
	c.Set(`b`, 2)
	c.Set(`c`, 3)
	ks := []string{}
	vs := []int{}
	c.Range(func(k string, v int) bool {
		ks = append(ks, k)
		vs = append(vs, v)
		return true
	})
	wantK := []string{`a`, `b`, `c`}
	wantV := []int{1, 2, 3}
	sort.Ints(vs)
	sort.Strings(ks)

	require.Equal(t, ks, wantK)
	require.Equal(t, vs, wantV)
	spew.Dump(ks, vs)
}

func TestCacheHasAny(t *testing.T) {
	c := NewCache[string, int]()

	// 初始状态，没有任何键
	require.False(t, c.HasAny("a", "b", "c"))

	// 添加一个键
	c.Set("a", 1)
	require.True(t, c.HasAny("a", "b", "c"))
	require.True(t, c.HasAny("a"))
	require.False(t, c.HasAny("b", "c"))

	// 添加更多键
	c.Set("b", 2)
	c.Set("c", 3)
	require.True(t, c.HasAny("a", "b", "c"))
	require.True(t, c.HasAny("a", "d", "e"))
	require.True(t, c.HasAny("d", "b", "e"))
	require.True(t, c.HasAny("d", "e", "c"))
	require.False(t, c.HasAny("d", "e", "f"))

	// 测试单个键
	require.True(t, c.HasAny("a"))
	require.False(t, c.HasAny("d"))

	// 测试空参数
	require.False(t, c.HasAny())
}
