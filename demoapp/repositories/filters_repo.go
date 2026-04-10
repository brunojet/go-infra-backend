package repositories

import (
	"github.com/brunojet/go-infra-backend/demoapp/models"
	internalrepos "github.com/brunojet/go-infra-backend/internal/ports/repositories"
	dbcontracts "github.com/brunojet/go-infra-backend/pkg/database/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
)

type FilterTypeRepo struct {
	contracts.Repository[models.FilterType]
}

func NewFilterTypeRepo(db dbcontracts.DatabaseAdapter) contracts.Repository[models.FilterType] {
	return &FilterTypeRepo{Repository: internalrepos.NewGormRepository[models.FilterType](db)}
}

type FilterRepo struct {
	contracts.Repository[models.Filter]
}

func NewFilterRepo(db dbcontracts.DatabaseAdapter) contracts.Repository[models.Filter] {
	return &FilterRepo{Repository: internalrepos.NewGormRepository[models.Filter](db)}
}
