// Package alpine provides high-performance compression for sequential numeric data.
// It supports both float64 (with ALP lossless compression) and int64 data,
// with multiple encoding modes optimized for different data patterns.
package alpine

import (
	"fmt"

	"github.com/ach968/alpine/internal"
)

// Mode represents the encoding strategy
type Mode = internal.Mode

const (
	// ModeFloatALP uses ALP + Predictive Delta for float64 data (lossless)
	ModeFloatALP = internal.ModeFloatALP

	// ModeIntSimpleDelta uses simple delta: value[i] - value[i-1]
	// Best for: Monotonically increasing/decreasing integers (timestamps, counters)
	ModeIntSimpleDelta = internal.ModeIntSimpleDelta

	// ModeIntXOR uses XOR difference: value[i] ^ value[i-1]
	// Best for: Database keys, hashes, data with bit patterns
	ModeIntXOR = internal.ModeIntXOR
)

// Options configures the encoding process
type Options struct {
	Mode        Mode // Encoding mode (default: ModeFloatALP for floats)
	RiceParam   int  // Golomb-Rice parameter (0 = auto-detect)
	ALPExponent int  // For ModeFloatALP: precision (-1 = auto-detect, 0 = integers)
}

// Encode compresses float64 data using the specified mode.
// For ModeFloatALP, uses ALP + Predictive Delta encoding (lossless).
func Encode(input []float64, opts Options) ([]byte, error) {
	if len(input) < 2 {
		return nil, fmt.Errorf("input must have at least 2 elements, got %d", len(input))
	}

	if opts.Mode == 0 {
		opts.Mode = ModeFloatALP // Default mode for floats
	}

	if opts.Mode != ModeFloatALP {
		return nil, fmt.Errorf("mode %v not supported for float64, use ModeFloatALP", opts.Mode)
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
		Mode:       opts.Mode,
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

// EncodeInt64 compresses int64 data using the specified mode.
// Supports ModeIntSimpleDelta and ModeIntXOR.
func EncodeInt64(input []int64, opts Options) ([]byte, error) {
	if len(input) < 2 {
		return nil, fmt.Errorf("input must have at least 2 elements, got %d", len(input))
	}

	if opts.Mode == 0 {
		opts.Mode = ModeIntSimpleDelta // Default mode for ints
	}

	if opts.Mode != ModeIntSimpleDelta && opts.Mode != ModeIntXOR {
		return nil, fmt.Errorf("unsupported mode %v for int64", opts.Mode)
	}

	var deltas []int64
	var first int64
	var second int64
	var err error

	// Step 1: Delta encoding based on mode
	switch opts.Mode {
	case ModeIntSimpleDelta:
		deltas, first, err = internal.SimpleDeltaEncode(input)
		if err != nil {
			return nil, fmt.Errorf("simple delta encode: %w", err)
		}
		second = input[1]
	case ModeIntXOR:
		deltas, first, err = internal.XORDeltaEncode(input)
		if err != nil {
			return nil, fmt.Errorf("xor delta encode: %w", err)
		}
		second = input[1]
	}

	// Step 2: Determine Rice parameter
	riceParam := opts.RiceParam
	if riceParam <= 0 {
		riceParam = internal.AutoRiceParam(deltas)
	}

	// Step 3: ZigZag encoding
	zigzagged, err := internal.ZigZagEncode(deltas)
	if err != nil {
		return nil, fmt.Errorf("zigzag encode: %w", err)
	}

	// Step 4: Golomb-Rice encoding
	packed, err := internal.GolombRiceEncode(zigzagged, riceParam)
	if err != nil {
		return nil, fmt.Errorf("golomb-rice encode: %w", err)
	}

	// Step 5: Build header
	header := &internal.Header{
		Mode:       opts.Mode,
		RiceParam:  riceParam,
		ALPExp:     0, // Not used for int modes
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

	if header.Mode != ModeFloatALP {
		return nil, fmt.Errorf("expected ModeFloatALP, got %v (use DecodeInt64 for int modes)", header.Mode)
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

// DecodeInt64 decompresses data produced by EncodeInt64.
func DecodeInt64(encoded []byte) ([]int64, error) {
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

	if header.Mode != ModeIntSimpleDelta && header.Mode != ModeIntXOR {
		return nil, fmt.Errorf("expected int mode, got %v", header.Mode)
	}

	payload := encoded[internal.HeaderSize:]
	deltaCount := header.ValueCount - 1

	// Decode Golomb-Rice
	zigzagged, err := internal.GolombRiceDecode(payload, 0, deltaCount, header.RiceParam)
	if err != nil {
		return nil, fmt.Errorf("golomb-rice decode: %w", err)
	}

	// Decode ZigZag
	deltas, err := internal.ZigZagDecode(zigzagged)
	if err != nil {
		return nil, fmt.Errorf("zigzag decode: %w", err)
	}

	// Decode based on mode
	var result []int64
	switch header.Mode {
	case ModeIntSimpleDelta:
		result, err = internal.SimpleDeltaDecode(deltas, header.First)
	case ModeIntXOR:
		result, err = internal.XORDeltaDecode(deltas, header.First)
	}

	if err != nil {
		return nil, fmt.Errorf("delta decode: %w", err)
	}

	return result, nil
}

// AutoRiceParam calculates the optimal Rice parameter for given deltas.
// This is a convenience function for advanced users who want to pre-calculate.
func AutoRiceParam(deltas []int64) int {
	return internal.AutoRiceParam(deltas)
}
