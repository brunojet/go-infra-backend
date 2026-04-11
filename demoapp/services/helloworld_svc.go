package services

import (
	"github.com/brunojet/go-infra-backend/demoapp/repositories"
	rpocts "github.com/brunojet/go-infra-backend/pkg/ports/backend/repositories/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/backend/services"
	svccts "github.com/brunojet/go-infra-backend/pkg/ports/backend/services/contracts"
	"github.com/brunojet/go-infra-backend/pkg/utils"
)

type HelloWorldDTO struct {
	ID      string
	Message string
}

// helloWorldMapper implements contracts.ServiceMapper[HelloWorldDTO, helloWorldRepo.HelloWorld]
type helloWorldMapper struct{}

func (helloWorldMapper) GetModelKey(id string) (map[string]any, error) {
	return map[string]any{helloWorldIDKey: id}, nil
}

func (helloWorldMapper) ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error) {
	return queryScopes, nil
}

func (helloWorldMapper) ToPostModel(dto HelloWorldDTO, model *repositories.HelloWorld) error {
	if model == nil {
		return errMapperNilModel
	}
	model.Message = utils.ToNullString(dto.Message)
	return nil
}

func (helloWorldMapper) ToPatchModel(dto HelloWorldDTO, model *repositories.HelloWorld) error {
	if model == nil {
		return errMapperNilModel
	}
	model.Message = utils.ToNullString(dto.Message)
	return nil
}

func (helloWorldMapper) ToDTO(mdl *repositories.HelloWorld, dto *HelloWorldDTO) error {
	if mdl == nil || dto == nil {
		return errMapperNilModel
	}
	dto.ID = mdl.ID
	dto.Message = utils.FromNullString(mdl.Message)
	return nil
}

type HelloWorldService interface {
	svccts.Service[HelloWorldDTO, HelloWorldDTO, HelloWorldDTO]
}

type helloWorldService struct {
	svccts.Service[HelloWorldDTO, HelloWorldDTO, HelloWorldDTO]
}

func NewHelloWorldService(repo rpocts.Repository[repositories.HelloWorld]) HelloWorldService {
	return &helloWorldService{
		Service: services.NewServiceImpl(repo, helloWorldMapper{}),
	}
}
