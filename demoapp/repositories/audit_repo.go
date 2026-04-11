package repositories

import (
	"github.com/brunojet/go-infra-backend/demoapp/models"
	internalrepos "github.com/brunojet/go-infra-backend/internal/ports/backend/repositories"
	dbcontracts "github.com/brunojet/go-infra-backend/pkg/infra/database/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/backend/repositories/contracts"
)

type AuditEventRepo struct {
	contracts.Repository[models.AuditEvent]
}

func NewAuditEventRepo(db dbcontracts.DatabaseAdapter) contracts.Repository[models.AuditEvent] {
	return &AuditEventRepo{Repository: internalrepos.NewGormRepository[models.AuditEvent](db)}
}

type AuditFieldChangeRepo struct {
	contracts.Repository[models.AuditFieldChange]
}

func NewAuditFieldChangeRepo(db dbcontracts.DatabaseAdapter) contracts.Repository[models.AuditFieldChange] {
	return &AuditFieldChangeRepo{Repository: internalrepos.NewGormRepository[models.AuditFieldChange](db)}
}
