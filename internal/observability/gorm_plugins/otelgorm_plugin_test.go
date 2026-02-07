package adapters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOtelGormPlugin_ReturnsPlugin(t *testing.T) {
	a := assert.New(t)
	p, err := NewOtelGormPlugin()
	a.NoError(err)
	a.NotNil(p)
}
