package services

import (
	"context"

	"github.com/brunojet/go-infra-backend/demoapp/core/dtos"
	"github.com/brunojet/go-infra-backend/demoapp/core/models"
	"github.com/brunojet/go-infra-backend/demoapp/core/repository/ports"
)

type HelloWorldService struct {
	repo ports.CRUDRepository[models.HelloWorld]
}

func NewHelloWorldService(repo ports.CRUDRepository[models.HelloWorld]) *HelloWorldService {
	return &HelloWorldService{repo: repo}
}

func (s *HelloWorldService) List(ctx context.Context) ([]dtos.HelloWorldResponse, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	responses := make([]dtos.HelloWorldResponse, len(items))
	for i, item := range items {
		responses[i] = dtos.ToHelloWorldResponse(item)
	}
	return responses, nil
}

func (s *HelloWorldService) Get(ctx context.Context, id uint) (*dtos.HelloWorldResponse, error) {
	item, err := s.repo.Get(ctx, id)
	if err != nil || item == nil {
		return nil, err
	}
	resp := dtos.ToHelloWorldResponse(*item)
	return &resp, nil
}

func (s *HelloWorldService) Create(ctx context.Context, req dtos.CreateHelloWorldRequest) (*dtos.HelloWorldResponse, error) {
	entity := models.HelloWorld{
		Message:  req.Message,
		Language: req.Language,
	}
	created, err := s.repo.Create(ctx, &entity)
	if err != nil {
		return nil, err
	}
	resp := dtos.ToHelloWorldResponse(*created)
	return &resp, nil
}

func (s *HelloWorldService) Patch(ctx context.Context, id uint, req dtos.PatchHelloWorldRequest) (*dtos.HelloWorldResponse, error) {
	updates := make(map[string]any)
	if req.Message != nil {
		updates["message"] = *req.Message
	}
	if req.Language != nil {
		updates["language"] = *req.Language
	}
	updated, err := s.repo.Patch(ctx, id, updates)
	if err != nil || updated == nil {
		return nil, err
	}
	resp := dtos.ToHelloWorldResponse(*updated)
	return &resp, nil
}

func (s *HelloWorldService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
