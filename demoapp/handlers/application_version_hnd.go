package handlers

import (
	"github.com/brunojet/go-infra-backend/demoapp/models"
	hnd "github.com/brunojet/go-infra-backend/internal/ports/handlers"
	hndcontracts "github.com/brunojet/go-infra-backend/pkg/ports/handlers/contracts"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

type ApplicationVersionHandler interface {
	hndcontracts.NestedGenericHandler[models.ApplicationVersion, models.ApplicationVersion]
}

type applicationVersionHandler struct {
	hndcontracts.NestedGenericHandler[models.ApplicationVersion, models.ApplicationVersion]
}

func NewApplicationVersionHandler(s svcContracts.NestedService[models.ApplicationVersion, models.ApplicationVersion]) ApplicationVersionHandler {
	return &applicationVersionHandler{
		NestedGenericHandler: hnd.NewNestedGenericHandler[models.ApplicationVersion](IDInt64Parameters, IDBase64Parameters, s),
	}
}
