package utils

// Package-level constants for internal/utils (compact encoding, sizes,
// and small helpers). Keep related magic numbers in one place so callers
// in this package use named symbols instead of literals.
const (
	// compact int special headers
	compactIntHeaderMask    = 0x80 // MSB indica header compactado
	compactIntRealBytesMask = 0x7F // 7 bits menos significativos: número de bytes

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
