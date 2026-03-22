package utils

import (
	"errors"
	"fmt"
)

var (
	//converter errors
	errValueOutOfRangeInt32       = errors.New("compact id value out of range for int32")
	errValueOutOfRangeInt16       = errors.New("compact id value out of range for int16")
	errUnsupportedConvertInt64ToT = errors.New("unsupported type for convertInt64ToT")
	_                             = errors.New("unsupported type for getIntTypeLen")

	// compact_ids related errors
	errCompositeKeyValuesRequired = errors.New("composite key values required")
	_                             = errors.New("unsupported type for compact encoding")
	errOffsetOutOfRange           = errors.New("offset out of range")
	errDstNil                     = errors.New("dst is nil")
	_                             = errors.New("value must be positive")

	// compact_int related errors
	errInvalidCompactInt        = errors.New("invalid compact int value: must be positive")
	errInvalidCompactIntPayload = errors.New("invalid compact int payload: insufficient data for declared length")

	// Centralized errors for compact_int.go
	errBufferInsuficiente       = errors.New("buffer insuficiente")
	errVarintCompactadoInvalido = errors.New("varint compactado inválido: realBytes excede limite do tipo")
)

func errInvalidCompositeKeyEncoding(err error) error {
	return fmt.Errorf("invalid composite key encoding: %w", err)
}

func errInsufficientBuffer(needed int) error {
	return fmt.Errorf("insufficient buffer capacity: need %d free bytes", needed)
}

func errInvalidCompactIntPrefix(prefix byte) error {
	return fmt.Errorf("invalid compact int prefix: %x", prefix)
}

func errCompactVarintTooLarge(needed int) error {
	return fmt.Errorf("compact varint too large: needs %d bytes", needed)
}
