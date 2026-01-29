package internal

import "errors"

func DeltaEncode(input []int64) (deltas []int64, first int64, second int64, err error) {
	if len(input) < 2 {
		return nil, 0, 0, errors.New("input must have at least 2 elements")
	}
	first = input[0]
	second = input[1]
	deltas = make([]int64, len(input)-2)
	for i := 2; i < len(input); i++ {
		predicted := input[i-1] + (input[i-1] - input[i-2])
		deltas[i-2] = input[i] - predicted
	}
	return deltas, first, second, nil
}

func DeltaDecode(deltas []int64, first int64, second int64) ([]int64, error) {
	result := make([]int64, len(deltas)+2)
	result[0] = first
	result[1] = second
	for i := 2; i < len(result); i++ {
		predicted := result[i-1] + (result[i-1] - result[i-2])
		result[i] = deltas[i-2] + predicted
	}
	return result, nil
}
