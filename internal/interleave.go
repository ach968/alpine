package internal

import "errors"

func ZigZagEncode(input []int64) ([]uint64, error) {
	if len(input) == 0 {
		return nil, errors.New("input cannot be empty")
	}
	encoded := make([]uint64, len(input))
	for i, n := range input {
		if n >= 0 {
			encoded[i] = uint64(n) * 2
		} else {
			encoded[i] = uint64(-n)*2 - 1
		}
	}
	return encoded, nil
}

func ZigZagDecode(input []uint64) ([]int64, error) {
	if len(input) == 0 {
		return nil, errors.New("input cannot be empty")
	}
	decoded := make([]int64, len(input))
	for i, n := range input {
		decoded[i] = int64(n>>1) ^ -int64(n&1)
	}
	return decoded, nil
}
