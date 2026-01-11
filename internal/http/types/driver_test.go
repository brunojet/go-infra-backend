package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeHTTPDriver(t *testing.T) {
	testCases := []struct {
		input    string
		expected HTTPDriver
	}{
		{"gin", HTTPDriverGin},
		{"chi", HTTPDriverChi},
		{"unknown", HTTPDriver("unknown")},
		{"", HTTPDriverGin},
	}
	for _, tc := range testCases {
		result := NormalizeHTTPDriver(tc.input)
		assert.Equal(t, tc.expected, result, "input: %s", tc.input)
	}
}

func TestHTTPDriver_IsSupported(t *testing.T) {
	testCases := []struct {
		driver    HTTPDriver
		supported bool
	}{
		{HTTPDriverGin, true},
		{HTTPDriverChi, true},
		{HTTPDriver("unknown"), false},
	}
	for _, tc := range testCases {
		assert.Equal(t, tc.supported, tc.driver.IsSupported(), "driver: %v", tc.driver)
	}
}
