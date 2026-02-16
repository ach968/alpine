package internal

import (
	"testing"
)

func TestAutoRiceParam_SmallValues(t *testing.T) {
	// Small deltas -> small rice param
	deltas := []int64{1, 2, 1, 3, 2, 1, 2}
	param := AutoRiceParam(deltas)

	// Median is 2, rounded to power of 2 = 2
	if param != 2 {
		t.Errorf("expected rice param 2 for small values, got %d", param)
	}
}

func TestAutoRiceParam_LargeValues(t *testing.T) {
	// Large deltas -> larger rice param
	deltas := []int64{100, 200, 150, 180, 120}
	param := AutoRiceParam(deltas)

	// Median is 150, rounded to power of 2 = 128, clamped to 64
	if param != 64 {
		t.Errorf("expected rice param 64 for large values, got %d", param)
	}
}

func TestAutoRiceParam_MixedSigns(t *testing.T) {
	// Mixed positive/negative deltas (absolute values used)
	deltas := []int64{-5, 5, -10, 10, -3, 3}
	param := AutoRiceParam(deltas)

	// Absolute values: 5, 5, 10, 10, 3, 3
	// Sorted: 3, 3, 5, 5, 10, 10
	// Median = (5+5)/2 = 5, rounded to power of 2 = 4
	if param != 4 {
		t.Errorf("expected rice param 4 for mixed signs, got %d", param)
	}
}

func TestAutoRiceParam_Empty(t *testing.T) {
	param := AutoRiceParam([]int64{})
	if param != 4 {
		t.Errorf("expected default rice param 4 for empty input, got %d", param)
	}
}

func TestAutoRiceParam_Uniform(t *testing.T) {
	// All same value
	deltas := []int64{8, 8, 8, 8, 8}
	param := AutoRiceParam(deltas)

	// Median is 8, rounded to power of 2 = 8
	if param != 8 {
		t.Errorf("expected rice param 8 for uniform values, got %d", param)
	}
}

func TestAutoRiceParam_Zeroes(t *testing.T) {
	// All zeros (perfectly predictable data)
	deltas := []int64{0, 0, 0, 0, 0}
	param := AutoRiceParam(deltas)

	// Median is 0, rounded to 1 (minimum)
	if param != 1 {
		t.Errorf("expected rice param 1 for all zeros, got %d", param)
	}
}

func TestRoundToPowerOf2(t *testing.T) {
	tests := []struct {
		input    uint64
		expected int
	}{
		{1, 1},
		{2, 2},
		{3, 4},
		{4, 4},
		{5, 4},
		{6, 8},
		{7, 8},
		{8, 8},
		{15, 16},
		{17, 16},
		{31, 32},
		{33, 32},
		{100, 128},
	}

	for _, tt := range tests {
		result := roundToPowerOf2(tt.input)
		if result != tt.expected {
			t.Errorf("roundToPowerOf2(%d): expected %d, got %d", tt.input, tt.expected, result)
		}
	}
}

func TestAutoRiceParam_SingleElement(t *testing.T) {
	deltas := []int64{42}
	param := AutoRiceParam(deltas)

	// Single value median is that value, rounded to power of 2
	// 42 rounded to nearest power of 2 = 32 or 64
	if param != 32 && param != 64 {
		t.Errorf("expected rice param 32 or 64 for single value, got %d", param)
	}
}

func TestFindMedian(t *testing.T) {
	tests := []struct {
		input    []uint64
		expected uint64
	}{
		{[]uint64{1, 2, 3}, 2},
		{[]uint64{1, 2, 3, 4}, 2},
		{[]uint64{5, 1, 3}, 3},
		{[]uint64{10}, 10},
		{[]uint64{1, 1, 1, 1}, 1},
	}

	for _, tt := range tests {
		result := findMedian(tt.input)
		if result != tt.expected {
			t.Errorf("findMedian(%v): expected %d, got %d", tt.input, tt.expected, result)
		}
	}
}
