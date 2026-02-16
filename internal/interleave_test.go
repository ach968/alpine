package internal

import (
	"math"
	"testing"
)

func TestZigZagEncode_Positive(t *testing.T) {
	input := []int64{0, 1, 2, 3}

	result, err := ZigZagEncode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []uint64{0, 2, 4, 6}
	if len(result) != len(expected) {
		t.Fatalf("length: expected %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("result[%d]: expected %d, got %d", i, v, result[i])
		}
	}
}

func TestZigZagEncode_Negative(t *testing.T) {
	input := []int64{-1, -2, -3}

	result, err := ZigZagEncode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []uint64{1, 3, 5}
	if len(result) != len(expected) {
		t.Fatalf("length: expected %d, got %d", len(expected), len(result))
	}

	// Verify actual values
	for i, v := range expected {
		if result[i] != v {
			t.Errorf("result[%d]: expected %d, got %d", i, v, result[i])
		}
	}
}

func TestZigZagDecode(t *testing.T) {
	input := []uint64{0, 1, 2, 3, 4, 5}

	result, err := ZigZagDecode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []int64{0, -1, 1, -2, 2, -3}
	if len(result) != len(expected) {
		t.Fatalf("length: expected %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("result[%d]: expected %d, got %d", i, v, result[i])
		}
	}
}

func TestZigZag_RoundTrip(t *testing.T) {
	original := []int64{-1000, -100, -10, 0, 10, 100, 1000}

	encoded, err := ZigZagEncode(original)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}

	decoded, err := ZigZagDecode(encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if len(decoded) != len(original) {
		t.Fatalf("length mismatch: expected %d, got %d", len(original), len(decoded))
	}

	for i := range original {
		if decoded[i] != original[i] {
			t.Errorf("round-trip[%d]: expected %d, got %d", i, original[i], decoded[i])
		}
	}
}

func TestZigZagEncode_Empty(t *testing.T) {
	_, err := ZigZagEncode([]int64{})
	if err == nil {
		t.Error("expected error for empty input, got nil")
	}
}

func TestZigZagDecode_Empty(t *testing.T) {
	_, err := ZigZagDecode([]uint64{})
	if err == nil {
		t.Error("expected error for empty input, got nil")
	}
}

func TestZigZag_BoundaryValues(t *testing.T) {
	// Test max int64 value
	input := []int64{math.MaxInt64}
	encoded, err := ZigZagEncode(input)
	if err != nil {
		t.Fatalf("failed to encode MaxInt64: %v", err)
	}

	decoded, err := ZigZagDecode(encoded)
	if err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if decoded[0] != math.MaxInt64 {
		t.Errorf("MaxInt64 round-trip failed: expected %d, got %d", math.MaxInt64, decoded[0])
	}

	// Test min int64 value (this might expose overflow bug)
	input = []int64{math.MinInt64}
	encoded, err = ZigZagEncode(input)
	if err != nil {
		t.Fatalf("failed to encode MinInt64: %v", err)
	}

	decoded, err = ZigZagDecode(encoded)
	if err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if decoded[0] != math.MinInt64 {
		t.Errorf("MinInt64 round-trip failed: expected %d, got %d", math.MinInt64, decoded[0])
	}
}

func TestZigZagDecode_LargeValues(t *testing.T) {
	// Test decoding larger values
	input := []uint64{1000, 1001, 2000, 2001, 10000}

	result, err := ZigZagDecode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []int64{500, -501, 1000, -1001, 5000}
	if len(result) != len(expected) {
		t.Fatalf("length: expected %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("result[%d]: expected %d, got %d", i, v, result[i])
		}
	}
}
