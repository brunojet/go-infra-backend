package models

import (
	"database/sql"
	"errors"

	porterrors "github.com/brunojet/go-infra-backend/pkg/ports/errors"
	"gorm.io/gorm"
)

const (
	tableApplicationProfileScreenshot = "application_profile_screenshot"
	tableApplicationProfileHistory    = "application_profile_history"

	profileStagePending    int16 = applicationStagePending
	profileStageReviewed   int16 = applicationStageReview
	profileStageProduction int16 = applicationStageProduction
	profileStageArchived   int16 = applicationStageArchived

	ColAppProfileID            = "application_profile_id"
	ColAppProfileApplicationID = "application_id"
	ColAppProfileStage         = "stage"
	ColAppProfileDeletedAt     = "deleted_at"

	ColAppCatalogApplicationID = "application_id"
	ColAppCatalogProfileID     = "application_profile_id"
	ColAppCatalogStage         = "stage"
)

var (
	errProfileStageInvalid             = porterrors.NewBusinessRuleError(errors.New("stage must be valid"))
	errProfileStageTransitionInvalid   = porterrors.NewBusinessRuleError(errors.New("invalid stage transition"))
	errProfileStageBackwardsTransition = porterrors.NewBusinessRuleError(errors.New("stage cannot transition backwards"))
	errProfileReviewAtNotAllowed       = porterrors.NewBusinessRuleError(errors.New("review_at must not be set on create"))
	errProfileProductionAtNotAllowed   = porterrors.NewBusinessRuleError(errors.New("production_at must not be set on create"))

	profileStageAllowedTransitions = map[int16]map[int16]struct{}{
		profileStagePending: {
			profileStageReviewed: {},
			profileStageArchived: {},
		},
		profileStageReviewed: {
			profileStageProduction: {},
			profileStageArchived:   {},
		},
		profileStageProduction: {
			profileStageArchived: {},
		},
		profileStageArchived: {},
	}
)

type ApplicationProfileScreenshot struct {
	ApplicationProfileId int64          `gorm:"column:application_profile_id;primaryKey;priority:1"`
	ApplicationImageId   int64          `gorm:"column:application_image_id;primaryKey;priority:2"`
	Position             int16          `gorm:"column:position;not null"`
	CreatedAt            sql.NullTime   `gorm:"autoCreateTime;index:idx_application_profile_screenshot_del_created,priority:2"`
	UpdatedAt            sql.NullTime   `gorm:"autoUpdateTime;index:idx_application_profile_screenshot_del_updated,priority:2"`
	DeletedAt            gorm.DeletedAt `gorm:"index:idx_application_profile_screenshot_del_created,priority:1;index:idx_application_profile_screenshot_del_updated,priority:1"`

	ApplicationProfile *ApplicationProfile `gorm:"-"`
	ApplicationImage   *ApplicationImage   `gorm:"-"`
}

func (ApplicationProfileScreenshot) TableName() string { return tableApplicationProfileScreenshot }

func (a *ApplicationProfileScreenshot) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ApplicationImage != nil {
		if err := a.ApplicationImage.GetOrCreate(tx); err != nil {
			return err
		}
		a.ApplicationImageId = a.ApplicationImage.ApplicationImageId
	}
	return nil
}

type ApplicationProfile struct {
	ApplicationProfileId int64          `gorm:"primaryKey;autoIncrement"`
	ApplicationId        int64          `gorm:"not null;index:idx_application_profile_stage_app,priority:1"`
	Stage                sql.NullInt16  `gorm:"not null;default:0;index:idx_application_profile_stage_app,priority:2"`
	Name                 sql.NullString `gorm:"size:128;index:ux_application_profile_name_app,priority:1"`
	Description          sql.NullString `gorm:"size:255"`
	ApplicationImageId   int64          `gorm:"not null"`
	ReviewAt             sql.NullTime
	ProductionAt         sql.NullTime
	CreatedAt            sql.NullTime   `gorm:"autoCreateTime;index:idx_application_profile_history_del_created,priority:2"`
	UpdatedAt            sql.NullTime   `gorm:"autoUpdateTime;index:idx_application_profile_history_del_updated,priority:2"`
	DeletedAt            gorm.DeletedAt `gorm:"index:idx_application_profile_history_del_created,priority:1;index:idx_application_profile_history_del_updated,priority:1;index:idx_application_profile_stage_app,priority:3"`
	Filters              []Filter       `gorm:"many2many:application_profile_filters;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`

	Application                   *Application
	ApplicationImage              *ApplicationImage
	ApplicationCatalogs           []ApplicationCatalog           `gorm:"foreignKey:ApplicationProfileId;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	ApplicationProfileScreenshots []ApplicationProfileScreenshot `gorm:"foreignKey:ApplicationProfileId;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}

func (ApplicationProfile) TableName() string {
	return tableApplicationProfileHistory
}

func (a ApplicationProfile) ValidateProfileStageTransition(current int16) error {
	target := a.Stage.Int16

	if target == current {
		return nil
	}

	if allowedTargets, ok := profileStageAllowedTransitions[current]; ok {
		if _, allowed := allowedTargets[target]; allowed {
			return nil
		}
	}

	if target < current {
		return errProfileStageBackwardsTransition
	}

	return errProfileStageTransitionInvalid
}

func (a ApplicationProfile) validateRequiredFieldsCreate() error {
	if a.Stage.Valid && a.Stage.Int16 != profileStagePending {
		return errProfileStageInvalid
	} else if a.ReviewAt.Valid {
		return errProfileReviewAtNotAllowed
	} else if a.ProductionAt.Valid {
		return errProfileProductionAtNotAllowed
	}

	return nil
}

func (a *ApplicationProfile) BeforeCreate(tx *gorm.DB) (err error) {
	if err := a.validateRequiredFieldsCreate(); err != nil {
		return err
	}

	if a.ApplicationImage != nil {
		if err := a.ApplicationImage.GetOrCreate(tx); err != nil {
			return err
		}
		a.ApplicationImageId = a.ApplicationImage.ApplicationImageId
	}

	return nil
}

func (a *ApplicationProfile) BeforeUpdate(tx *gorm.DB) (err error) {
	now := sql.NullTime{Time: tx.NowFunc(), Valid: true}

	switch a.Stage.Int16 {
	case profileStageReviewed:
		a.ReviewAt = now
	case profileStageProduction:
		a.ProductionAt = now
	case profileStageArchived:
		a.DeletedAt = gorm.DeletedAt{Time: now.Time, Valid: true}
	}

	return nil
}
