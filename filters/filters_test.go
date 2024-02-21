package filters

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestRange(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6}
	ans := []int{1, 2, 3, 4}
	var b = []int{}
	Range(a, func(i int) bool {
		if i < 5 {
			b = append(b, i)
			return true
		}
		return false
	})
	require.Equal(t, ans, b)
	spew.Dump(b)
}

func TestFilter(t *testing.T) {
	ans := []int{2, 4, 6}
	a := []int{1, 2, 3, 4, 5, 6}
	b := Filter(a, func(i int) bool {
		return i%2 == 0
	})
	require.Equal(t, ans, b)
	spew.Dump(b)
}

func TestFilterNoSpace(t *testing.T) {
	ans1 := []string{"1", "2", "3"}
	a := []string{"", "1", "", "2", "", "3", ""}
	b := Filter(a, NoSpace)
	require.Equal(t, ans1, b)
	spew.Dump(b)
}

func TestFilterMap(t *testing.T) {
	ans := map[string]int{
		"2": 2,
		"4": 4,
		"6": 6,
	}
	a := map[string]int{
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
		"5": 5,
		"6": 6,
	}
	b := FilterMap(a, func(k string, i int) bool {
		return i%2 == 0
	})
	require.Equal(t, ans, b)
	spew.Dump(b)
}

func TestMap(t *testing.T) {
	ans := []string{"2", "4", "6", "end"}
	a := []int{1, 2, 3, 4, 5, 6}
	b := Map(a, func(i int) ([]string, bool) {
		if i == 6 {
			return []string{fmt.Sprintf(`%v`, i), `end`}, true
		}
		if i%2 == 0 {
			return []string{fmt.Sprintf(`%v`, i)}, true
		} else {
			return nil, false
		}
	})
	require.Equal(t, ans, b)
	spew.Dump(b)
}

func TestFirstInt(t *testing.T) {
	ans1, ans2 := 1, 0
	a := []int{1, 2, 3, 4, 5, 6}
	b, ok := First(a)
	require.True(t, ok)
	require.Equal(t, ans1, b)
	spew.Dump(b)
	c := []int{}
	d, ok := First(c)
	require.False(t, ok)
	require.Equal(t, ans2, d)
	spew.Dump(d)
}

func TestFirstString(t *testing.T) {
	ans1, ans2 := "1", ""
	a := []string{"1", "2", "3", "4", "5", "6"}
	b, ok := First(a)
	require.True(t, ok)
	require.Equal(t, ans1, b)
	spew.Dump(b)
	c := []string{}
	d, ok := First(c)
	require.False(t, ok)
	require.Equal(t, ans2, d)
	spew.Dump(d)
}

func TestORInt(t *testing.T) {
	a, b, c := 1, 0, -1
	res := OR(a, b)
	require.Equal(t, res, a)
	spew.Dump(res)
	res = OR(b, c)
	require.Equal(t, res, c)
	spew.Dump(res)
	res = OR(b, b)
	require.Equal(t, res, b)
	spew.Dump(res)
	res = OR(a, b, c)
	require.Equal(t, res, a)
	spew.Dump(res)

}

func TestORString(t *testing.T) {
	a, b, c := "1", "", "-1"
	res := OR(a, b)
	require.Equal(t, res, a)
	spew.Dump(res)
	res = OR(b, c)
	require.Equal(t, res, c)
	spew.Dump(res)
	res = OR(b, b)
	require.Equal(t, res, b)
	spew.Dump(res)
	res = OR(a, b, c)
	require.Equal(t, res, a)
	spew.Dump(res)
}

func TestORBool(t *testing.T) {
	a, b, c := true, false, true
	res := OR(a, b)
	require.Equal(t, res, a)
	spew.Dump(res)
	res = OR(b, c)
	require.Equal(t, res, c)
	spew.Dump(res)
	res = OR(b, b)
	require.Equal(t, res, b)
	spew.Dump(res)
	res = OR(a, b, c)
	require.Equal(t, res, a)
	spew.Dump(res)
}

func TestAddOrUpdateSlice(t *testing.T) {
	v1 := 1
	v2 := 2
	v3 := 3
	vs := []int{v1, v2, v3}
	vs = AddOrUpdateSlice(vs, []int{2, 3}...)
	require.Equal(t, []int{v1, v2, v3}, vs)
	spew.Dump(vs)
	vs = AddOrUpdateSlice(vs, []int{2, 4}...)
	require.Equal(t, []int{v1, v2, v3, 4}, vs)
	spew.Dump(vs)
}

func TestGetDefault(t *testing.T) {
	testCase := []struct {
		name     string
		val      any
		def      any
		expected any
	}{
		{
			name:     "default val string",
			val:      "a",
			def:      "b",
			expected: "a",
		},
		{
			name:     "default null val string",
			val:      "",
			def:      "b",
			expected: "b",
		},
		{
			name:     "default val int",
			val:      1,
			def:      2,
			expected: 1,
		},
		{
			name:     "default null val int",
			val:      0,
			def:      2,
			expected: 2,
		},
		{
			name:     "default val []slice",
			val:      []string{"a", "b"},
			def:      []string{"b"},
			expected: []string{"a", "b"},
		},
		{
			name:     "empty val []sting",
			val:      nil,
			def:      []string{"a"},
			expected: []string{"a"},
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			res := GetDefault(tc.val, tc.def)
			if !reflect.DeepEqual(res, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, res)
			}
		})
	}
}

func TestMax(t *testing.T) {
	testCase := []struct {
		name     string
		val      []int
		expected int
	}{
		{
			name:     "t1",
			val:      []int{1, 2, 3, 4, 5},
			expected: 5,
		},
		{
			name:     "t2",
			val:      []int{},
			expected: 0,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			res := Max(tc.val...)
			if res != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, res)
			}
		})
	}
}

func TestMin(t *testing.T) {
	testCase := []struct {
		name     string
		val      []int
		expected int
	}{
		{
			name:     "t1",
			val:      []int{1, 2, 3, 4, 5},
			expected: 1,
		},
		{
			name:     "t2",
			val:      []int{},
			expected: 0,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			res := Min(tc.val...)
			if res != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, res)
			}
		})
	}
}

func TestSum(t *testing.T) {
	testCase := []struct {
		name     string
		val      []int
		expected int
	}{
		{
			name:     "t1",
			val:      []int{-1, 2, 3, 4, 5},
			expected: 13,
		},
		{
			name:     "t2",
			val:      []int{},
			expected: 0,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			res := Sum(tc.val...)
			if res != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, res)
			}
		})
	}
}
