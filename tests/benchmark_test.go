package alpine_test

import (
	"testing"

	"github.com/ach968/alpine"
)

func generateTimeSeries(n int) []float64 {
	data := make([]float64, n)
	value := 100.0
	for i := 0; i < n; i++ {
		data[i] = value
		value += 0.5 + float64(i%10)*0.1
	}
	return data
}

func BenchmarkEncode_1K(b *testing.B) {
	data := generateTimeSeries(1000)
	opts := alpine.Options{RiceParam: 4, ALPExponent: -1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := alpine.Encode(data, opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecode_1K(b *testing.B) {
	data := generateTimeSeries(1000)
	opts := alpine.Options{RiceParam: 4, ALPExponent: -1}
	encoded, _ := alpine.Encode(data, opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := alpine.Decode(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncode_10K(b *testing.B) {
	data := generateTimeSeries(10000)
	opts := alpine.Options{RiceParam: 4, ALPExponent: -1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := alpine.Encode(data, opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecode_10K(b *testing.B) {
	data := generateTimeSeries(10000)
	opts := alpine.Options{RiceParam: 4, ALPExponent: -1}
	encoded, _ := alpine.Encode(data, opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := alpine.Decode(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}
