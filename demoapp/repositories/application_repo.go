package repositories

import (
	"github.com/brunojet/go-infra-backend/demoapp/models"
	dbcontracts "github.com/brunojet/go-infra-backend/pkg/infra/database/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/backend/repositories"
	"github.com/brunojet/go-infra-backend/pkg/ports/backend/repositories/contracts"
	"gorm.io/gorm"
)

const (
	WhereNameEq          = models.ColName + " = ?"
	WhereApplicationIDEq = models.ColApplicationID + " = ?"
	WherePackageNameEq   = models.ColPackageName + " = ?"
)

type ApplicationRepository interface {
	contracts.Repository[models.Application]
	ValidateAppOwnerCreate(tx *gorm.DB, customerId, applicationName string) error
	ValidateAppOwnerUpdate(tx *gorm.DB, customerId string, applicationId int64) error
}

type applicationRepository struct {
	contracts.Repository[models.Application]
}

func NewApplicationRepo(db dbcontracts.DatabaseAdapter) contracts.Repository[models.Application] {
	return &applicationRepository{Repository: repositories.NewGormRepository[models.Application](db)}
}

// validateAppOwnerCreate checks whether there is an existing Application with the
// same name owned by a different customer. candidateCustomer is the customer
// identifier to validate (useful for Create and Update flows). It requires the
// caller to run inside a transaction so SELECT ... FOR UPDATE is effective.
func (r *applicationRepository) ValidateAppOwnerCreate(tx *gorm.DB, customerId, applicationName string) error {
	return r.validateAppOwner(tx, customerId, WhereNameEq, applicationName)
}

// ValidateAppOwnerUpdate validates updates that would change customer ownership or name.
// It uses the Application instance's `CustomerId` as the candidate value and
// also enforces name ownership rule during updates.
func (r *applicationRepository) ValidateAppOwnerUpdate(tx *gorm.DB, customerId string, applicationId int64) error {
	return r.validateAppOwner(tx, customerId, WhereApplicationIDEq, applicationId)
}

func (r *applicationRepository) validateAppOwner(tx *gorm.DB, customerId string, whereSql string, whereArgs ...any) error {
	if err := repositories.ValidateTxWithUpdateLock(tx, repositories.LockValidationSpec[models.Application]{
		SelectColumns: []string{models.ColCustomerID},
		WhereSQL:      whereSql,
		WhereArgs:     whereArgs,
		BlockIfFound: func(existing *models.Application) error {
			if existing.CustomerId.String != customerId {
				return ErrCustomerIdViolation
			}
			return nil
		},
	}); err != nil {
		return err
	}
	return nil
}

type ApplicationConfigurationRepository interface {
	contracts.Repository[models.ApplicationConfiguration]
	ValidateAppConfigurationOwner(tx *gorm.DB, applicationId int64, packageName string) error
}

type applicationConfigurationRepository struct {
	contracts.Repository[models.ApplicationConfiguration]
}

func NewApplicationConfigurationRepo(db dbcontracts.DatabaseAdapter) ApplicationConfigurationRepository {
	return &applicationConfigurationRepository{Repository: repositories.NewGormRepository[models.ApplicationConfiguration](db)}
}

func (a applicationConfigurationRepository) ValidateAppConfigurationOwner(tx *gorm.DB, applicationId int64, packageName string) error {
	return repositories.ValidateTxWithUpdateLock(tx, repositories.LockValidationSpec[models.ApplicationConfiguration]{
		SelectColumns: []string{models.ColApplicationID},
		WhereSQL:      WherePackageNameEq,
		WhereArgs:     []any{packageName},
		BlockIfFound: func(existing *models.ApplicationConfiguration) error {
			if existing.ApplicationId != applicationId {
				return ErrApplicationIdViolation
			}
			return nil
		},
	})
}
