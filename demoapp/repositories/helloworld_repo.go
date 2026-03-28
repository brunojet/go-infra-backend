package repositories

import (
	"context"
	"database/sql"

	"github.com/brunojet/go-infra-backend/internal/ports/repositories"
	dbcontracts "github.com/brunojet/go-infra-backend/pkg/database/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"gorm.io/gorm"
)

type HelloWorld struct {
	ID      string `gorm:"primaryKey"`
	Message sql.NullString
}

func (h HelloWorld) TableName() string {
	return "hello_world"
}

func (HelloWorld) WhereOnConflict(tx *gorm.DB) *gorm.DB {
	return tx
}

type HelloWorldRepo struct {
	contracts.Repository[HelloWorld]
}

func NewHelloWorldRepo(db dbcontracts.DatabaseAdapter) contracts.Repository[HelloWorld] {
	return &HelloWorldRepo{
		Repository: repositories.NewGormRepository[HelloWorld](db),
	}
}

func (h *HelloWorldRepo) FindByMessage(ctx context.Context, message string, out *HelloWorld) error {
	return repositories.MapDbError(h.GormDB().WithContext(ctx).Where("message = ?", message).First(out).Error)
}
