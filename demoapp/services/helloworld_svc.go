package services

import (
	"github.com/brunojet/go-infra-backend/demoapp/repositories"
	"github.com/brunojet/go-infra-backend/internal/ports/services"
	"github.com/brunojet/go-infra-backend/internal/utils"
	rpocts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	svccts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
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
	svccts.Service[HelloWorldDTO, HelloWorldDTO, HelloWorldDTO, repositories.HelloWorld]
}

type helloWorldService struct {
	svccts.Service[HelloWorldDTO, HelloWorldDTO, HelloWorldDTO, repositories.HelloWorld]
}

func NewHelloWorldService(repo rpocts.Repository[repositories.HelloWorld]) HelloWorldService {
	return &helloWorldService{
		Service: services.NewServiceImpl(repo, helloWorldMapper{}),
	}
}
