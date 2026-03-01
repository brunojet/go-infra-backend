package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/brunojet/go-infra-backend/demoapp/models"
	internalrepos "github.com/brunojet/go-infra-backend/internal/ports/repositories"
	dbcontracts "github.com/brunojet/go-infra-backend/pkg/database/contracts"
	portsrepos "github.com/brunojet/go-infra-backend/pkg/ports/repositories"
	"github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	whereProfileIDNeq    = models.ColAppProfileID + " <> ?"
	whereProfileIDEq     = models.ColAppProfileID + " = ?"
	whereProfileAppIDEq  = models.ColAppProfileApplicationID + " = ?"
	whereProfileStageEq  = models.ColAppProfileStage + " = ?"
	whereProfileDelIsNil = models.ColAppProfileDeletedAt + " IS NULL"
)

var (
	errArchiveProfileInvalidProfileID = portsrepos.NewBusinessRuleError(errors.New("application_profile_id must be valid"))
	errArchiveProfileInvalidAppID     = portsrepos.NewBusinessRuleError(errors.New("application_id must be valid"))
	errArchiveProfileStageRequired    = portsrepos.NewBusinessRuleError(errors.New("stage is required"))
	errArchiveProfileModelRequired    = portsrepos.NewBusinessRuleError(errors.New("profile must be valid"))
)

type ApplicationProfileRepository interface {
	contracts.Repository[models.ApplicationProfile]
	LoadCurrentStage(ctx context.Context, profileID int64) (int16, error)
	FindCurrentProductionProfileID(ctx context.Context, applicationID int64) (int64, error)
	ArchiveStageDuplicates(ctx context.Context, inOut *models.ApplicationProfile) error
}

type ApplicationProfileRepo struct {
	contracts.Repository[models.ApplicationProfile]
}

func NewApplicationProfileRepo(db dbcontracts.DatabaseAdapter) ApplicationProfileRepository {
	return &ApplicationProfileRepo{Repository: internalrepos.NewGormRepository[models.ApplicationProfile](db)}
}

func (r *ApplicationProfileRepo) LoadCurrentStage(ctx context.Context, profileID int64) (int16, error) {
	tx, err := portsrepos.TxFromContext(ctx)
	if err != nil {
		return 0, err
	}

	var stage int16
	result := tx.Model(&models.ApplicationProfile{}).
		Where(whereProfileIDEq, profileID).
		Limit(1).
		Pluck(models.ColAppProfileStage, &stage)
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
		return errArchiveProfileModelRequired
	}
	if inOut.ApplicationProfileId == 0 {
		return errArchiveProfileInvalidProfileID
	}
	if inOut.ApplicationId == 0 {
		return errArchiveProfileInvalidAppID
	}
	if !inOut.Stage.Valid {
		return errArchiveProfileStageRequired
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
		models.ColAppProfileStage:     models.ApplicationStageArchived,
		models.ColAppProfileDeletedAt: sql.NullTime{Time: tx.NowFunc(), Valid: true},
	}).Error
}
