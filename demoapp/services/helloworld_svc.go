package services

import (
	helloWorldRepo "github.com/brunojet/go-infra-backend/demoapp/repositories"
	repoContracts "github.com/brunojet/go-infra-backend/internal/ports/repositories/contracts"
	svc "github.com/brunojet/go-infra-backend/internal/ports/services"
	svcContracts "github.com/brunojet/go-infra-backend/internal/ports/services/contracts"
)

type HelloWorldDTO struct {
	ID      string
	Message string
}

// helloWorldMapper implements contracts.ServiceMapper[HelloWorldDTO, helloWorldRepo.HelloWorld]
type helloWorldMapper struct{}

func (helloWorldMapper) GetModelKey(id string) (map[string]any, error) {
	return map[string]any{"id": id}, nil
}

func (helloWorldMapper) ToModel(dto *HelloWorldDTO) (helloWorldRepo.HelloWorld, error) {
	if dto == nil {
		return helloWorldRepo.HelloWorld{}, nil
	}
	var mdl helloWorldRepo.HelloWorld
	mdl.Message = svc.ToNullString(dto.Message)
	return mdl, nil
}

func (helloWorldMapper) ToDTO(mdl *helloWorldRepo.HelloWorld, dto *HelloWorldDTO) {
	if dto == nil {
		return
	}
	if mdl == nil {
		*dto = HelloWorldDTO{}
		return
	}
	dto.ID = mdl.ID
	dto.Message = svc.FromNullString(mdl.Message)
}

type HelloWorldService struct {
	svcContracts.Service[HelloWorldDTO, helloWorldRepo.HelloWorld]
}

func NewHelloWorldService(repo repoContracts.Repository[helloWorldRepo.HelloWorld]) svcContracts.Service[HelloWorldDTO, helloWorldRepo.HelloWorld] {
	return &HelloWorldService{
		Service: svc.NewServiceImpl(repo, helloWorldMapper{}),
	}
}
