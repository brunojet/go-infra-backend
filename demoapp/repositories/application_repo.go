package repositories

import (
	"github.com/brunojet/go-infra-backend/demoapp/models"
	internalrepos "github.com/brunojet/go-infra-backend/internal/ports/repositories"
	dbcontracts "github.com/brunojet/go-infra-backend/pkg/database/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
)

type ApplicationRepo struct {
	contracts.Repository[models.Application]
}

func NewApplicationRepo(db dbcontracts.DatabaseAdapter) contracts.Repository[models.Application] {
	return &ApplicationRepo{Repository: internalrepos.NewGormRepository[models.Application](db)}
}

type ApplicationConfigurationRepo struct {
	contracts.Repository[models.ApplicationConfiguration]
}

func NewApplicationConfigurationRepo(db dbcontracts.DatabaseAdapter) contracts.Repository[models.ApplicationConfiguration] {
	return &ApplicationConfigurationRepo{Repository: internalrepos.NewGormRepository[models.ApplicationConfiguration](db)}
}
