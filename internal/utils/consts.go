package utils

// Package-level constants for internal/utils (compact encoding, sizes,
// and small helpers). Keep related magic numbers in one place so callers
// in this package use named symbols instead of literals.
const (
	// compact varint encoding constants
	compactVarintPrefixBase = 0x80
	compactVarintPrefixMax  = 0x87
	compactVarintMsb        = 0x80
	compactVarintMaxBytes   = 7

	// byte/shift helpers
	zeroByte            = 0x00
	shiftInt16ToTopByte = 24
	shiftInt64ToTopByte = 56
	zigZagSignShift     = 63
	byteSizeShift       = 8

	// composite id encoding
	compositeMaxPartSize = 8

	// integer sizes
	int16Size = 2
	int32Size = 4
	int64Size = 8

	// JSON helpers
	nullJSONString = "null"
)

var (
	nullJSONBytes = []byte(nullJSONString)
)
