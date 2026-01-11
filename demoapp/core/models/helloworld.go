package models

import (
	"github.com/brunojet/go-infra-backend/demoapp/core/models/ports"
)

// HelloWorld é o modelo GORM para a entidade HelloWorld.
type HelloWorld struct {
	ID uint `gorm:"primaryKey;autoIncrement"`
	ports.Audit
	Message  string `gorm:"not null"`
	Language string `gorm:"not null"`
}
