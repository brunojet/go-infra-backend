package utils

import (
	"math"
)

var (
	convertMinInt32 = int64(math.MinInt32)
	convertMaxInt32 = int64(math.MaxInt32)
	convertMinInt16 = int64(math.MinInt16)
	convertMaxInt16 = int64(math.MaxInt16)
)

func convertInt64ToT[T AnyInt](src int64) (T, error) {
	var zero T
	switch any(zero).(type) {
	case int64:
		return any(src).(T), nil
	case int32:
		if src < convertMinInt32 || src > convertMaxInt32 {
			return zero, errValueOutOfRangeInt32
		}
		return any(int32(src)).(T), nil
	case int16:
		if src < convertMinInt16 || src > convertMaxInt16 {
			return zero, errValueOutOfRangeInt16
		}
		return any(int16(src)).(T), nil
	default:
		return zero, errUnsupportedConvertInt64ToT
	}
}

// getIntTypeLen returns the byte length for supported integer types when only
// an empty interface is available at runtime. Moved from compact_ids.go.
func getIntTypeLen[T AnyInt]() int {
	var out T
	switch any(out).(type) {
	case int16:
		return int16Size
	case int32:
		return int32Size
	default:
		return int64Size
	}
}
