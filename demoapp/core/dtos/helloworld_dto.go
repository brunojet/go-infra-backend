package dtos

import (
	dtoports "github.com/brunojet/go-infra-backend/demoapp/core/dtos/ports"
	"github.com/brunojet/go-infra-backend/demoapp/core/models"
)

type CreateHelloWorldRequest struct {
	Message  string `json:"mensagem"`
	Language string `json:"lingua"`
}

type PatchHelloWorldRequest struct {
	Message  *string `json:"mensagem"`
	Language *string `json:"lingua"`
}

type HelloWorldResponse struct {
	ID       uint   `json:"id"`
	Message  string `json:"mensagem"`
	Language string `json:"lingua"`
	dtoports.AuditDTO
}

func ToHelloWorldResponse(a models.HelloWorld) HelloWorldResponse {
	return HelloWorldResponse{
		ID:       a.ID,
		Message:  a.Message,
		Language: a.Language,
		AuditDTO: dtoports.ToAuditResponse(a.Audit),
	}
}
