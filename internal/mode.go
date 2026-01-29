package internal

// Mode represents the encoding strategy used
type Mode int

const (
	// ModeFloatALP uses ALP + Predictive Delta for float64 data
	ModeFloatALP Mode = iota

	// ModeIntSimpleDelta uses simple delta: value[i] - value[i-1]
	// Best for: Monotonically increasing/decreasing integers
	ModeIntSimpleDelta

	// ModeIntXOR uses XOR difference: value[i] ^ value[i-1]
	// Best for: Database keys, hashes, data with bit patterns
	ModeIntXOR
)

// ModeFromByte converts a byte to Mode
func ModeFromByte(b byte) Mode {
	return Mode(b)
}

// Byte returns the byte representation of Mode
func (m Mode) Byte() byte {
	return byte(m)
}
