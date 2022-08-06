package cache

import (
	"testing"
	"time"

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
