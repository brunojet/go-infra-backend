package services

import (
	bffsvccts "github.com/brunojet/go-infra-backend/pkg/ports/bff/services/contracts"
)

// HelloWorldBffE is the upstream DTO for the helloworld resource.
// CE = RE = UE because demoapp's API is symmetric (same shape for create, read, update).
type HelloWorldBffE struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func (HelloWorldBffE) ResourceName() string { return "hello-worlds" }

// HelloWorldBffDTO is the domain DTO exposed to handlers.
type HelloWorldBffDTO struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// HelloWorldBffMapper implements BffServiceMapper for the helloworld resource.
// Mapping is a trivial pass-through because the domain and upstream shapes are identical.
type HelloWorldBffMapper struct{}

var _ bffsvccts.BffServiceMapper[
	HelloWorldBffDTO, HelloWorldBffDTO, HelloWorldBffDTO,
	HelloWorldBffE, HelloWorldBffE, HelloWorldBffE,
] = HelloWorldBffMapper{}

func (HelloWorldBffMapper) ToUpstreamPost(dto HelloWorldBffDTO, up *HelloWorldBffE) error {
	up.Message = dto.Message
	return nil
}

func (HelloWorldBffMapper) ToUpstreamPatch(dto HelloWorldBffDTO, up *HelloWorldBffE) error {
	up.Message = dto.Message
	return nil
}

func (HelloWorldBffMapper) ToDomainDTO(up *HelloWorldBffE, dto *HelloWorldBffDTO) error {
	dto.ID = up.ID
	dto.Message = up.Message
	return nil
}

func (HelloWorldBffMapper) GetUpstreamID(id string) (string, error) {
	return id, nil
}

func (HelloWorldBffMapper) ApplyQueryScopes(scopes map[string]any) (map[string]any, error) {
	return scopes, nil
}

// ExtractUpstreamTotal is not used for the routes exposed by demobff (no List).
func (HelloWorldBffMapper) ExtractUpstreamTotal(upstream []HelloWorldBffE) int64 {
	return int64(len(upstream))
}

// ExtractUpstreamError passes through the error as-is; helloworld has no
// custom upstream error contract.
func (HelloWorldBffMapper) ExtractUpstreamError(err error) error {
	return err
}
