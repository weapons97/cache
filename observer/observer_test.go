package observer

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	obs := New(data)

	if !reflect.DeepEqual(obs.Data(), data) {
		t.Errorf("New() = %v, want %v", obs.Data(), data)
	}
}

func TestFrom(t *testing.T) {
	data := []string{"a", "b", "c"}
	obs := From(data)

	if !reflect.DeepEqual(obs.Data(), data) {
		t.Errorf("From() = %v, want %v", obs.Data(), data)
	}
}

func TestLen(t *testing.T) {
	tests := []struct {
		name     string
		data     []int
		expected int
	}{
		{"normal", []int{1, 2, 3}, 3},
		{"empty", []int{}, 0},
		{"single", []int{1}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obs := New(tt.data)
			if obs.Len() != tt.expected {
				t.Errorf("Len() = %d, want %d", obs.Len(), tt.expected)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		data     []int
		expected bool
	}{
		{"not_empty", []int{1, 2, 3}, false},
		{"empty", []int{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obs := New(tt.data)
			if obs.IsEmpty() != tt.expected {
				t.Errorf("IsEmpty() = %v, want %v", obs.IsEmpty(), tt.expected)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6}
	obs := New(data)

	// 过滤偶数
	filtered := obs.Filter(func(x int) bool {
		return x%2 == 0
	})

	expected := []int{2, 4, 6}
	if !reflect.DeepEqual(filtered.Data(), expected) {
		t.Errorf("Filter() = %v, want %v", filtered.Data(), expected)
	}
}

func TestMap(t *testing.T) {
	data := []int{1, 2, 3, 4}
	obs := New(data)

	// 将每个数字映射为两个相同的数字
	mapped := obs.Map(func(x int) ([]int, bool) {
		return []int{x, x}, true
	})

	expected := []int{1, 1, 2, 2, 3, 3, 4, 4}
	if !reflect.DeepEqual(mapped.Data(), expected) {
		t.Errorf("Map() = %v, want %v", mapped.Data(), expected)
	}
}

func TestFirst(t *testing.T) {
	tests := []struct {
		name     string
		data     []int
		expected int
		found    bool
	}{
		{"normal", []int{1, 2, 3}, 1, true},
		{"empty", []int{}, 0, false},
		{"single", []int{5}, 5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obs := New(tt.data)
			result, found := obs.First()
			if found != tt.found {
				t.Errorf("First() found = %v, want %v", found, tt.found)
			}
			if found && result != tt.expected {
				t.Errorf("First() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestOR(t *testing.T) {
	tests := []struct {
		name     string
		data     []int
		expected int
	}{
		{"first_non_zero", []int{0, 1, 2, 3}, 1},
		{"all_zero", []int{0, 0, 0}, 0},
		{"first_non_zero", []int{1, 0, 2}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obs := New(tt.data)
			result := obs.OR()
			if result != tt.expected {
				t.Errorf("OR() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestAddOrUpdate(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		items    []int
		expected []int
	}{
		{"add_new", []int{1, 2, 3}, []int{4, 5}, []int{1, 2, 3, 4, 5}},
		{"update_existing", []int{1, 2, 3}, []int{2, 4}, []int{1, 2, 3, 4}},
		{"empty_input", []int{}, []int{1, 2, 3}, []int{1, 2, 3}},
		{"empty_items", []int{1, 2, 3}, []int{}, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obs := New(tt.input)
			result := obs.AddOrUpdate(tt.items...)
			if !reflect.DeepEqual(result.Data(), tt.expected) {
				t.Errorf("AddOrUpdate() = %v, want %v", result.Data(), tt.expected)
			}
		})
	}
}

func TestGetDefault(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		default_ int
		expected int
	}{
		{"with_data", []int{1, 2, 3}, 0, 1},
		{"empty_data", []int{}, 42, 42},
		{"zero_data", []int{0, 1, 2}, 42, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obs := New(tt.input)
			result := obs.GetDefault(tt.default_)
			if result != tt.expected {
				t.Errorf("GetDefault(%d) = %d, want %d", tt.default_, result, tt.expected)
			}
		})
	}
}
