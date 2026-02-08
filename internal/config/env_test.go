package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet_WithEnvAndDefault(t *testing.T) {
	a := assert.New(t)
	key := "TEST_ENV_KEY"
	_ = os.Unsetenv(key)
	a.Equal("def", GetEnv(key, "def"))

	_ = os.Setenv(key, "val")
	a.Equal("val", GetEnv(key, "def"))
	_ = os.Unsetenv(key)
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
		if err := ValidateEndpoint(c); err != nil {
			t.Fatalf("expected valid endpoint %q, got err=%v", c, err)
		}
		if !IsValidEndpoint(c) {
			t.Fatalf("expected IsValidEndpoint(%q)=true", c)
		}
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
		if err := ValidateEndpoint(c); err == nil {
			t.Fatalf("expected invalid endpoint %q", c)
		}
		if c != "" && IsValidEndpoint(c) {
			t.Fatalf("expected IsValidEndpoint(%q)=false", c)
		}
	}
}
