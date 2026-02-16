package internal

import "errors"

// Returns deltas starting from idx 1. Faster than predictive delta but less effective for linear data.
func SimpleDeltaEncode(input []int64) (deltas []int64, first int64, err error) {
	if len(input) < 1 {
		return nil, 0, errors.New("input must have at least 1 element")
	}

	first = input[0]

	if len(input) == 1 {
		return []int64{}, first, nil
	}

	deltas = make([]int64, len(input)-1)
	for i := 1; i < len(input); i++ {
		deltas[i-1] = input[i] - input[i-1]
	}

	return deltas, first, nil
}

func SimpleDeltaDecode(deltas []int64, first int64) ([]int64, error) {
	result := make([]int64, len(deltas)+1)
	result[0] = first

	for i := 1; i < len(result); i++ {
		result[i] = result[i-1] + deltas[i-1]
	}

	return result, nil
}
