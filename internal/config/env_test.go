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
	a.Equal("def", Get(key, "def"))

	_ = os.Setenv(key, "val")
	a.Equal("val", Get(key, "def"))
	_ = os.Unsetenv(key)
}
