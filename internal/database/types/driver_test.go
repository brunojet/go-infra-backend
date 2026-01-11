package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeDBDriver(t *testing.T) {
	testCases := []struct {
		input    string
		expected DBDriver
	}{
		{"mysql", DBDriverMySQL},
		{"sqlite", DBDriverSQLite},
		{"unknown", DBDriver("unknown")},
		{"", DBDriverMySQL},
	}
	for _, tc := range testCases {
		result := NormalizeDBDriver(tc.input)
		assert.Equal(t, tc.expected, result, "input: %s", tc.input)
	}
}

func TestDBDriver_IsSupported(t *testing.T) {
	testCases := []struct {
		driver    DBDriver
		supported bool
	}{
		{DBDriverMySQL, true},
		{DBDriverSQLite, true},
		{DBDriver("unknown"), false},
	}
	for _, tc := range testCases {
		assert.Equal(t, tc.supported, tc.driver.IsSupported(), "driver: %v", tc.driver)
	}
}
