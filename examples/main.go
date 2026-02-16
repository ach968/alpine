package main

import (
	"fmt"
	"log"

	"github.com/ach968/alpine"
)

func builder_pattern() {
	// Sample time-series data
	data := []float64{10.5, 11.2, 12.8, 13.1, 14.5}

	encoded, err := alpine.NewFloatEncoder(data).WithAutoPrecision().WithAutoRiceParam().Encode()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Original: %d bytes, Compressed: %d bytes\n",
		len(data)*8, len(encoded))

	// Decode
	decoded, err := alpine.Decode(encoded)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Round-trip successful: %v\n", decoded)
}

func main() {
	builder_pattern()
}
