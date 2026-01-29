package alpaca_test

import (
	"testing"

	"github.com/ach968/alpaca"
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
	opts := alpaca.Options{RiceParam: 4, ALPExponent: -1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := alpaca.Encode(data, opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecode_1K(b *testing.B) {
	data := generateTimeSeries(1000)
	opts := alpaca.Options{RiceParam: 4, ALPExponent: -1}
	encoded, _ := alpaca.Encode(data, opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := alpaca.Decode(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncode_10K(b *testing.B) {
	data := generateTimeSeries(10000)
	opts := alpaca.Options{RiceParam: 4, ALPExponent: -1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := alpaca.Encode(data, opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecode_10K(b *testing.B) {
	data := generateTimeSeries(10000)
	opts := alpaca.Options{RiceParam: 4, ALPExponent: -1}
	encoded, _ := alpaca.Encode(data, opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := alpaca.Decode(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}
