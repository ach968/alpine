package internal

import (
	"testing"
)

func TestXORDeltaEncode_Basic(t *testing.T) {
	input := []int64{10, 15, 7, 23}

	deltas, first, err := XORDeltaEncode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if first != 10 {
		t.Errorf("first: expected 10, got %d", first)
	}

	// XOR calculations:
	// 15 ^ 10 = 1010 ^ 1111 = 0101 = 5
	// 7 ^ 15 = 1111 ^ 0111 = 1000 = 8
	// 23 ^ 7 = 0111 ^ 10111 = 10000 = 16
	expected := []int64{5, 8, 16}
	if len(deltas) != len(expected) {
		t.Fatalf("deltas length: expected %d, got %d", len(expected), len(deltas))
	}

	for i, v := range expected {
		if deltas[i] != v {
			t.Errorf("deltas[%d]: expected %d, got %d", i, v, deltas[i])
		}
	}
}

func TestXORDeltaEncode_SequentialIDs(t *testing.T) {
	// Database sequential IDs (consecutive integers)
	input := []int64{1000, 1001, 1002, 1003, 1004}

	deltas, first, err := XORDeltaEncode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if first != 1000 {
		t.Errorf("first: expected 1000, got %d", first)
	}

	// Consecutive integers have small XOR differences
	// These small values compress very well with Golomb-Rice
	for i, d := range deltas {
		if d < 0 {
			t.Errorf("deltas[%d]: XOR should produce non-negative, got %d", i, d)
		}
	}
}

func TestXORDeltaEncode_SingleValue(t *testing.T) {
	input := []int64{0xDEADBEEF}

	deltas, first, err := XORDeltaEncode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if first != 0xDEADBEEF {
		t.Errorf("first: expected %d, got %d", 0xDEADBEEF, first)
	}

	if len(deltas) != 0 {
		t.Errorf("expected empty deltas for single value, got %v", deltas)
	}
}

func TestXORDeltaDecode_Basic(t *testing.T) {
	// 15 ^ 10 = 5
	// 7 ^ 15 = 8
	// 23 ^ 7 = 16
	deltas := []int64{5, 8, 16}
	first := int64(10)

	result, err := XORDeltaDecode(deltas, first)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []int64{10, 15, 7, 23}
	if len(result) != len(expected) {
		t.Fatalf("length: expected %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("result[%d]: expected %d, got %d", i, v, result[i])
		}
	}
}

func TestXORDelta_RoundTrip(t *testing.T) {
	// Test with various patterns
	original := []int64{
		0x123456789ABCDEF0,
		0x123456789ABCDEF1,
		0x123456789ABCDEFF,
		0x123456789ABC0000,
	}

	deltas, first, err := XORDeltaEncode(original)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}

	decoded, err := XORDeltaDecode(deltas, first)
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

func TestXORDeltaEncode_Empty(t *testing.T) {
	_, _, err := XORDeltaEncode([]int64{})
	if err == nil {
		t.Error("expected error for empty input, got nil")
	}
}

func TestXORDelta_Symmetric(t *testing.T) {
	// XOR is symmetric: a ^ b = c, then c ^ b = a and c ^ a = b
	input := []int64{100, 200, 300, 400}

	deltas, first, err := XORDeltaEncode(input)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}

	// Verify symmetry: first ^ deltas[0] = input[1]
	if first^deltas[0] != input[1] {
		t.Errorf("symmetry check failed: %d ^ %d = %d, expected %d",
			first, deltas[0], first^deltas[0], input[1])
	}
}
