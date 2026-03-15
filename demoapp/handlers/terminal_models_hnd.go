package handlers

import (
	"github.com/brunojet/go-infra-backend/demoapp/models"
	hnd "github.com/brunojet/go-infra-backend/pkg/ports/handlers"
	hndcontracts "github.com/brunojet/go-infra-backend/pkg/ports/handlers/contracts"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

type TerminalModelHandler interface {
	hndcontracts.GenericHandler[models.TerminalModel, models.TerminalModel]
}

type terminalModelHandler struct {
	hndcontracts.GenericHandler[models.TerminalModel, models.TerminalModel]
}

func NewTerminalModelsHandler(s svcContracts.Service[models.TerminalModel, models.TerminalModel]) TerminalModelHandler {
	return &terminalModelHandler{
		GenericHandler: hnd.NewGenericHandler[models.TerminalModel](IDInt64Parameters, s),
	}
}

type TerminalModelConfigurationHandler interface {
	hndcontracts.NestedGenericHandler[models.TerminalModelConfiguration, models.TerminalModelConfiguration]
}

type terminalModelConfigurationHandler struct {
	hndcontracts.NestedGenericHandler[models.TerminalModelConfiguration, models.TerminalModelConfiguration]
}

func NewTerminalModelConfigurationHandler(s svcContracts.NestedService[models.TerminalModelConfiguration, models.TerminalModelConfiguration]) TerminalModelConfigurationHandler {
	return &terminalModelConfigurationHandler{
		NestedGenericHandler: hnd.NewNestedGenericHandler[models.TerminalModelConfiguration](IDInt64Parameters, IDBase64Parameters, s),
	}
}
