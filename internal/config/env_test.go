package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGet_WithEnvAndDefault(t *testing.T) {
	key := "TEST_ENV_KEY"
	t.Setenv(key, "")
	assert.Equal(t, "def", GetEnv(key, "def"))

	t.Setenv(key, "val")
	assert.Equal(t, "val", GetEnv(key, "def"))

	t.Setenv(key, "")
	assert.Equal(t, "def", GetEnv(key, "def"))
}

func TestValidateEndpoint_Valid(t *testing.T) {
	cases := []string{
		":0",
		":80",
		"localhost:5432",
		"my-host:123",
		"my-host.example:65535",
		"127.0.0.1:8080",
		"0.0.0.0:1",
		"[::1]:443",
	}
	for _, c := range cases {
		require.NoError(t, ValidateEndpoint(c), "case=%q", c)
		assert.True(t, IsValidEndpoint(c), "case=%q", c)
	}
}

func TestValidateEndpoint_Invalid(t *testing.T) {
	cases := []string{
		"",
		"localhost",
		"localhost:",
		"localhost:99999",
		":65536",
		"http://localhost:5432",
		"localhost:80/path",
		"[::1]",
		"::1:80", // must be bracketed
		"-bad:80",
		"bad-:80",
	}
	for _, c := range cases {
		require.Error(t, ValidateEndpoint(c), "case=%q", c)
		if c != "" {
			assert.False(t, IsValidEndpoint(c), "case=%q", c)
		}
	}
}

func TestValidateHost_Valid(t *testing.T) {
	cases := []string{
		"localhost",
		"db",
		"db.example.com",
		"127.0.0.1",
		"[::1]",
		"::1",
	}
	for _, c := range cases {
		require.NoError(t, ValidateHost(c), "case=%q", c)
		assert.True(t, IsValidHost(c), "case=%q", c)
	}
}

func TestValidateHost_Invalid(t *testing.T) {
	cases := []string{"", "-nope", "nope-", "a..b", "a b"}
	for _, c := range cases {
		require.Error(t, ValidateHost(c), "case=%q", c)
		assert.False(t, IsValidHost(c), "case=%q", c)
	}
}

func TestValidatePort_Valid(t *testing.T) {
	cases := []string{"0", "1", "80", "443", "65535"}
	for _, c := range cases {
		require.NoError(t, ValidatePort(c), "case=%q", c)
		assert.True(t, IsValidPort(c), "case=%q", c)
	}
}

func TestValidatePort_Invalid(t *testing.T) {
	cases := []string{"", "-1", "65536", "99999", "abc"}
	for _, c := range cases {
		require.Error(t, ValidatePort(c), "case=%q", c)
		assert.False(t, IsValidPort(c), "case=%q", c)
	}
}
