package filters

import "testing"

func TestIterateSets(t *testing.T) {
	testCase := []struct {
		name string
		val  []int
		size int
	}{
		{
			`int slice 2`,
			[]int{1, 2, 3, 4, 5, 6, 7, 8},
			2,
		},
		{
			`int slice 4`,
			[]int{1, 2, 3, 4, 5, 6, 7, 8},
			4,
		},
		{
			`int slice 3`,
			[]int{1, 2, 3, 4, 5, 6, 7, 8},
			3,
		},
		{
			`int slice 6`,
			[]int{1, 2, 3, 4, 5, 6, 7, 8},
			6,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			IterateSets(tc.val, tc.size, func(ints []int) {
				t.Logf(`%v`, ints)
			})
		})
	}

}

func TestIteratePartitions(t *testing.T) {
	testCase := []struct {
		name string
		val  []int
		size int
	}{
		{
			`int slice 2`,
			[]int{1, 2, 3, 4, 5, 6, 7, 8},
			2,
		},
		{
			`int slice 4`,
			[]int{1, 2, 3, 4, 5, 6, 7, 8},
			4,
		},
		{
			`int slice 3`,
			[]int{1, 2, 3, 4, 5, 6, 7, 8},
			3,
		},
		{
			`int slice 6`,
			[]int{1, 2, 3, 4, 5, 6, 7, 8},
			6,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			IteratePartitions(tc.val, tc.size, func(ints [][]int) {
				t.Logf(`%v`, ints)
			})
		})
	}
}
