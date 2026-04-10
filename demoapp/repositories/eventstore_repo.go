package repositories

import (
	"github.com/brunojet/go-infra-backend/demoapp/models"
	internalrepos "github.com/brunojet/go-infra-backend/internal/ports/repositories"
	dbcontracts "github.com/brunojet/go-infra-backend/pkg/database/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
)

type EventRepo struct {
	contracts.Repository[models.Event]
}

func NewEventRepo(db dbcontracts.DatabaseAdapter) contracts.Repository[models.Event] {
	return &EventRepo{Repository: internalrepos.NewGormRepository[models.Event](db)}
}

type SnapshotRepo struct {
	contracts.Repository[models.Snapshot]
}

func NewSnapshotRepo(db dbcontracts.DatabaseAdapter) contracts.Repository[models.Snapshot] {
	return &SnapshotRepo{Repository: internalrepos.NewGormRepository[models.Snapshot](db)}
}
