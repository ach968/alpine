package internal

import (
	"testing"
)

func TestGolombRiceEncode_SmallValues(t *testing.T) {
	input := []uint64{1, 2, 3, 4, 5}
	m := 4 // Rice parameter (power of 2 for efficiency)

	packed, err := GolombRiceEncode(input, m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(packed.Data) == 0 {
		t.Error("expected non-empty data")
	}

	if packed.ValueCount != len(input) {
		t.Errorf("ValueCount: expected %d, got %d", len(input), packed.ValueCount)
	}

	if packed.BitCount <= 0 {
		t.Error("expected positive BitCount")
	}
}

func TestGolombRiceDecode(t *testing.T) {
	input := []uint64{1, 2, 3, 4, 5}
	m := 4

	packed, err := GolombRiceEncode(input, m)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}

	decoded, err := GolombRiceDecode(packed.Data, packed.BitCount, packed.ValueCount, m)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if len(decoded) != len(input) {
		t.Fatalf("length: expected %d, got %d", len(input), len(decoded))
	}

	for i, v := range input {
		if decoded[i] != v {
			t.Errorf("decoded[%d]: expected %d, got %d", i, v, decoded[i])
		}
	}
}

func TestGolombRice_VariousParams(t *testing.T) {
	input := []uint64{10, 20, 30, 40, 50}
	params := []int{1, 2, 4, 8, 16}

	for _, m := range params {
		packed, err := GolombRiceEncode(input, m)
		if err != nil {
			t.Errorf("m=%d: encode error: %v", m, err)
			continue
		}

		decoded, err := GolombRiceDecode(packed.Data, packed.BitCount, packed.ValueCount, m)
		if err != nil {
			t.Errorf("m=%d: decode error: %v", m, err)
			continue
		}

		if len(decoded) != len(input) {
			t.Errorf("m=%d: length mismatch", m)
			continue
		}

		for i := range input {
			if decoded[i] != input[i] {
				t.Errorf("m=%d: decoded[%d]: expected %d, got %d", m, i, input[i], decoded[i])
				break
			}
		}
	}
}

func TestGolombRiceEncode_Empty(t *testing.T) {
	_, err := GolombRiceEncode([]uint64{}, 4)
	if err == nil {
		t.Error("expected error for empty input, got nil")
	}
}

func TestGolombRiceEncode_InvalidParam(t *testing.T) {
	input := []uint64{1, 2, 3}

	_, err := GolombRiceEncode(input, 0)
	if err == nil {
		t.Error("expected error for m=0, got nil")
	}

	_, err = GolombRiceEncode(input, -1)
	if err == nil {
		t.Error("expected error for negative m, got nil")
	}
}
