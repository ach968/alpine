package internal

import (
	"testing"
)

func TestSimpleDeltaEncode_Basic(t *testing.T) {
	input := []int64{10, 15, 20, 25, 30}

	deltas, first, err := SimpleDeltaEncode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if first != 10 {
		t.Errorf("first: expected 10, got %d", first)
	}

	// Deltas should be: 5, 5, 5, 5
	expected := []int64{5, 5, 5, 5}
	if len(deltas) != len(expected) {
		t.Fatalf("deltas length: expected %d, got %d", len(expected), len(deltas))
	}

	for i, v := range expected {
		if deltas[i] != v {
			t.Errorf("deltas[%d]: expected %d, got %d", i, v, deltas[i])
		}
	}
}

func TestSimpleDeltaEncode_SingleValue(t *testing.T) {
	input := []int64{42}

	deltas, first, err := SimpleDeltaEncode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if first != 42 {
		t.Errorf("first: expected 42, got %d", first)
	}

	if len(deltas) != 0 {
		t.Errorf("expected empty deltas for single value, got %v", deltas)
	}
}

func TestSimpleDeltaEncode_Negative(t *testing.T) {
	input := []int64{-10, -5, 0, 5, 10}

	deltas, first, err := SimpleDeltaEncode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if first != -10 {
		t.Errorf("first: expected -10, got %d", first)
	}

	// Deltas should be: 5, 5, 5, 5
	expected := []int64{5, 5, 5, 5}
	for i, v := range expected {
		if deltas[i] != v {
			t.Errorf("deltas[%d]: expected %d, got %d", i, v, deltas[i])
		}
	}
}

func TestSimpleDeltaEncode_Varied(t *testing.T) {
	input := []int64{100, 105, 103, 108, 110}

	deltas, first, err := SimpleDeltaEncode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if first != 100 {
		t.Errorf("first: expected 100, got %d", first)
	}

	// Deltas should be: 5, -2, 5, 2
	expected := []int64{5, -2, 5, 2}
	for i, v := range expected {
		if deltas[i] != v {
			t.Errorf("deltas[%d]: expected %d, got %d", i, v, deltas[i])
		}
	}
}

func TestSimpleDeltaDecode_Basic(t *testing.T) {
	deltas := []int64{5, 5, 5, 5}
	first := int64(10)

	result, err := SimpleDeltaDecode(deltas, first)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []int64{10, 15, 20, 25, 30}
	if len(result) != len(expected) {
		t.Fatalf("length: expected %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("result[%d]: expected %d, got %d", i, v, result[i])
		}
	}
}

func TestSimpleDeltaDecode_Empty(t *testing.T) {
	deltas := []int64{}
	first := int64(42)

	result, err := SimpleDeltaDecode(deltas, first)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 || result[0] != 42 {
		t.Errorf("expected [42], got %v", result)
	}
}

func TestSimpleDelta_RoundTrip(t *testing.T) {
	original := []int64{1000, 1005, 1010, 1008, 1015, 1020}

	deltas, first, err := SimpleDeltaEncode(original)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}

	decoded, err := SimpleDeltaDecode(deltas, first)
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

func TestSimpleDeltaEncode_Empty(t *testing.T) {
	_, _, err := SimpleDeltaEncode([]int64{})
	if err == nil {
		t.Error("expected error for empty input, got nil")
	}
}
