package repositories

import (
	"context"
	"database/sql"

	"github.com/brunojet/go-infra-backend/demoapp/models"
	internalrepos "github.com/brunojet/go-infra-backend/internal/ports/repositories"
	dbcontracts "github.com/brunojet/go-infra-backend/pkg/database/contracts"
	portsrepos "github.com/brunojet/go-infra-backend/pkg/ports/repositories"
	"github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	whereVersionIDNeq   = models.ColAppVersionID + " <> ?"
	whereVersionAppIDEq = models.ColAppVersionApplicationID + " = ?"
	whereVersionTmIDEq  = models.ColAppVersionTerminalModelConfigurationID + " = ?"
	whereVersionStageEq = models.ColAppVersionStage + " = ?"
	whereVersionDelNil  = models.ColDeletedAt + " IS NULL"
)

type ApplicationVersionRepository interface {
	contracts.Repository[models.ApplicationVersion]
	LoadCurrentStage(ctx context.Context, scope map[string]any) (int16, error)
	FindStageVersionID(ctx context.Context, applicationID int64, configurationTerminalModelID int64, stage int16) (int64, error)
	ArchiveStageDuplicates(ctx context.Context, inOut *models.ApplicationVersion) error
}

type ApplicationVersionRepo struct {
	contracts.Repository[models.ApplicationVersion]
}

func NewApplicationVersionRepo(db dbcontracts.DatabaseAdapter) ApplicationVersionRepository {
	return &ApplicationVersionRepo{Repository: internalrepos.NewGormRepository[models.ApplicationVersion](db)}
}

func (r *ApplicationVersionRepo) LoadCurrentStage(ctx context.Context, scope map[string]any) (int16, error) {
	tx, err := portsrepos.TxFromContext(ctx)
	if err != nil {
		return 0, err
	}
	var stage int16
	result := tx.Model(&models.ApplicationVersion{}).
		Where(scope).
		Limit(1).
		Pluck(models.ColAppVersionStage, &stage)
	if err := result.Error; err != nil {
		return 0, err
	}
	if result.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return stage, nil
}

func (r *ApplicationVersionRepo) FindStageVersionID(ctx context.Context, applicationID int64, configurationTerminalModelID int64, stage int16) (int64, error) {
	switch stage {
	case models.ApplicationStagePending, models.ApplicationStagePilot, models.ApplicationStageProduction, models.ApplicationStageArchived:
	default:
		return 0, gorm.ErrRecordNotFound
	}
	tx, err := portsrepos.TxFromContext(ctx)
	if err != nil {
		return 0, err
	}
	var versionID int64
	result := tx.Model(&models.ApplicationVersion{}).
		Where(models.ColAppVersionApplicationID+" = ?", applicationID).
		Where(models.ColAppVersionTerminalModelConfigurationID+" = ?", configurationTerminalModelID).
		Where(models.ColAppVersionStage+" = ?", stage).
		Where(models.ColDeletedAt+" IS NULL").
		Order(clause.OrderByColumn{Column: clause.Column{Name: models.ColAppVersionID}, Desc: true}).
		Limit(1).
		Pluck(models.ColAppVersionID, &versionID)
	if err := result.Error; err != nil {
		return 0, err
	}
	if result.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return versionID, nil
}

func validateArchiveVersionRequiredFields(inOut *models.ApplicationVersion) error {
	if inOut == nil {
		return ErrInvalidApplicationVersionModel
	}
	if inOut.ApplicationVersionId == 0 {
		return ErrInvalidApplicationVersionID
	}
	if inOut.ApplicationId == 0 {
		return ErrInvalidApplicationID
	}
	if inOut.TerminalModelConfigurationId == 0 {
		return ErrInvalidConfigurationTerminalModelID
	}
	if !inOut.Stage.Valid {
		return ErrInvalidStage
	}
	return nil
}

func (r *ApplicationVersionRepo) ArchiveStageDuplicates(ctx context.Context, inOut *models.ApplicationVersion) error {
	if err := validateArchiveVersionRequiredFields(inOut); err != nil {
		return err
	}

	tx, err := portsrepos.TxFromContext(ctx)
	if err != nil {
		return err
	}

	query := tx.Session(&gorm.Session{SkipHooks: true}).Model(&models.ApplicationVersion{}).
		Where(whereVersionIDNeq, inOut.ApplicationVersionId).
		Where(whereVersionAppIDEq, inOut.ApplicationId).
		Where(whereVersionTmIDEq, inOut.TerminalModelConfigurationId).
		Where(whereVersionStageEq, inOut.Stage.Int16).
		Where(whereVersionDelNil)

	return query.Updates(map[string]any{
		models.ColAppVersionStage: models.ApplicationStageArchived,
		models.ColDeletedAt:       sql.NullTime{Time: tx.NowFunc(), Valid: true},
	}).Error
}
