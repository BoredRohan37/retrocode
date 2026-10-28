// Package rpc implements the retrocode wire framing:
// 4-byte little-endian length prefix + UTF-8 JSON payload.
// Bit 31 of the length word marks a zstd-compressed payload.
package rpc

import (
	"encoding/binary"
	"fmt"
)

const (
	CompressedBit uint32 = 0x80000000
	MaxFrame      uint32 = 16 * 1024 * 1024
)

// Encode prepends the 4-byte header to payload.
func Encode(payload []byte, compressed bool) ([]byte, error) {
	if uint32(len(payload)) > MaxFrame {
		return nil, fmt.Errorf("frame too large: %d bytes", len(payload))
	}
	length := uint32(len(payload))
	if compressed {
		length |= CompressedBit
	}
	out := make([]byte, 4+len(payload))
	binary.LittleEndian.PutUint32(out[:4], length)
	copy(out[4:], payload)
	return out, nil
}

// DecodeHeader returns payload length and compression flag.
func DecodeHeader(h []byte) (uint32, bool, error) {
	if len(h) != 4 {
		return 0, false, fmt.Errorf("header must be 4 bytes, got %d", len(h))
	}
	raw := binary.LittleEndian.Uint32(h)
	length := raw &^ CompressedBit
	if length > MaxFrame {
		return 0, false, fmt.Errorf("frame too large: %d bytes", length)
	}
	return length, raw&CompressedBit != 0, nil
}
