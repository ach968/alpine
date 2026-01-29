package internal

import "errors"

// XORDeltaEncode computes XOR deltas: deltas[i] = values[i] ^ values[i-1]
// Returns deltas starting from index 1, plus the first value separately.
// Best for: Database keys, hashes, encryption keys, data with bit patterns.
func XORDeltaEncode(input []int64) (deltas []int64, first int64, err error) {
	if len(input) < 1 {
		return nil, 0, errors.New("input must have at least 1 element")
	}

	first = input[0]

	if len(input) == 1 {
		return []int64{}, first, nil
	}

	deltas = make([]int64, len(input)-1)
	for i := 1; i < len(input); i++ {
		deltas[i-1] = input[i] ^ input[i-1]
	}

	return deltas, first, nil
}

// XORDeltaDecode reconstructs original values from XOR deltas.
// values[0] = first
// values[i] = values[i-1] ^ deltas[i-1]
func XORDeltaDecode(deltas []int64, first int64) ([]int64, error) {
	result := make([]int64, len(deltas)+1)
	result[0] = first

	for i := 1; i < len(result); i++ {
		result[i] = result[i-1] ^ deltas[i-1]
	}

	return result, nil
}
