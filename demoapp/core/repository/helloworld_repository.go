package repository

import (
	"github.com/brunojet/go-infra-backend/demoapp/core/models"
	"github.com/brunojet/go-infra-backend/demoapp/core/repository/ports"
	"gorm.io/gorm"
)

type HelloWorldRepository = ports.GormCRUDRepository[models.HelloWorld]

func NewHelloWorldRepository(db *gorm.DB) *HelloWorldRepository {
	return ports.NewGormCRUDRepository[models.HelloWorld](db)
}
