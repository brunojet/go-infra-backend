package utils

// AnyInt is a local constraint for integer types supported by the compact encoder.
type AnyInt interface {
	~int64 | ~int32 | ~int16
}
