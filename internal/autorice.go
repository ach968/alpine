package internal

import (
	"math/bits"
)

// AutoRiceParam calculates the optimal Rice parameter for given data.
// Uses the median of absolute values as a heuristic, rounded to nearest power of 2.
// This provides good compression without requiring manual tuning.
func AutoRiceParam(deltas []int64) int {
	if len(deltas) == 0 {
		return 4 // Default reasonable value
	}

	// Collect absolute values
	absValues := make([]uint64, len(deltas))
	for i, d := range deltas {
		if d < 0 {
			absValues[i] = uint64(-d)
		} else {
			absValues[i] = uint64(d)
		}
	}

	// Find median
	median := findMedian(absValues)

	// Round to nearest power of 2 for efficiency and clamp to [1, 64]
	riceParam := min(max(roundToPowerOf2(median), 1), 64)

	return riceParam
}

// findMedian finds the median of a slice
func findMedian(values []uint64) uint64 {
	if len(values) == 0 {
		return 0
	}

	sorted := make([]uint64, len(values))
	copy(sorted, values)

	// Insertion sort (efficient for small slices)
	for i := 1; i < len(sorted); i++ {
		key := sorted[i]
		j := i - 1
		for j >= 0 && sorted[j] > key {
			sorted[j+1] = sorted[j]
			j--
		}
		sorted[j+1] = key
	}

	mid := len(sorted) / 2
	if len(sorted)%2 == 0 {
		return (sorted[mid-1] + sorted[mid]) / 2
	}
	return sorted[mid]
}

// roundToPowerOf2 rounds x to the nearest power of 2
func roundToPowerOf2(x uint64) int {
	if x <= 1 {
		return 1
	}

	highestBit := bits.Len64(x) - 1

	lower := uint64(1) << highestBit
	upper := uint64(1) << (highestBit + 1)

	if x-lower < upper-x {
		return int(lower)
	}
	return int(upper)
}
