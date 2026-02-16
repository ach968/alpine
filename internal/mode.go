package internal

// Mode represents the encoding strategy used
type Mode int

const (
	// ModeFloat uses ALP + Predictive Delta for float64 data
	ModeFloat Mode = iota

	// ModeInt uses simple delta: value[i] - value[i-1]
	// Best for: Monotonically increasing/decreasing integers
	ModeInt
)

// ModeFromByte converts a byte to Mode
func ModeFromByte(b byte) Mode {
	return Mode(b)
}

// Byte returns the byte representation of Mode
func (m Mode) Byte() byte {
	return byte(m)
}
