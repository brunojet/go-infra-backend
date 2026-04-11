package gormplugins

import (
	internalgorm "github.com/brunojet/go-infra-backend/internal/infra/observability/gorm_plugins"
	"gorm.io/gorm"
)

func NewOtelGormPlugin() gorm.Plugin {
	return internalgorm.NewOtelGormPlugin()
}
