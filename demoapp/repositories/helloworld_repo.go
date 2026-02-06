package repositories

import (
	"context"
	"database/sql"

	"github.com/brunojet/go-infra-backend/internal/ports/repositories"
	"github.com/brunojet/go-infra-backend/internal/ports/repositories/contracts"
	"gorm.io/gorm"
)

type HelloWorld struct {
	ID      string `gorm:"primaryKey"`
	Message sql.NullString
}

func (h HelloWorld) TableName() string {
	return "hello_world"
}

type HelloWorldRepo struct {
	contracts.Repository[HelloWorld]
}

func NewHelloWorldRepo(db *gorm.DB) contracts.Repository[HelloWorld] {
	return &HelloWorldRepo{
		Repository: repositories.NewGormRepository[HelloWorld](db),
	}
}

func (h *HelloWorldRepo) FindByMessage(ctx context.Context, message string, out *HelloWorld) error {
	return repositories.MapDbError(h.DB().WithContext(ctx).Where("message = ?", message).First(out).Error)
}
