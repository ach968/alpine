package alpine

import (
	"testing"
)

func TestEncode_AutoRiceParam(t *testing.T) {
	input := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	opts := Options{
		Mode:      ModeFloat,
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

func TestEncode_WrongMode(t *testing.T) {
	input := []float64{1.0, 2.0, 3.0}

	// Try to use int mode with float data
	opts := Options{Mode: ModeInt}
	_, err := Encode(input, opts)
	if err == nil {
		t.Error("expected error for wrong mode, got nil")
	}
}

func TestDecode_WrongMode(t *testing.T) {
	// ModeInt is for int64, not float64 - should fail
	input := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18}

	// Try to decode as float - should fail since mode is int
	_, err := Decode(input)
	if err == nil {
		t.Error("expected error for decoding int data with Decode(), got nil")
	}
}

func TestModeConstants(t *testing.T) {
	if ModeFloat != 0 {
		t.Errorf("ModeFloat expected 0, got %d", ModeFloat)
	}
	if ModeInt != 1 {
		t.Errorf("ModeInt expected 1, got %d", ModeInt)
	}
}

func TestEncode_InputTooShort(t *testing.T) {
	// Float encode with < 2 elements
	_, err := Encode([]float64{1.0}, Options{Mode: ModeFloat})
	if err == nil {
		t.Error("expected error for short float input, got nil")
	}

	_, err = Encode([]float64{}, Options{Mode: ModeFloat})
	if err == nil {
		t.Error("expected error for empty float input, got nil")
	}
}

func TestDecode_DataTooShort(t *testing.T) {
	_, err := Decode([]byte{0x01, 0x02})
	if err == nil {
		t.Error("expected error for short data, got nil")
	}
}

func TestEncode_DefaultMode(t *testing.T) {
	// When mode is 0 (not set), should default to ModeFloat
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
