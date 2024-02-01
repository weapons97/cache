package cache

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
)

func TestSetUnion(t *testing.T) {
	s := NewSet[string]()
	s.Add(`a`)
	s.Add(`b`)
	s2 := NewSet[string]()
	s2.Add(`b`)
	s2.Add(`d`)
	s3 := s.Union(s2)
	wantS3 := []string{`a`, `b`, `d`}

	ans := s3.List()
	sort.Strings(ans)
	require.Equal(t, wantS3, ans)
	spew.Dump(s.List(), s2.List(), s3.List())
}

func TestIntersection(t *testing.T) {
	s := NewSet[string]()
	s.Add(`a`)
	s.Add(`b`)
	s2 := NewSet[string]()
	s2.Add(`b`)
	s2.Add(`d`)
	s3 := s.Intersection(s2)
	wantS3 := []string{`b`}

	ans := s3.List()
	sort.Strings(ans)
	require.Equal(t, wantS3, ans)
	spew.Dump(s.List(), s2.List(), s3.List())
}

func TestDifference(t *testing.T) {
	s := NewSetInits([]int{1, 2, 3, 8})

	s2 := NewSetInits([]int{2, 3, 4})

	s3 := s.Difference(s2)
	wantS3 := []int{1, 8}
	ans := s3.List()
	sort.Ints(ans)

	require.Equal(t, wantS3, ans)
	spew.Dump(s.List(), s2.List(), s3.List())
}

func TestIsEqual(t *testing.T) {
	s := NewSetInits([]int{1, 2, 3, 8})
	s2 := NewSetInits([]int{2, 3, 4})
	s3 := NewSetInits([]int{4, 2, 3})

	require.False(t, s.IsEqual(s2))
	require.True(t, s2.IsEqual(s3))
}

func TestHas(t *testing.T) {
	s := NewSetInits([]int{1, 2, 3, 8})
	s2 := NewSetInits([]int{2, 3})
	s3 := NewSetInits([]int{4, 2})

	require.True(t, s.Has(s2.List()...))
	require.False(t, s.Has(s3.List()...))
}

func TestSetSub(t *testing.T) {
	s := NewSet[string]()
	s.Add(`a`)
	s.Add(`b`)
	s2 := NewSet[string]()
	s2.Add(`b`)
	s2.Add(`d`)
	s3 := s.Sub(s2)
	wantS3 := []string{`a`}
	ans := s3.List()
	sort.Strings(ans)
	require.Equal(t, wantS3, ans)
	spew.Dump(s.List(), s2.List(), s3.List())
}

func TestSetPop(t *testing.T) {
	req := []string{`a`, `b`}
	s := NewSetInits(req)
	ans := []string{}
	for i := 0; i < 10; i++ {
		x := s.Pop()
		if x != "" {
			ans = append(ans, x)
		}
	}
	sort.Strings(ans)
	require.Equal(t, req, ans)
	spew.Dump(ans)
}

func TestSetHas(t *testing.T) {
	req := []string{`a`, `b`}
	s := NewSetInits(req)
	s.Has(`a`)
	require.True(t, s.Has(`a`))
	require.True(t, s.Has(`b`))
	require.True(t, s.Has(`a`, `b`))
	require.False(t, s.Has(`a`, `b`, `c`))
	require.False(t, s.Has(`c`))
}

func TestSeparate(t *testing.T) {
	s := NewSetInits([]int{1, 2, 3, 8})
	s2 := NewSetInits([]int{2, 3})
	s3 := NewSetInits([]int{4, 2})

	s4 := s.Copy()
	s.Separate(s2)
	s4.Separate(s3)
	wants1 := []int{1, 8}
	wants2 := []int{1, 3, 8}

	ans1 := s.List()
	ans2 := s4.List()
	sort.Ints(ans1)
	sort.Ints(ans2)
	require.Equal(t, wants1, ans1)
	require.Equal(t, wants2, ans2)
}

func TestCopy(t *testing.T) {
	s := NewSetInits([]int{1, 2, 3, 8})

	s4 := s.Copy()
	ans1 := s.List()
	ans2 := s4.List()
	sort.Ints(ans1)
	sort.Ints(ans2)

	require.Equal(t, ans1, ans2)
}

func TestIsEmpty(t *testing.T) {
	s := NewSetInits([]int{1, 2, 3, 8})
	s.Clear()

	require.Equal(t, 0, s.Size())
	require.True(t, s.IsEmpty())
}

func TestIsSubset(t *testing.T) {
	s := NewSetInits([]int{1, 2, 3, 8})
	s2 := NewSetInits([]int{2, 3})

	require.True(t, s.IsSubset(s2))
}

func TestMerge(t *testing.T) {
	s := NewSetInits([]int{1, 2, 3, 8})
	s2 := NewSetInits([]int{2, 3, 7})
	s3 := NewSetInits([]int{1, 2, 3, 7, 8})

	s.Merge(s2)
	ans1 := s.List()
	ans2 := s3.List()
	sort.Ints(ans1)
	sort.Ints(ans2)

	require.Equal(t, ans1, ans2)
}
