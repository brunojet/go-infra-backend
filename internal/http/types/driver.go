package types

import "strings"

// HTTPDriver identifies the HTTP adapter/runtime implementation.
//
// It is intentionally string-based to keep env/config wiring simple.
type HTTPDriver string

const (
	HTTPDriverGin HTTPDriver = "gin"
	HTTPDriverChi HTTPDriver = "chi"
)

// NormalizeHTTPDriver trims/normalizes the input and applies the default.
//
// Empty values default to gin.
func NormalizeHTTPDriver(raw string) HTTPDriver {
	v := HTTPDriver(strings.TrimSpace(strings.ToLower(raw)))
	if v == "" {
		return HTTPDriverGin
	}
	return v
}

func (d HTTPDriver) IsSupported() bool {
	switch d {
	case HTTPDriverGin, HTTPDriverChi:
		return true
	default:
		return false
	}
}
