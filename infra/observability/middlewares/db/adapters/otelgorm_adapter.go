package adapters

import (
	"gorm.io/gorm"

	// Import the tracing plugin from the official gorm opentelemetry module.
	// The module provides tracing, metrics and logging hooks. We choose the
	// tracing plugin here by default; callers can register it on their DB
	// instance with `db.Use()`.
	tracing "gorm.io/plugin/opentelemetry/tracing"
)

// NewOtelGormPlugin returns the GORM OpenTelemetry tracing plugin instance.
// Callers should register it with `db.Use(...)` on their `*gorm.DB`.
func NewOtelGormPlugin() (gorm.Plugin, error) {
	return tracing.NewPlugin(), nil
}
