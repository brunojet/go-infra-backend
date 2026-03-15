package handlers

import (
	"github.com/brunojet/go-infra-backend/demoapp/models"
	hndcontracts "github.com/brunojet/go-infra-backend/pkg/ports/handlers"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

type ApplicationProfileHandler interface {
	hndcontracts.NestedGenericHandler[models.ApplicationProfile, models.ApplicationProfile]
}

type applicationProfileHandler struct {
	hndcontracts.NestedGenericHandler[models.ApplicationProfile, models.ApplicationProfile]
}

func NewApplicationProfileHandler(s svcContracts.NestedService[models.ApplicationProfile, models.ApplicationProfile]) ApplicationProfileHandler {
	return &applicationProfileHandler{
		NestedGenericHandler: hndcontracts.NewNestedGenericHandler[models.ApplicationProfile](IDInt64Parameters, IDInt64Parameters, s),
	}
}
