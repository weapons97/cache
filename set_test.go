package cache

import (
	"sort"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
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

func TestSetHasAny(t *testing.T) {
	req := []string{`a`, `b`}
	s := NewSetInits(req)

	// 测试单个元素
	require.True(t, s.HasAny(`a`))
	require.True(t, s.HasAny(`b`))
	require.False(t, s.HasAny(`c`))

	// 测试多个元素 - 任何一个存在就返回 true
	require.True(t, s.HasAny(`a`, `b`))
	require.True(t, s.HasAny(`a`, `c`))
	require.True(t, s.HasAny(`c`, `b`))
	require.False(t, s.HasAny(`c`, `d`))

	// 测试空参数
	require.False(t, s.HasAny())

	// 测试整数类型
	intSet := NewSetInits([]int{1, 2, 3, 8})
	require.True(t, intSet.HasAny(1))
	require.True(t, intSet.HasAny(2, 3))
	require.True(t, intSet.HasAny(4, 2))
	require.False(t, intSet.HasAny(4, 5))
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

func TestPowerSet(t *testing.T) {
	// Setup test input
	s := NewSet[int]()
	s.Add(1)
	s.Add(2)
	s.Add(3)

	// Get the power set
	powerSet := s.PowerSet()

	// Verify the number of subsets (should be 2^n = 8 for n=3)
	require.Equal(t, 8, len(powerSet), "Power set of 3 elements should have 8 subsets")

	// Collect all subsets for verification
	var subsets [][]int
	for _, subset := range powerSet {
		items := subset.List()
		sort.Ints(items)
		subsets = append(subsets, items)
	}

	// Expected subsets (including empty set)
	expected := [][]int{
		{},
		{1},
		{2},
		{3},
		{1, 2},
		{1, 3},
		{2, 3},
		{1, 2, 3},
	}

	// Sort both actual and expected for comparison
	sortSubsets := func(subsets [][]int) {
		sort.Slice(subsets, func(i, j int) bool {
			if len(subsets[i]) != len(subsets[j]) {
				return len(subsets[i]) < len(subsets[j])
			}
			for k := 0; k < len(subsets[i]); k++ {
				if subsets[i][k] != subsets[j][k] {
					return subsets[i][k] < subsets[j][k]
				}
			}
			return false
		})
	}

	sortSubsets(subsets)
	sortSubsets(expected)

	// Verify all subsets are present
	require.Equal(t, expected, subsets, "Power set should contain all possible subsets")

	spew.Dump(subsets)
}
