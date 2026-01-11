package env_test

import (
	"os"
	"testing"

	"github.com/brunojet/go-infra-backend/internal/config/adapters/env"
)

func TestEnvSource_Lookup(t *testing.T) {
	os.Setenv("FOO", "bar")
	src := env.New()
	val, ok := src.Lookup("FOO")
	if !ok || val != "bar" {
		t.Errorf("expected env FOO=bar, got %v %v", ok, val)
	}
	_, ok = src.Lookup("NOPE")
	if ok {
		t.Error("expected false for missing key")
	}
}
