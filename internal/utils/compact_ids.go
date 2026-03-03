package utils

import (
	"encoding/base64"
	"fmt"
)

// EncodeCompactInt64s encodes one or more int64 values into a base64url-safe string.
func EncodeCompactInt64s(values ...int64) string {
	buf := make([]byte, 0, len(values)*10)
	for _, value := range values {
		buf = appendVarint(buf, value)
	}
	return base64.RawURLEncoding.EncodeToString(buf)
}

// DecodeCompactInt64s decodes a base64url-safe compact id and validates the exact arity.
func DecodeCompactInt64s(encoded string, expected int) ([]int64, error) {
	if expected <= 0 {
		return nil, fmt.Errorf("expected must be greater than zero")
	}

	data, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("invalid compact id encoding: %w", err)
	}

	values := make([]int64, 0, expected)
	index := 0
	for index < len(data) {
		value, next, ok := readVarint(data, index)
		if !ok {
			return nil, fmt.Errorf("invalid compact id payload")
		}
		values = append(values, value)
		index = next
	}

	if len(values) != expected {
		return nil, fmt.Errorf("invalid compact id arity: expected %d values, got %d", expected, len(values))
	}

	for _, value := range values {
		if value <= 0 {
			return nil, fmt.Errorf("invalid compact id value")
		}
	}

	return values, nil
}

func appendVarint(dst []byte, value int64) []byte {
	ux := uint64((value << 1) ^ (value >> 63))
	for ux >= 0x80 {
		dst = append(dst, byte(ux)|0x80)
		ux >>= 7
	}
	return append(dst, byte(ux))
}

func readVarint(src []byte, start int) (int64, int, bool) {
	var (
		ux    uint64
		shift uint
	)

	for i := start; i < len(src); i++ {
		b := src[i]
		ux |= uint64(b&0x7F) << shift
		if b < 0x80 {
			value := int64((ux >> 1) ^ uint64((int64(ux&1)<<63)>>63))
			return value, i + 1, true
		}
		shift += 7
		if shift >= 64 {
			return 0, 0, false
		}
	}

	return 0, 0, false
}
