package internal

import (
	"testing"
)

func TestHeader_Marshal(t *testing.T) {
	h := &Header{
		Mode:       ModeInt,
		RiceParam:  8,
		ALPExp:     0,
		First:      100,
		Second:     200,
		ValueCount: 1000,
	}

	data := h.Marshal()

	if len(data) != HeaderSize {
		t.Errorf("expected header size %d, got %d", HeaderSize, len(data))
	}

	if data[0] != byte(ModeInt) {
		t.Errorf("expected mode %d, got %d", ModeInt, data[0])
	}

	if data[1] != 8 {
		t.Errorf("expected rice param 8, got %d", data[1])
	}
}

func TestHeader_Unmarshal(t *testing.T) {
	original := &Header{
		Mode:       ModeFloat,
		RiceParam:  4,
		ALPExp:     5,
		First:      42,
		Second:     100,
		ValueCount: 500,
	}

	data := original.Marshal()
	decoded, err := Unmarshal(data)

	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if decoded.Mode != original.Mode {
		t.Errorf("mode: expected %d, got %d", original.Mode, decoded.Mode)
	}

	if decoded.RiceParam != original.RiceParam {
		t.Errorf("rice param: expected %d, got %d", original.RiceParam, decoded.RiceParam)
	}

	if decoded.First != original.First {
		t.Errorf("first: expected %d, got %d", original.First, decoded.First)
	}

	if decoded.Second != original.Second {
		t.Errorf("second: expected %d, got %d", original.Second, decoded.Second)
	}

	if decoded.ValueCount != original.ValueCount {
		t.Errorf("value count: expected %d, got %d", original.ValueCount, decoded.ValueCount)
	}
}

func TestHeader_Validate(t *testing.T) {
	tests := []struct {
		name    string
		header  *Header
		wantErr bool
	}{
		{
			name: "valid header",
			header: &Header{
				RiceParam:  4,
				ValueCount: 10,
			},
			wantErr: false,
		},
		{
			name: "zero rice param",
			header: &Header{
				RiceParam:  0,
				ValueCount: 10,
			},
			wantErr: true,
		},
		{
			name: "single value",
			header: &Header{
				RiceParam:  4,
				ValueCount: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.header.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHeader_Unmarshal_TooShort(t *testing.T) {
	data := make([]byte, HeaderSize-1)
	_, err := Unmarshal(data)
	if err == nil {
		t.Error("expected error for short data")
	}
}

func TestHeader_RoundTrip(t *testing.T) {
	tests := []struct {
		name string
		mode Mode
	}{
		{"Float", ModeFloat},
		{"Int", ModeInt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Header{
				Mode:       tt.mode,
				RiceParam:  8,
				ALPExp:     2,
				First:      12345,
				Second:     67890,
				ValueCount: 100,
			}

			data := h.Marshal()
			decoded, err := Unmarshal(data)
			if err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}

			if decoded.Mode != h.Mode {
				t.Errorf("mode mismatch: expected %d, got %d", h.Mode, decoded.Mode)
			}
		})
	}
}
