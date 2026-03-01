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
)

const (
	whereVersionIDNeq   = models.ColAppVersionID + " <> ?"
	whereVersionAppIDEq = models.ColAppVersionApplicationID + " = ?"
	whereVersionTmIDEq  = models.ColAppVersionTerminalModelConfigurationID + " = ?"
	whereVersionStageEq = models.ColAppVersionStage + " = ?"
	whereVersionDelNil  = models.ColAppVersionDeletedAt + " IS NULL"
)

var (
	errArchiveVersionModelRequired = portsrepos.NewBusinessRuleError(errors.New("version must be valid"))
	errArchiveVersionInvalidID     = portsrepos.NewBusinessRuleError(errors.New("application_version_id must be valid"))
	errArchiveVersionInvalidAppID  = portsrepos.NewBusinessRuleError(errors.New("application_id must be valid"))
	errArchiveVersionInvalidTmID   = portsrepos.NewBusinessRuleError(errors.New("terminal_model_configuration_id must be valid"))
	errArchiveVersionStageRequired = portsrepos.NewBusinessRuleError(errors.New("stage is required"))
)

type ApplicationVersionRepository interface {
	contracts.Repository[models.ApplicationVersion]
	LoadCurrentStage(ctx context.Context, versionID int64) (int16, error)
	ArchiveStageDuplicates(ctx context.Context, inOut *models.ApplicationVersion) error
}

type ApplicationVersionRepo struct {
	contracts.Repository[models.ApplicationVersion]
}

func NewApplicationVersionRepo(db dbcontracts.DatabaseAdapter) ApplicationVersionRepository {
	return &ApplicationVersionRepo{Repository: internalrepos.NewGormRepository[models.ApplicationVersion](db)}
}

func (r *ApplicationVersionRepo) LoadCurrentStage(ctx context.Context, versionID int64) (int16, error) {
	tx, err := portsrepos.TxFromContext(ctx)
	if err != nil {
		return 0, err
	}

	var stage int16
	result := tx.Model(&models.ApplicationVersion{}).
		Where(models.ColAppVersionID+" = ?", versionID).
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

func validateArchiveVersionRequiredFields(inOut *models.ApplicationVersion) error {
	if inOut == nil {
		return errArchiveVersionModelRequired
	}
	if inOut.ApplicationVersionId == 0 {
		return errArchiveVersionInvalidID
	}
	if inOut.ApplicationId == 0 {
		return errArchiveVersionInvalidAppID
	}
	if inOut.TerminalModelConfigurationId == 0 {
		return errArchiveVersionInvalidTmID
	}
	if !inOut.Stage.Valid {
		return errArchiveVersionStageRequired
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
		models.ColAppVersionStage:     models.ApplicationStageArchived,
		models.ColAppVersionDeletedAt: sql.NullTime{Time: tx.NowFunc(), Valid: true},
	}).Error
}
