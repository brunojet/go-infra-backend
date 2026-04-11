package gormplugins

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewOtelGormPlugin_NotNil(t *testing.T) {
	p := NewOtelGormPlugin()
	require.NotNil(t, p)
}
