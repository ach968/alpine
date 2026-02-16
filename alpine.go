// Package alpine provides high-performance compression for sequential numeric data.
// It supports both float64 (with ALP lossless compression) and int64 data,
// with multiple encoding modes optimized for different data patterns.
package alpine

import (
	"fmt"

	"github.com/ach968/alpine/internal"
)

// Mode represents the encoding strategy
type Mode int

const (
	// ModeFloat uses ALP + Predictive Delta for float64 data (lossless)
	ModeFloat Mode = 0

	// ModeInt uses simple delta: value[i] - value[i-1]
	// Best for: Monotonically increasing/decreasing integers (timestamps, counters)
	ModeInt Mode = 1
)

// Options configures the encoding process
type Options struct {
	Mode        Mode // Encoding mode (default: ModeFloat for floats)
	RiceParam   int  // Golomb-Rice parameter (0 = auto-detect)
	ALPExponent int  // For ModeFloat: precision (-1 = auto-detect, 0 = integers)
}

// FloatEncoder is a builder for encoding float64 data
type FloatEncoder struct {
	data          []float64
	riceParam     int
	precision     int
	autoRiceParam bool
	autoPrecision bool
}

// NewFloatEncoder creates a new FloatEncoder with the given data
func NewFloatEncoder(data []float64) *FloatEncoder {
	return &FloatEncoder{
		data:      data,
		riceParam: 0,
		precision: 0,
	}
}

// WithRiceParam sets the Golomb-Rice parameter for encoding
func (e *FloatEncoder) WithRiceParam(param int) *FloatEncoder {
	e.riceParam = param
	e.autoRiceParam = false
	return e
}

// WithPrecision sets the ALP precision exponent (0 = auto, -1 = auto-detect)
// Valid values are 0-17
func (e *FloatEncoder) WithPrecision(precision int) *FloatEncoder {
	e.precision = precision
	e.autoPrecision = false
	return e
}

// WithAutoRiceParam enables automatic Rice parameter detection
func (e *FloatEncoder) WithAutoRiceParam() *FloatEncoder {
	e.autoRiceParam = true
	return e
}

// WithAutoPrecision enables automatic precision detection
func (e *FloatEncoder) WithAutoPrecision() *FloatEncoder {
	e.autoPrecision = true
	return e
}

// Encode compresses the float64 data and returns the encoded bytes
func (e *FloatEncoder) Encode() ([]byte, error) {
	if len(e.data) < 2 {
		return nil, fmt.Errorf("input must have at least 2 elements, got %d", len(e.data))
	}

	riceParam := e.riceParam
	exponent := e.precision

	if e.autoPrecision {
		exponent = -1
	} else if exponent == 0 {
		exponent = -1
	}

	scaled, exp, err := internal.ALPEncode(e.data, exponent)
	if err != nil {
		return nil, fmt.Errorf("alp encode: %w", err)
	}

	deltas, first, second, err := internal.DeltaEncode(scaled)
	if err != nil {
		return nil, fmt.Errorf("delta encode: %w", err)
	}

	if riceParam <= 0 || e.autoRiceParam {
		riceParam = internal.AutoRiceParam(deltas)
	}

	var zigzagged []uint64
	if len(deltas) > 0 {
		zigzagged, err = internal.ZigZagEncode(deltas)
		if err != nil {
			return nil, fmt.Errorf("zigzag encode: %w", err)
		}
	}

	var packed internal.PackedData
	if len(zigzagged) > 0 {
		packed, err = internal.GolombRiceEncode(zigzagged, riceParam)
		if err != nil {
			return nil, fmt.Errorf("golomb-rice encode: %w", err)
		}
	}

	header := &internal.Header{
		Mode:       internal.ModeFloat,
		RiceParam:  riceParam,
		ALPExp:     exp,
		First:      first,
		Second:     second,
		ValueCount: len(e.data),
	}

	output := make([]byte, internal.HeaderSize+len(packed.Data))
	copy(output, header.Marshal())
	copy(output[internal.HeaderSize:], packed.Data)

	return output, nil
}

// IntEncoder is a builder for encoding int64 data
type IntEncoder struct {
	data          []int64
	riceParam     int
	autoRiceParam bool
}

// NewIntEncoder creates a new IntEncoder with the given data
func NewIntEncoder(data []int64) *IntEncoder {
	return &IntEncoder{
		data:      data,
		riceParam: 0,
	}
}

// WithRiceParam sets the Golomb-Rice parameter for encoding
func (e *IntEncoder) WithRiceParam(param int) *IntEncoder {
	e.riceParam = param
	e.autoRiceParam = false
	return e
}

// WithAutoRiceParam enables automatic Rice parameter detection
func (e *IntEncoder) WithAutoRiceParam() *IntEncoder {
	e.autoRiceParam = true
	return e
}

// Encode compresses the int64 data using predictive delta encoding and returns the encoded bytes
func (e *IntEncoder) Encode() ([]byte, error) {
	if len(e.data) < 2 {
		return nil, fmt.Errorf("input must have at least 2 elements, got %d", len(e.data))
	}

	riceParam := e.riceParam
	if riceParam <= 0 || e.autoRiceParam {
		riceParam = internal.AutoRiceParam(e.data)
	}

	deltas, first, second, err := internal.DeltaEncode(e.data)
	if err != nil {
		return nil, fmt.Errorf("delta encode: %w", err)
	}

	var zigzagged []uint64
	if len(deltas) > 0 {
		zigzagged, err = internal.ZigZagEncode(deltas)
		if err != nil {
			return nil, fmt.Errorf("zigzag encode: %w", err)
		}
	}

	var packed internal.PackedData
	if len(zigzagged) > 0 {
		packed, err = internal.GolombRiceEncode(zigzagged, riceParam)
		if err != nil {
			return nil, fmt.Errorf("golomb-rice encode: %w", err)
		}
	}

	header := &internal.Header{
		Mode:       internal.ModeInt,
		RiceParam:  riceParam,
		ALPExp:     0,
		First:      first,
		Second:     second,
		ValueCount: len(e.data),
	}

	output := make([]byte, internal.HeaderSize+len(packed.Data))
	copy(output, header.Marshal())
	copy(output[internal.HeaderSize:], packed.Data)

	return output, nil
}

// Decoder is a builder for decoding compressed data
type Decoder struct {
	encoded []byte
}

// NewDecoder creates a new Decoder with the given encoded data
func NewDecoder(encoded []byte) *Decoder {
	return &Decoder{
		encoded: encoded,
	}
}

// DecodeFloat decodes the encoded data as float64 values
func (d *Decoder) DecodeFloat() ([]float64, error) {
	if len(d.encoded) < internal.HeaderSize {
		return nil, fmt.Errorf("data too short: need at least %d bytes, got %d", internal.HeaderSize, len(d.encoded))
	}

	header, err := internal.Unmarshal(d.encoded)
	if err != nil {
		return nil, fmt.Errorf("unmarshal header: %w", err)
	}

	if err := header.Validate(); err != nil {
		return nil, fmt.Errorf("invalid header: %w", err)
	}

	if header.Mode != internal.ModeFloat {
		return nil, fmt.Errorf("expected ModeFloat, got %v", header.Mode)
	}

	payload := d.encoded[internal.HeaderSize:]
	deltaCount := header.ValueCount - 2

	var zigzagged []uint64
	if deltaCount > 0 {
		var err error
		zigzagged, err = internal.GolombRiceDecode(payload, 0, deltaCount, header.RiceParam)
		if err != nil {
			return nil, fmt.Errorf("golomb-rice decode: %w", err)
		}
	}

	var deltas []int64
	if len(zigzagged) > 0 {
		var err error
		deltas, err = internal.ZigZagDecode(zigzagged)
		if err != nil {
			return nil, fmt.Errorf("zigzag decode: %w", err)
		}
	}

	scaled, err := internal.DeltaDecode(deltas, header.First, header.Second)
	if err != nil {
		return nil, fmt.Errorf("delta decode: %w", err)
	}

	result := internal.ALPDecode(scaled, header.ALPExp)

	return result, nil
}

// DecodeInt decodes the encoded data as int64 values
func (d *Decoder) DecodeInt() ([]int64, error) {
	if len(d.encoded) < internal.HeaderSize {
		return nil, fmt.Errorf("data too short: need at least %d bytes, got %d", internal.HeaderSize, len(d.encoded))
	}

	header, err := internal.Unmarshal(d.encoded)
	if err != nil {
		return nil, fmt.Errorf("unmarshal header: %w", err)
	}

	if err := header.Validate(); err != nil {
		return nil, fmt.Errorf("invalid header: %w", err)
	}

	if header.Mode != internal.ModeInt {
		return nil, fmt.Errorf("expected ModeInt, got %v", header.Mode)
	}

	payload := d.encoded[internal.HeaderSize:]
	deltaCount := header.ValueCount - 2

	var zigzagged []uint64
	if deltaCount > 0 {
		var err error
		zigzagged, err = internal.GolombRiceDecode(payload, 0, deltaCount, header.RiceParam)
		if err != nil {
			return nil, fmt.Errorf("golomb-rice decode: %w", err)
		}
	}

	var deltas []int64
	if len(zigzagged) > 0 {
		var err error
		deltas, err = internal.ZigZagDecode(zigzagged)
		if err != nil {
			return nil, fmt.Errorf("zigzag decode: %w", err)
		}
	}

	result, err := internal.DeltaDecode(deltas, header.First, header.Second)
	if err != nil {
		return nil, fmt.Errorf("delta decode: %w", err)
	}

	return result, nil
}

// Encode compresses float64 data using the specified mode.
// For ModeFloat, uses ALP + Predictive Delta encoding (lossless).
func Encode(input []float64, opts Options) ([]byte, error) {
	if len(input) < 2 {
		return nil, fmt.Errorf("input must have at least 2 elements, got %d", len(input))
	}

	if opts.Mode == 0 {
		opts.Mode = ModeFloat // Default mode for floats
	}

	if opts.ALPExponent == 0 {
		opts.ALPExponent = -1 // Default: auto-detect precision
	}

	if opts.Mode != ModeFloat {
		return nil, fmt.Errorf("mode %v not supported for float64, use ModeFloat", opts.Mode)
	}

	// Step 1: ALP encoding
	scaled, exponent, err := internal.ALPEncode(input, opts.ALPExponent)
	if err != nil {
		return nil, fmt.Errorf("alp encode: %w", err)
	}

	// Step 2: Predictive delta encoding
	deltas, first, second, err := internal.DeltaEncode(scaled)
	if err != nil {
		return nil, fmt.Errorf("delta encode: %w", err)
	}

	// Step 3: Determine Rice parameter
	riceParam := opts.RiceParam
	if riceParam <= 0 {
		riceParam = internal.AutoRiceParam(deltas)
	}

	// Step 4: ZigZag encoding (skip if no deltas)
	var zigzagged []uint64
	if len(deltas) > 0 {
		zigzagged, err = internal.ZigZagEncode(deltas)
		if err != nil {
			return nil, fmt.Errorf("zigzag encode: %w", err)
		}
	}

	// Step 5: Golomb-Rice encoding (skip if no zigzagged values)
	var packed internal.PackedData
	if len(zigzagged) > 0 {
		packed, err = internal.GolombRiceEncode(zigzagged, riceParam)
		if err != nil {
			return nil, fmt.Errorf("golomb-rice encode: %w", err)
		}
	}

	// Step 6: Build header
	header := &internal.Header{
		Mode:       internal.Mode(opts.Mode),
		RiceParam:  riceParam,
		ALPExp:     exponent,
		First:      first,
		Second:     second,
		ValueCount: len(input),
	}

	// Combine header and payload
	output := make([]byte, internal.HeaderSize+len(packed.Data))
	copy(output, header.Marshal())
	copy(output[internal.HeaderSize:], packed.Data)

	return output, nil
}

// Decode decompresses data produced by Encode (float64).
func Decode(encoded []byte) ([]float64, error) {
	if len(encoded) < internal.HeaderSize {
		return nil, fmt.Errorf("data too short: need at least %d bytes, got %d", internal.HeaderSize, len(encoded))
	}

	header, err := internal.Unmarshal(encoded)
	if err != nil {
		return nil, fmt.Errorf("unmarshal header: %w", err)
	}

	if err := header.Validate(); err != nil {
		return nil, fmt.Errorf("invalid header: %w", err)
	}

	if header.Mode != internal.ModeFloat {
		return nil, fmt.Errorf("expected ModeFloat, got %v", header.Mode)
	}

	payload := encoded[internal.HeaderSize:]
	deltaCount := header.ValueCount - 2

	// Decode Golomb-Rice (skip if no deltas)
	var zigzagged []uint64
	if deltaCount > 0 {
		var err error
		zigzagged, err = internal.GolombRiceDecode(payload, 0, deltaCount, header.RiceParam)
		if err != nil {
			return nil, fmt.Errorf("golomb-rice decode: %w", err)
		}
	}

	// Decode ZigZag (skip if empty)
	var deltas []int64
	if len(zigzagged) > 0 {
		var err error
		deltas, err = internal.ZigZagDecode(zigzagged)
		if err != nil {
			return nil, fmt.Errorf("zigzag decode: %w", err)
		}
	}

	// Decode predictive delta
	scaled, err := internal.DeltaDecode(deltas, header.First, header.Second)
	if err != nil {
		return nil, fmt.Errorf("delta decode: %w", err)
	}

	// ALP decode
	result := internal.ALPDecode(scaled, header.ALPExp)

	return result, nil
}

// AutoRiceParam calculates the optimal Rice parameter for given deltas.
// This is a convenience function for advanced users who want to pre-calculate.
func AutoRiceParam(deltas []int64) int {
	return internal.AutoRiceParam(deltas)
}
