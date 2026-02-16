# alpine

High-performance compression for time-series data in Go. Optimized for float64 and int64 sequences.

## Features

- **Lossless float compression** using [ALP](https://github.com/cwida/ALP) (Adaptive Lossless floating-Point)
- **Integer support** for timestamps, counters, and sequential data
- **Predictive Delta encoding** for optimal time-series compression
- **Auto-optimization** - automatic rice parameter and precision detection
- **Zero dependencies** - pure Go implementation
- **Builder pattern** - fluent API for configuration
- **Fast** - optimized for time-series workloads

## Installation

```bash
go get github.com/ach968/alpine
```

## Quick Start

### Float Data (Time-Series)

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/ach968/alpine"
)

func main() {
    // Sample time-series data
    data := []float64{10.5, 11.2, 12.8, 13.1, 14.5}
    
    // Encode with auto-detected parameters (default behavior)
    encoded := alpine.NewFloatEncoder(data).Encode()
    if encoded == nil {
        log.Fatal("encode failed")
    }
    
    fmt.Printf("Original: %d bytes, Encoded: %d bytes\n", 
        len(data)*8, len(encoded))
    
    // Decode
    decoded := alpine.NewDecoder(encoded).DecodeFloat()
    if decoded == nil {
        log.Fatal("decode failed")
    }
    
    fmt.Printf("Round-trip successful: %v\n", decoded)
}
```

### Integer Data (Timestamps)

```go
// Timestamps
timestamps := []int64{1700000000, 1700000060, 1700000120}
encoded := alpine.NewIntEncoder(timestamps).Encode()
decoded := alpine.NewDecoder(encoded).DecodeInt()
```

### Builder Pattern

```go
// Custom configuration
encoded := alpine.NewFloatEncoder(data).
    WithRiceParam(8).           // Set specific Rice parameter
    WithPrecision(2).            // 2 decimal places
    WithAutoRiceParam().        // Override: auto-detect Rice param
    Encode()

// All auto-detected
encoded := alpine.NewFloatEncoder(data).
    WithAutoPrecision().
    WithAutoRiceParam().
    Encode()
```

## API Reference

### Builder Types

```go
// FloatEncoder - encodes float64 data
alpine.NewFloatEncoder(data []float64) *FloatEncoder

// IntEncoder - encodes int64 data  
alpine.NewIntEncoder(data []int64) *IntEncoder

// Decoder - decodes both float and int data
alpine.NewDecoder(encoded []byte) *Decoder
```

### FloatEncoder Methods

```go
func (e *FloatEncoder) WithRiceParam(param int) *FloatEncoder
func (e *FloatEncoder) WithPrecision(precision int) *FloatEncoder
func (e *FloatEncoder) WithAutoRiceParam() *FloatEncoder
func (e *FloatEncoder) WithAutoPrecision() *FloatEncoder
func (e *FloatEncoder) Encode() ([]byte, error)
```

### IntEncoder Methods

```go
func (e *IntEncoder) WithRiceParam(param int) *IntEncoder
func (e *IntEncoder) WithAutoRiceParam() *IntEncoder
func (e *IntEncoder) Encode() ([]byte, error)
```

### Decoder Methods

```go
func (d *Decoder) DecodeFloat() ([]float64, error)
func (d *Decoder) DecodeInt() ([]int64, error)
```

### Backwards Compatibility

The legacy Options API is still supported:

```go
opts := alpine.Options{
    RiceParam:   0,  // Auto-detect
    ALPExponent: -1, // Auto-detect precision
}
encoded, _ := alpine.Encode(data, opts)
decoded, _ := alpine.Decode(encoded)
```

## How It Works

### Encoding Pipeline

```
[]float64 -> ALP Scale (detect precision, multiply by 10^p)
          -> Predictive Delta Encode
          -> ZigZag (signed -> unsigned)
          -> Golomb-Rice Encode
          -> []byte (with header)

[]int64 -> Predictive Delta Encode
        -> ZigZag (signed -> unsigned)
        -> Golomb-Rice Encode
        -> []byte (with header)
```

### Auto-Parameter Selection

The library automatically detects optimal parameters:

- **Rice parameter**: Calculated from the median of absolute delta values, rounded to nearest power of 2. This heuristic is efficient because Golomb-Rice encoding is optimal when m ≈ median(|deltas|), and powers of 2 enable bit-shift optimization.
- **Precision**: Detected by testing round-trip accuracy (1-17 decimal places)

## References

- **ALP**: Adaptive Lossless Floating-Point Compression - [https://github.com/cwida/ALP](https://github.com/cwida/ALP) (Azim Afroozeh, Leonardo Kuffó, Peter Boncz - ACM SIGMOD 2024)
- **Delta Encoding**: [https://en.wikipedia.org/wiki/Delta_encoding](https://en.wikipedia.org/wiki/Delta_encoding)
- **Golomb-Rice Coding**: [https://en.wikipedia.org/wiki/Golomb_coding](https://en.wikipedia.org/wiki/Golomb_coding)

## License

MIT License - see [LICENSE](LICENSE) file
