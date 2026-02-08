package gormplugins

import (
	internalgorm "github.com/brunojet/go-infra-backend/internal/observability/gorm_plugins"
	"gorm.io/gorm"
)

// NewOtelGormPlugin returns the GORM OpenTelemetry tracing plugin instance.
// Callers should register it with `db.Use(...)` on their `*gorm.DB`.
func NewOtelGormPlugin() gorm.Plugin {
	return internalgorm.NewOtelGormPlugin()
}
