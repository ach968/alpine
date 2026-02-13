# alpine

High-performance compression for sequential numeric data in Go. Purpose-built for time-series, database keys, and IoT metrics.

## Features

- **Lossless float compression** using ALP (Adaptive Lossless floating-Point)
- **Multiple integer modes**: SimpleDelta (timestamps), XOR (database keys)
- **Auto-optimization** - automatic rice parameter selection
- **Zero dependencies** - pure Go implementation
- **Simple API** - encode/decode with one function call
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
    
    // Encode with auto-detected rice parameter
    opts := alpine.Options{
        RiceParam: 0,  // 0 = auto-detect
    }
    
    encoded, err := alpine.Encode(data, opts)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Original: %d bytes, Encoded: %d bytes\n", 
        len(data)*8, len(encoded))
    
    // Decode
    decoded, err := alpine.Decode(encoded)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Round-trip successful: %v\n", decoded)
}
```

### Integer Data - Timestamps

```go
timestamps := []int64{1700000000, 1700000060, 1700000120}
encoded, _ := alpine.EncodeInt64(timestamps, alpine.Options{
    Mode: alpine.ModeIntSimpleDelta,
})
decoded, _ := alpine.DecodeInt64(encoded)
```

### Integer Data - Database Keys

```go
keys := []int64{0x1234, 0x1235, 0x1236}
encoded, _ := alpine.EncodeInt64(keys, alpine.Options{
    Mode: alpine.ModeIntXOR,
})
decoded, _ := alpine.DecodeInt64(encoded)
```

## Mode Selection Guide

| Mode | Best For | Algorithm |
|------|----------|-----------|
| ModeFloatALP | Float time-series | ALP + Predictive Delta |
| ModeIntSimpleDelta | Sequential integers | value[i] - value[i-1] |
| ModeIntXOR | Database keys, hashes | value[i] ^ value[i-1] |

**Choosing a mode:**
- **ModeFloatALP**: For any floating-point time-series data (temperatures, prices, metrics)
- **ModeIntSimpleDelta**: For monotonically increasing integers like Unix timestamps, counters
- **ModeIntXOR**: For database keys, hashes, or data with similar bit patterns

## API Reference

### Types

```go
type Options struct {
    Mode        Mode // Encoding mode (default: auto-selected based on data type)
    RiceParam   int  // Golomb-Rice parameter (0 = auto-detect)
    ALPExponent int  // For ModeFloatALP: precision (-1 = auto-detect, 0 = integers)
}

type Mode int

const (
    ModeFloatALP       Mode // Float data with ALP
    ModeIntSimpleDelta Mode // Integer simple delta
    ModeIntXOR         Mode // Integer XOR delta
)
```

### Functions

#### Encode

```go
func Encode(input []float64, opts Options) ([]byte, error)
```

Compresses float64 data using ALP + Delta + Golomb-Rice encoding.

**Parameters:**
- `input`: Slice of float64 values (minimum 2 elements)
- `opts`: Encoding options (Mode must be ModeFloatALP)

**Returns:**
- Encoded byte slice with header
- Error if input is invalid

#### Decode

```go
func Decode(encoded []byte) ([]float64, error)
```

Decompresses data produced by Encode.

#### EncodeInt64

```go
func EncodeInt64(input []int64, opts Options) ([]byte, error)
```

Compresses int64 data using the specified mode (ModeIntSimpleDelta or ModeIntXOR).

**Parameters:**
- `input`: Slice of int64 values (minimum 2 elements)
- `opts`: Encoding options

**Returns:**
- Encoded byte slice with header
- Error if input is invalid

#### DecodeInt64

```go
func DecodeInt64(encoded []byte) ([]int64, error)
```

Decompresses data produced by EncodeInt64.

#### AutoRiceParam

```go
func AutoRiceParam(deltas []int64) int
```

Calculates the optimal Rice parameter for given deltas. Used internally when RiceParam=0.

## How It Works

### Float Encoding Pipeline

```
[]float64 -> ALP Scale (detect precision, multiply by 10^p)
          -> Delta Encode (predictive delta)
          -> ZigZag (signed -> unsigned)
          -> Golomb-Rice Encode
          -> []byte (with header)
```

### Integer Encoding Pipeline

```
[]int64 -> Delta Encode (Simple or XOR)
        -> ZigZag (signed -> unsigned)
        -> Golomb-Rice Encode
        -> []byte (with header)
```

## Auto-Parameter Selection

When `RiceParam: 0` is specified, the library automatically calculates the optimal Rice parameter:

```go
// Automatic parameter selection
opts := alpine.Options{
    RiceParam: 0,  // Library finds optimal value
}
```

The optimal parameter minimizes the encoded size based on the data distribution.

## License

MIT License - see [LICENSE](LICENSE) file
