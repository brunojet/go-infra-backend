package repositories

import (
	"context"
	"database/sql"

	"github.com/brunojet/go-infra-backend/demoapp/models"
	dbcontracts "github.com/brunojet/go-infra-backend/pkg/infra/database/contracts"
	portsrepos "github.com/brunojet/go-infra-backend/pkg/ports/backend/repositories"
	"github.com/brunojet/go-infra-backend/pkg/ports/backend/repositories/contracts"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	whereProfileIDNeq    = models.ColAppProfileID + " <> ?"
	whereProfileAppIDEq  = models.ColApplicationID + " = ?"
	whereProfileStageEq  = models.ColApplicationStage + " = ?"
	whereProfileDelIsNil = models.ColDeletedAt + " IS NULL"
)

type ApplicationProfileRepository interface {
	contracts.Repository[models.ApplicationProfile]
	LoadCurrentStage(ctx context.Context, scope map[string]any) (int16, error)
	FindCurrentProductionProfileID(ctx context.Context, applicationID int64) (int64, error)
	AssociateFilters(ctx context.Context, model *models.ApplicationProfile) error
	ArchiveStageDuplicates(ctx context.Context, inOut *models.ApplicationProfile) error
}

type ApplicationProfileRepo struct {
	contracts.Repository[models.ApplicationProfile]
}

func NewApplicationProfileRepo(db dbcontracts.DatabaseAdapter) ApplicationProfileRepository {
	return &ApplicationProfileRepo{Repository: portsrepos.NewGormRepository[models.ApplicationProfile](db)}
}

func (r *ApplicationProfileRepo) LoadCurrentStage(ctx context.Context, scope map[string]any) (int16, error) {
	tx, err := portsrepos.TxFromContext(ctx)
	if err != nil {
		return 0, err
	}
	var stage int16
	result := tx.Model(&models.ApplicationProfile{}).
		Where(scope).
		Limit(1).
		Pluck(models.ColApplicationStage, &stage)
	if err := result.Error; err != nil {
		return 0, err
	}
	if result.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return stage, nil
}

func (r *ApplicationProfileRepo) FindCurrentProductionProfileID(ctx context.Context, applicationID int64) (int64, error) {
	tx, err := portsrepos.TxFromContext(ctx)
	if err != nil {
		return 0, err
	}
	var profileID int64
	result := tx.Model(&models.ApplicationProfile{}).
		Where(whereProfileAppIDEq, applicationID).
		Where(whereProfileStageEq, models.ApplicationStageProduction).
		Order(clause.OrderByColumn{Column: clause.Column{Name: models.ColAppProfileID}, Desc: true}).
		Limit(1).
		Pluck(models.ColAppProfileID, &profileID)
	if err := result.Error; err != nil {
		return 0, err
	}
	if result.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return profileID, nil
}

func validateArchiveProfileRequiredFields(inOut *models.ApplicationProfile) error {
	if inOut == nil {
		return ErrInvalidApplicationProfileModel
	}
	if inOut.ApplicationProfileId == 0 {
		return ErrInvalidApplicationProfileID
	}
	if inOut.ApplicationId == 0 {
		return ErrInvalidApplicationID
	}
	if !inOut.Stage.Valid {
		return ErrInvalidStage
	}
	return nil
}

func (r *ApplicationProfileRepo) ArchiveStageDuplicates(ctx context.Context, inOut *models.ApplicationProfile) error {
	if err := validateArchiveProfileRequiredFields(inOut); err != nil {
		return err
	}
	tx, err := portsrepos.TxFromContext(ctx)
	if err != nil {
		return err
	}
	query := tx.Session(&gorm.Session{SkipHooks: true}).Model(&models.ApplicationProfile{}).
		Where(whereProfileIDNeq, inOut.ApplicationProfileId).
		Where(whereProfileAppIDEq, inOut.ApplicationId).
		Where(whereProfileStageEq, inOut.Stage.Int16).
		Where(whereProfileDelIsNil)
	return query.Updates(map[string]any{
		models.ColApplicationStage: models.ApplicationStageArchived,
		models.ColDeletedAt:        sql.NullTime{Time: tx.NowFunc(), Valid: true},
	}).Error
}

func (r *ApplicationProfileRepo) AssociateFilters(ctx context.Context, model *models.ApplicationProfile) error {
	if len(model.Filters) == 0 {
		return nil
	}
	tx, err := portsrepos.TxFromContext(ctx)
	if err != nil {
		return err
	}
	return tx.Model(model).Association("Filters").Replace(model.Filters)
}
