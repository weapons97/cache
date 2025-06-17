package filters

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestRange(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6}
	ans := []int{1, 2, 3, 4}
	b := Range(a, func(i int) (int, bool) {
		if i < 5 {
			return i, true
		}
		return 0, false
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

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		item     int
		expected bool
	}{
		{"contains", []int{1, 2, 3, 4, 5}, 3, true},
		{"not_contains", []int{1, 2, 3, 4, 5}, 6, false},
		{"empty_slice", []int{}, 1, false},
		{"single_item", []int{1}, 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Contains(tt.slice, tt.item)
			if result != tt.expected {
				t.Errorf("Contains(%v, %d) = %v, want %v", tt.slice, tt.item, result, tt.expected)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"no_duplicates", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"with_duplicates", []int{1, 2, 2, 3, 3, 4}, []int{1, 2, 3, 4}},
		{"empty_slice", []int{}, []int{}},
		{"single_item", []int{1}, []int{1}},
		{"all_duplicates", []int{1, 1, 1, 1}, []int{1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Unique(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Unique(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"normal", []int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
		{"empty", []int{}, []int{}},
		{"single", []int{1}, []int{1}},
		{"two_items", []int{1, 2}, []int{2, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Reverse(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Reverse(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestChunk(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		size     int
		expected [][]int
	}{
		{"normal_chunks", []int{1, 2, 3, 4, 5, 6}, 2, [][]int{{1, 2}, {3, 4}, {5, 6}}},
		{"uneven_chunks", []int{1, 2, 3, 4, 5}, 2, [][]int{{1, 2}, {3, 4}, {5}}},
		{"empty_slice", []int{}, 3, [][]int{}},
		{"size_larger_than_slice", []int{1, 2}, 5, [][]int{{1, 2}}},
		{"size_one", []int{1, 2, 3}, 1, [][]int{{1}, {2}, {3}}},
		{"invalid_size", []int{1, 2, 3}, 0, [][]int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Chunk(tt.input, tt.size)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Chunk(%v, %d) = %v, want %v", tt.input, tt.size, result, tt.expected)
			}
		})
	}
}
