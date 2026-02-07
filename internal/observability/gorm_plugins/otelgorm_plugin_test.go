package adapters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOtelGormPlugin_ReturnsPlugin(t *testing.T) {
	p := NewOtelGormPlugin()
	assert.NotNil(t, p)
}
