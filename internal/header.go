package internal

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// Header format:
// Offset  Size  Field
// 0       1B    Mode
// 1       1B    Rice parameter
// 2       1B    ALP exponent (or reserved for int modes)
// 3       1B    Reserved
// 4       8B    First value (int64, big-endian)
// 12      8B    Second value (int64, big-endian)
// 20      4B    Value count (uint32, big-endian)
// 24      ...   Payload

const HeaderSize = 24

type Header struct {
	Mode       Mode
	RiceParam  int
	ALPExp     int
	First      int64
	Second     int64
	ValueCount int
}

func (h *Header) Marshal() []byte {
	buf := make([]byte, HeaderSize)
	buf[0] = h.Mode.Byte()
	buf[1] = byte(h.RiceParam)
	buf[2] = byte(h.ALPExp)
	buf[3] = 0 // Reserved
	binary.BigEndian.PutUint64(buf[4:12], uint64(h.First))
	binary.BigEndian.PutUint64(buf[12:20], uint64(h.Second))
	binary.BigEndian.PutUint32(buf[20:24], uint32(h.ValueCount))
	return buf
}

func Unmarshal(data []byte) (*Header, error) {
	if len(data) < HeaderSize {
		return nil, fmt.Errorf("data too short: need at least %d bytes, got %d", HeaderSize, len(data))
	}

	h := &Header{
		Mode:       ModeFromByte(data[0]),
		RiceParam:  int(data[1]),
		ALPExp:     int(data[2]),
		First:      int64(binary.BigEndian.Uint64(data[4:12])),
		Second:     int64(binary.BigEndian.Uint64(data[12:20])),
		ValueCount: int(binary.BigEndian.Uint32(data[20:24])),
	}

	return h, nil
}

// Validate checks if the header is valid
func (h *Header) Validate() error {
	if h.RiceParam <= 0 {
		return errors.New("rice parameter must be positive")
	}

	if h.ValueCount < 2 {
		return errors.New("value count must be at least 2")
	}

	return nil
}
