package internal

import (
	"errors"
	"math"
	"math/bits"
)

// PackedData holds the encoded byte data along with metadata needed for decoding
type PackedData struct {
	Data       []byte
	BitCount   int
	ValueCount int
}

func GolombRiceEncode(input []uint64, m int) (PackedData, error) {
	if len(input) == 0 {
		return PackedData{}, errors.New("input cannot be empty")
	}
	if m <= 0 {
		return PackedData{}, errors.New("m must be positive")
	}

	output := make([]byte, 0)
	var currentByte byte
	var bitPos int = 0
	totalBits := 0

	addBit := func(bit byte) {
		currentByte |= (bit & 1) << (7 - bitPos)
		bitPos++
		totalBits++

		if bitPos == 8 {
			output = append(output, currentByte)
			currentByte = 0
			bitPos = 0
		}
	}

	isPow2 := m > 0 && (m&(m-1)) == 0

	for i := range input {
		var q uint64
		var r uint64

		// m=2^k speedup with bit shifting
		if isPow2 {
			k := bits.TrailingZeros(uint(m))
			q = input[i] >> k
			r = input[i] & uint64(m-1)
		} else {
			q = input[i] / uint64(m)
			r = input[i] % uint64(m)
		}

		// Unary code for quotient. q zeros followed by 1
		for range q {
			addBit(0)
		}
		addBit(1)

		// Append binary representation of remainder
		bitsNeeded := int(math.Ceil(math.Log2(float64(m))))
		for j := bitsNeeded - 1; j >= 0; j-- {
			if (r>>j)&1 == 1 {
				addBit(1)
			} else {
				addBit(0)
			}
		}
	}

	// Flush final byte
	if bitPos > 0 {
		currentByte &= (0xFF << (8 - bitPos))
		output = append(output, currentByte)
	}

	return PackedData{Data: output, BitCount: totalBits, ValueCount: len(input)}, nil
}

func GolombRiceDecode(data []byte, bitCount int, valueCount int, m int) ([]uint64, error) {
	if len(data) == 0 {
		return nil, errors.New("data cannot be empty")
	}
	if m <= 0 {
		return nil, errors.New("m must be positive")
	}
	if valueCount <= 0 {
		return nil, errors.New("valueCount must be positive")
	}

	byteIdx := 0
	bitIdx := 0

	readBit := func() (byte, error) {
		if byteIdx >= len(data) {
			return 0, errors.New("unexpected end of data")
		}
		bit := (data[byteIdx] >> (7 - bitIdx)) & 1
		bitIdx++
		if bitIdx == 8 {
			bitIdx = 0
			byteIdx++
		}
		return bit, nil
	}

	isPow2 := m > 0 && (m&(m-1)) == 0
	bitsNeeded := int(math.Ceil(math.Log2(float64(m))))

	result := make([]uint64, 0, valueCount)

	for range valueCount {
		// Read quotient
		var q uint64
		for {
			bit, err := readBit()
			if err != nil {
				return nil, err
			}
			if bit == 1 {
				break
			}
			q++
		}

		// Read remainder
		var r uint64
		for j := 0; j < bitsNeeded; j++ {
			bit, err := readBit()
			if err != nil {
				return nil, err
			}
			r = (r << 1) | uint64(bit)
		}

		// Reconstruct value
		var value uint64
		if isPow2 {
			k := bits.TrailingZeros(uint(m))
			value = (q << k) | r
		} else {
			value = q*uint64(m) + r
		}

		result = append(result, value)
	}

	return result, nil
}
