package repositories

import (
	"context"
	"database/sql"

	"github.com/brunojet/go-infra-backend/debugassert"
	"github.com/brunojet/go-infra-backend/demoapp/models"
	dbcts "github.com/brunojet/go-infra-backend/pkg/database/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/repositories"
	"github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"gorm.io/gorm"
)

type ApplicationCatalogRepository interface {
	contracts.Repository[models.ApplicationCatalog]
}

type applicationCatalogRepo struct {
	contracts.Repository[models.ApplicationCatalog]
}

func NewApplicationCatalogRepo(db dbcts.DatabaseAdapter) ApplicationCatalogRepository {
	return &applicationCatalogRepo{
		Repository: repositories.NewGormRepository[models.ApplicationCatalog](db),
	}
}

// hydrateApplicationVersionID fills the ApplicationVersionId field of an ApplicationCatalog
// if it is not already set, by querying the database for a matching record based on
// ApplicationId, TerminalModelConfigurationId, and Stage fields.
// If the field is already set, it does nothing. Returns an error if the query fails.
func (r *applicationCatalogRepo) hydrateApplicationVersionID(db *gorm.DB, a *models.ApplicationCatalog) error {
	debugassert.Assert(a != nil, "application catalog pointer is nil")
	if a.ApplicationVersionId != nil { // already set, no need to hydrate
		return nil
	}
	existing := make([]sql.NullInt64, 0, 1)
	if err := db.Model(&models.ApplicationCatalog{}).
		Where(whereApplicationIDEq, a.ApplicationId).
		Where(whereTerminalModelIDEq, a.TerminalModelConfigurationId).
		Where(whereStageEq, a.Stage).
		Limit(1).
		Pluck(models.ColAppVersionID, &existing).Error; err != nil {
		if a.Stage == models.CatalogStageReview {
			return nil // for review stage, it's acceptable to not have an existing version, so we ignore not found errors
		}
		return err
	}

	if len(existing) > 0 && existing[0].Valid {
		existingVersionID := existing[0].Int64
		a.ApplicationVersionId = &existingVersionID
	}

	return nil
}

func (r *applicationCatalogRepo) Create(ctx context.Context, model *models.ApplicationCatalog) error {
	db := r.Repository.DbFromContext(ctx)
	if err := r.hydrateApplicationVersionID(db, model); err != nil {
		return err
	}

	return r.Repository.Create(ctx, model)
}
