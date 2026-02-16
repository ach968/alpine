package internal

import (
	"testing"
)

func TestDelta_RoundTrip(t *testing.T) {
	original := []int64{100, 105, 110, 108, 115}

	deltas, first, second, err := DeltaEncode(original)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}

	decoded, err := DeltaDecode(deltas, first, second)
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

func TestDeltaEncode_TooFewValues(t *testing.T) {
	_, _, _, err := DeltaEncode([]int64{1})
	if err == nil {
		t.Error("expected error for single value, got nil")
	}
}

func TestDeltaEncode_Empty(t *testing.T) {
	_, _, _, err := DeltaEncode([]int64{})
	if err == nil {
		t.Error("expected error for empty input, got nil")
	}
}

func TestDeltaEncode_Basic(t *testing.T) {
	input := []int64{10, 15, 20, 25}

	deltas, first, second, err := DeltaEncode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if first != 10 {
		t.Errorf("first: expected 10, got %d", first)
	}

	if second != 15 {
		t.Errorf("second: expected 15, got %d", second)
	}

	// Predictive delta: value[i] - (value[i-1] + (value[i-1] - value[i-2]))
	// For linear data like this (step of 5), deltas should be zero
	// predicted[2] = 15 + (15-10) = 20, delta = 20-20 = 0
	// predicted[3] = 20 + (20-15) = 25, delta = 25-25 = 0
	expectedDeltas := []int64{0, 0}
	if len(deltas) != len(expectedDeltas) {
		t.Fatalf("deltas length: expected %d, got %d", len(expectedDeltas), len(deltas))
	}

	// Verify actual delta values
	for i, expected := range expectedDeltas {
		if deltas[i] != expected {
			t.Errorf("deltas[%d]: expected %d, got %d", i, expected, deltas[i])
		}
	}
}

func TestDeltaEncode_NonLinear(t *testing.T) {
	// Non-linear data: 10, 20, 35, 50 (increases by 10, then 15, then 15)
	input := []int64{10, 20, 35, 50}

	deltas, first, second, err := DeltaEncode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if first != 10 || second != 20 {
		t.Errorf("expected first=10, second=20, got first=%d, second=%d", first, second)
	}

	// predicted[2] = 20 + (20-10) = 30, delta = 35-30 = 5
	// predicted[3] = 35 + (35-20) = 50, delta = 50-50 = 0
	expectedDeltas := []int64{5, 0}
	if len(deltas) != len(expectedDeltas) {
		t.Fatalf("deltas length: expected %d, got %d", len(expectedDeltas), len(deltas))
	}

	for i, expected := range expectedDeltas {
		if deltas[i] != expected {
			t.Errorf("deltas[%d]: expected %d, got %d", i, expected, deltas[i])
		}
	}
}

func TestDeltaEncode_ExactlyTwoValues(t *testing.T) {
	input := []int64{100, 200}

	deltas, first, second, err := DeltaEncode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if first != 100 || second != 200 {
		t.Errorf("expected first=100, second=200, got first=%d, second=%d", first, second)
	}

	// Should return empty deltas slice (no values to delta encode)
	if len(deltas) != 0 {
		t.Errorf("expected empty deltas for 2 values, got %v", deltas)
	}
}

func TestDeltaEncode_Negative(t *testing.T) {
	input := []int64{-10, -5, 0, 5}

	deltas, first, second, err := DeltaEncode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if first != -10 || second != -5 {
		t.Errorf("unexpected first/second values")
	}

	// Linear progression with step of 5, deltas should be 0
	expectedDeltas := []int64{0, 0}
	for i, expected := range expectedDeltas {
		if deltas[i] != expected {
			t.Errorf("deltas[%d]: expected %d, got %d", i, expected, deltas[i])
		}
	}
}

func TestDeltaDecode(t *testing.T) {
	// Explicit test data instead of relying on TestDeltaEncode_Basic
	deltas := []int64{0, 0} // Deltas from linear progression 10, 15, 20, 25
	first := int64(10)
	second := int64(15)

	result, err := DeltaDecode(deltas, first, second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Decode should reconstruct: 10, 15, 20, 25
	expected := []int64{10, 15, 20, 25}
	if len(result) != len(expected) {
		t.Fatalf("length: expected %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("result[%d]: expected %d, got %d", i, v, result[i])
		}
	}
}

func TestDeltaDecode_EmptyDeltas(t *testing.T) {
	deltas := []int64{}
	first := int64(42)
	second := int64(100)

	result, err := DeltaDecode(deltas, first, second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should return just the first two values
	expected := []int64{42, 100}
	if len(result) != len(expected) {
		t.Fatalf("length: expected %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("result[%d]: expected %d, got %d", i, v, result[i])
		}
	}
}
