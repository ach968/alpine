package alpine

import (
	"testing"
)

func TestEncodeInt64_SimpleDelta(t *testing.T) {
	input := []int64{100, 105, 110, 115, 120}
	opts := Options{
		Mode:      ModeIntSimpleDelta,
		RiceParam: 4,
	}

	encoded, err := EncodeInt64(input, opts)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}

	if len(encoded) == 0 {
		t.Error("expected non-empty encoded data")
	}

	decoded, err := DecodeInt64(encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if len(decoded) != len(input) {
		t.Fatalf("length mismatch: expected %d, got %d", len(input), len(decoded))
	}

	for i := range input {
		if decoded[i] != input[i] {
			t.Errorf("round-trip[%d]: expected %d, got %d", i, input[i], decoded[i])
		}
	}
}

func TestEncodeInt64_XOR(t *testing.T) {
	input := []int64{0x1234, 0x1235, 0x1236, 0x1237}
	opts := Options{
		Mode:      ModeIntXOR,
		RiceParam: 4,
	}

	encoded, err := EncodeInt64(input, opts)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}

	decoded, err := DecodeInt64(encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if len(decoded) != len(input) {
		t.Fatalf("length mismatch: expected %d, got %d", len(input), len(decoded))
	}

	for i := range input {
		if decoded[i] != input[i] {
			t.Errorf("round-trip[%d]: expected %d, got %d", i, input[i], decoded[i])
		}
	}
}

func TestEncode_AutoRiceParam(t *testing.T) {
	input := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	opts := Options{
		Mode:      ModeFloatALP,
		RiceParam: 0, // Auto-detect
	}

	encoded, err := Encode(input, opts)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}

	decoded, err := Decode(encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	for i := range input {
		if decoded[i] != input[i] {
			t.Errorf("round-trip[%d]: expected %f, got %f", i, input[i], decoded[i])
		}
	}
}

func TestEncodeInt64_AutoRiceParam(t *testing.T) {
	input := []int64{1000, 1001, 1002, 1003, 1004}
	opts := Options{
		Mode:      ModeIntSimpleDelta,
		RiceParam: 0, // Auto-detect
	}

	encoded, err := EncodeInt64(input, opts)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}

	decoded, err := DecodeInt64(encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	for i := range input {
		if decoded[i] != input[i] {
			t.Errorf("round-trip[%d]: expected %d, got %d", i, input[i], decoded[i])
		}
	}
}

func TestEncode_WrongMode(t *testing.T) {
	input := []float64{1.0, 2.0, 3.0}

	// Try to use int mode with float data
	opts := Options{Mode: ModeIntSimpleDelta}
	_, err := Encode(input, opts)
	if err == nil {
		t.Error("expected error for wrong mode, got nil")
	}
}

func TestDecode_WrongMode(t *testing.T) {
	input := []int64{1, 2, 3, 4, 5}

	// Encode as int
	opts := Options{Mode: ModeIntSimpleDelta, RiceParam: 4}
	encoded, _ := EncodeInt64(input, opts)

	// Try to decode as float
	_, err := Decode(encoded)
	if err == nil {
		t.Error("expected error for decoding int data with Decode(), got nil")
	}
}

func TestModeConstants(t *testing.T) {
	// Verify mode values are what we expect
	if ModeFloatALP != 0 {
		t.Errorf("ModeFloatALP expected 0, got %d", ModeFloatALP)
	}
	if ModeIntSimpleDelta != 1 {
		t.Errorf("ModeIntSimpleDelta expected 1, got %d", ModeIntSimpleDelta)
	}
	if ModeIntXOR != 2 {
		t.Errorf("ModeIntXOR expected 2, got %d", ModeIntXOR)
	}
}

func TestEncode_InputTooShort(t *testing.T) {
	// Float encode with < 2 elements
	_, err := Encode([]float64{1.0}, Options{Mode: ModeFloatALP})
	if err == nil {
		t.Error("expected error for short float input, got nil")
	}

	_, err = Encode([]float64{}, Options{Mode: ModeFloatALP})
	if err == nil {
		t.Error("expected error for empty float input, got nil")
	}

	// Int encode with < 2 elements
	_, err = EncodeInt64([]int64{1}, Options{Mode: ModeIntSimpleDelta})
	if err == nil {
		t.Error("expected error for short int input, got nil")
	}

	_, err = EncodeInt64([]int64{}, Options{Mode: ModeIntSimpleDelta})
	if err == nil {
		t.Error("expected error for empty int input, got nil")
	}
}

func TestDecode_DataTooShort(t *testing.T) {
	_, err := Decode([]byte{0x01, 0x02})
	if err == nil {
		t.Error("expected error for short data, got nil")
	}

	_, err = DecodeInt64([]byte{0x01, 0x02})
	if err == nil {
		t.Error("expected error for short data, got nil")
	}
}

func TestEncodeInt64_DefaultMode(t *testing.T) {
	// When mode is 0 (not set), should default to ModeIntSimpleDelta
	input := []int64{100, 101, 102, 103}
	opts := Options{
		Mode:      0, // Default
		RiceParam: 4,
	}

	encoded, err := EncodeInt64(input, opts)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}

	decoded, err := DecodeInt64(encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	for i := range input {
		if decoded[i] != input[i] {
			t.Errorf("round-trip[%d]: expected %d, got %d", i, input[i], decoded[i])
		}
	}
}

func TestEncode_DefaultMode(t *testing.T) {
	// When mode is 0 (not set), should default to ModeFloatALP
	input := []float64{1.0, 2.0, 3.0, 4.0}
	opts := Options{
		Mode:      0, // Default
		RiceParam: 4,
	}

	encoded, err := Encode(input, opts)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}

	decoded, err := Decode(encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	for i := range input {
		if decoded[i] != input[i] {
			t.Errorf("round-trip[%d]: expected %f, got %f", i, input[i], decoded[i])
		}
	}
}

func TestAutoRiceParam(t *testing.T) {
	deltas := []int64{1, 2, 3, 4, 5}
	param := AutoRiceParam(deltas)
	if param <= 0 {
		t.Errorf("expected positive rice param, got %d", param)
	}
}
