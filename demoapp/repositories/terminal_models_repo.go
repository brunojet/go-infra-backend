package repositories

import (
	"github.com/brunojet/go-infra-backend/demoapp/models"
	internalrepos "github.com/brunojet/go-infra-backend/internal/ports/backend/repositories"
	dbcontracts "github.com/brunojet/go-infra-backend/pkg/infra/database/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/backend/repositories/contracts"
)

type TerminalModelRepo struct {
	contracts.Repository[models.TerminalModel]
}

func NewTerminalModelRepo(db dbcontracts.DatabaseAdapter) contracts.Repository[models.TerminalModel] {
	return &TerminalModelRepo{Repository: internalrepos.NewGormRepository[models.TerminalModel](db)}
}

type TerminalModelConfigurationRepo struct {
	contracts.Repository[models.TerminalModelConfiguration]
}

func NewTerminalModelConfigurationRepo(db dbcontracts.DatabaseAdapter) contracts.Repository[models.TerminalModelConfiguration] {
	return &TerminalModelConfigurationRepo{Repository: internalrepos.NewGormRepository[models.TerminalModelConfiguration](db)}
}
