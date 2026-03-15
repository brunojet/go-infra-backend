package handlers

import (
	"github.com/brunojet/go-infra-backend/demoapp/models"
	hnd "github.com/brunojet/go-infra-backend/pkg/ports/handlers"
	hndcontracts "github.com/brunojet/go-infra-backend/pkg/ports/handlers/contracts"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

type ApplicationHandler interface {
	hndcontracts.GenericHandler[models.Application, models.Application]
}

type applicationHandler struct {
	hndcontracts.GenericHandler[models.Application, models.Application]
}

func NewapplicationHandler(s svcContracts.Service[models.Application, models.Application]) ApplicationHandler {
	return &applicationHandler{
		GenericHandler: hnd.NewGenericHandler[models.Application](IDInt64Parameters, s),
	}
}

type ApplicationConfigurationHandler interface {
	hndcontracts.NestedGenericHandler[models.ApplicationConfiguration, models.ApplicationConfiguration]
}

type applicationConfigurationHandler struct {
	hndcontracts.NestedGenericHandler[models.ApplicationConfiguration, models.ApplicationConfiguration]
}

func NewapplicationConfigurationHandler(s svcContracts.NestedService[models.ApplicationConfiguration, models.ApplicationConfiguration]) ApplicationConfigurationHandler {
	return &applicationConfigurationHandler{
		NestedGenericHandler: hnd.NewNestedGenericHandler[models.ApplicationConfiguration](IDInt64Parameters, IDInt64Parameters, s),
	}
}
