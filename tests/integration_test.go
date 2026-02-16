package alpine_test

import (
	"testing"

	"github.com/ach968/alpine"
)

func TestEncode_Basic(t *testing.T) {
	input := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	opts := alpine.Options{
		Mode:        alpine.ModeFloat,
		RiceParam:   4,
		ALPExponent: 0,
	}

	encoded, err := alpine.Encode(input, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(encoded) == 0 {
		t.Error("expected non-empty encoded data")
	}

	// Should have at least header size (24 bytes)
	if len(encoded) < 24 {
		t.Errorf("expected at least %d bytes (header), got %d", 24, len(encoded))
	}
}

func TestDecode_Basic(t *testing.T) {
	original := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	opts := alpine.Options{
		Mode:        alpine.ModeFloat,
		RiceParam:   4,
		ALPExponent: 0,
	}

	encoded, err := alpine.Encode(original, opts)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}

	decoded, err := alpine.Decode(encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if len(decoded) != len(original) {
		t.Fatalf("length mismatch: expected %d, got %d", len(original), len(decoded))
	}

	for i := range original {
		if decoded[i] != original[i] {
			t.Errorf("decoded[%d]: expected %f, got %f", i, original[i], decoded[i])
		}
	}
}

func TestRoundTrip_Decimals(t *testing.T) {
	original := []float64{3.14159, 2.71828, 1.41421, 1.73205}
	opts := alpine.Options{
		Mode:        alpine.ModeFloat,
		RiceParam:   8,
		ALPExponent: -1,
	}

	encoded, err := alpine.Encode(original, opts)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}

	decoded, err := alpine.Decode(encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if len(decoded) != len(original) {
		t.Fatalf("length mismatch: expected %d, got %d", len(original), len(decoded))
	}

	for i := range original {
		if decoded[i] != original[i] {
			t.Errorf("round-trip[%d]: expected %f, got %f", i, original[i], decoded[i])
		}
	}
}

func TestEncode_Empty(t *testing.T) {
	_, err := alpine.Encode([]float64{}, alpine.Options{
		Mode:      alpine.ModeFloat,
		RiceParam: 4,
	})
	if err == nil {
		t.Error("expected error for empty input, got nil")
	}
}

func TestEncode_SingleValue(t *testing.T) {
	_, err := alpine.Encode([]float64{42.0}, alpine.Options{
		Mode:      alpine.ModeFloat,
		RiceParam: 4,
	})
	if err == nil {
		t.Error("expected error for single value, got nil")
	}
}

func TestEncode_AutoRiceParam(t *testing.T) {
	// RiceParam=0 means auto-detect
	input := []float64{1.0, 2.0, 3.0}
	opts := alpine.Options{
		Mode:      alpine.ModeFloat,
		RiceParam: 0, // Auto-detect
	}

	encoded, err := alpine.Encode(input, opts)
	if err != nil {
		t.Fatalf("unexpected error with auto rice param: %v", err)
	}

	if len(encoded) == 0 {
		t.Error("expected non-empty encoded data")
	}
}

func TestDecode_DataTooShort(t *testing.T) {
	_, err := alpine.Decode([]byte{1, 2, 3})
	if err == nil {
		t.Error("expected error for short data, got nil")
	}
}

func TestDecode_InvalidValueCount(t *testing.T) {
	// Create data with invalid value count (< 2)
	// Header format: [1B: Mode] [1B: Rice] [1B: ALP] [1B: Reserved] [8B: First] [8B: Second] [4B: Count]
	data := make([]byte, 24)
	data[0] = 0 // ModeFloat
	data[1] = 4 // Rice param
	data[2] = 0 // ALP Exponent
	// Count at bytes 20-23 = 1 (too few)
	data[23] = 1

	_, err := alpine.Decode(data)
	if err == nil {
		t.Error("expected error for invalid value count, got nil")
	}
}

func TestRoundTrip_TwoValues(t *testing.T) {
	// Test with exactly 2 values
	original := []float64{10.5, 20.25}
	opts := alpine.Options{
		Mode:        alpine.ModeFloat,
		RiceParam:   4,
		ALPExponent: -1,
	}

	encoded, err := alpine.Encode(original, opts)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}

	decoded, err := alpine.Decode(encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if len(decoded) != len(original) {
		t.Fatalf("length mismatch: expected %d, got %d", len(original), len(decoded))
	}

	for i := range original {
		if decoded[i] != original[i] {
			t.Errorf("round-trip[%d]: expected %f, got %f", i, original[i], decoded[i])
		}
	}
}
