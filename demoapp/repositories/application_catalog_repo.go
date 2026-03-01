package repositories

import (
	"github.com/brunojet/go-infra-backend/demoapp/models"
	internalrepos "github.com/brunojet/go-infra-backend/internal/ports/repositories"
	dbcontracts "github.com/brunojet/go-infra-backend/pkg/database/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
)

type ApplicationCatalogRepository interface {
	contracts.Repository[models.ApplicationCatalog]
}

type ApplicationCatalogRepo struct {
	contracts.Repository[models.ApplicationCatalog]
	profileRepo ApplicationProfileRepository
}

func NewApplicationCatalogRepo(db dbcontracts.DatabaseAdapter) ApplicationCatalogRepository {
	return &ApplicationCatalogRepo{
		Repository:  internalrepos.NewGormRepository[models.ApplicationCatalog](db),
		profileRepo: NewApplicationProfileRepo(db),
	}
}
