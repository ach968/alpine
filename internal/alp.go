package internal

import (
	"errors"
	"fmt"
	"math"
)

var pow10Table [18]float64

func init() {
	pow10Table[0] = 1
	for i := 1; i <= 17; i++ {
		pow10Table[i] = pow10Table[i-1] * 10
	}
}

func ALPEncode(input []float64, exponent int) ([]int64, int, error) {
	if len(input) == 0 {
		return nil, 0, errors.New("input cannot be empty")
	}

	if exponent < 0 {
		exponent = detectPrecision(input)
	}

	// Add bounds check
	if exponent >= len(pow10Table) {
		return nil, 0, fmt.Errorf("exponent %d exceeds maximum %d", exponent, len(pow10Table)-1)
	}

	multiplier := pow10Table[exponent]
	result := make([]int64, len(input))
	for i, val := range input {
		result[i] = int64(math.Round(val * multiplier))
	}

	return result, exponent, nil
}

func ALPDecode(input []int64, exponent int) []float64 {
	if exponent < 0 || exponent >= len(pow10Table) {
		// Return empty or handle error - for now return empty slice
		return []float64{}
	}
	multiplier := pow10Table[exponent]
	result := make([]float64, len(input))
	for i, val := range input {
		result[i] = float64(val) / multiplier
	}
	return result
}

func detectPrecision(data []float64) int {
	const maxExp = 17
	const maxInt64 = float64(math.MaxInt64)

	for p := 1; p <= maxExp; p++ {
		multiplier := pow10Table[p]
		allMatch := true
		for _, val := range data {
			scaled := val * multiplier
			if scaled > maxInt64 || scaled < -maxInt64 {
				allMatch = false
				break
			}
			rounded := math.Round(scaled)
			// Check if rounding changes the value meaningfully, slight speed optimization
			if math.Abs(scaled-rounded) >= 0.5 {
				allMatch = false
				break
			}
			// Verify round-trip
			restored := rounded / multiplier
			if restored != val {
				allMatch = false
				break
			}
		}
		if allMatch {
			return p
		}
	}
	return 0
}
