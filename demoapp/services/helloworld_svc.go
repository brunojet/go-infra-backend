package services

import (
	helloWorldRepo "github.com/brunojet/go-infra-backend/demoapp/repositories"
	svc "github.com/brunojet/go-infra-backend/internal/ports/services"
	"github.com/brunojet/go-infra-backend/internal/utils"
	repoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

const helloWorldIDKey = "id"

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

func (helloWorldMapper) ToModel(dto *HelloWorldDTO, model *helloWorldRepo.HelloWorld) {
	model.Message = utils.ToNullString(dto.Message)
}

func (helloWorldMapper) ToDTO(mdl *helloWorldRepo.HelloWorld, dto *HelloWorldDTO) {
	dto.ID = mdl.ID
	dto.Message = utils.FromNullString(mdl.Message)
}

type HelloWorldService struct {
	svcContracts.Service[HelloWorldDTO, helloWorldRepo.HelloWorld]
}

func NewHelloWorldService(repo repoContracts.Repository[helloWorldRepo.HelloWorld]) svcContracts.Service[HelloWorldDTO, helloWorldRepo.HelloWorld] {
	return &HelloWorldService{
		Service: svc.NewServiceImpl(repo, helloWorldMapper{}),
	}
}
