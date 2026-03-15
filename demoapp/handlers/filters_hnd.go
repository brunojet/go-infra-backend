package handlers

import (
	"github.com/brunojet/go-infra-backend/demoapp/models"
	hnd "github.com/brunojet/go-infra-backend/pkg/ports/handlers"
	hndcontracts "github.com/brunojet/go-infra-backend/pkg/ports/handlers/contracts"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

type FilterTypeHandler interface {
	hndcontracts.GenericHandler[models.FilterType, models.FilterType]
}

type filterTypeHandler struct {
	hndcontracts.GenericHandler[models.FilterType, models.FilterType]
}

func NewFiltersTypeHandler(s svcContracts.Service[models.FilterType, models.FilterType]) FilterTypeHandler {
	return &filterTypeHandler{
		GenericHandler: hnd.NewGenericHandler[models.FilterType](IDInt64Parameters, s),
	}
}

type FilterHandler interface {
	hndcontracts.NestedGenericHandler[models.Filter, models.Filter]
}

type filterHandler struct {
	hndcontracts.NestedGenericHandler[models.Filter, models.Filter]
}

func NewFiltersHandler(s svcContracts.NestedService[models.Filter, models.Filter]) FilterHandler {
	return &filterHandler{
		NestedGenericHandler: hnd.NewNestedGenericHandler[models.Filter](IDInt64Parameters, IDBase64Parameters, s),
	}
}
