package exporters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConsoleLoggerExporter_ReturnsExporter(t *testing.T) {
	a := assert.New(t)
	e, err := NewConsoleLoggerExporter()
	a.NoError(err)
	a.NotNil(e)
}
