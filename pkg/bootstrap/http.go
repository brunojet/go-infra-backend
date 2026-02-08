package bootstrap

import (
	internalbootstrap "github.com/brunojet/go-infra-backend/internal/bootstrap"
	"github.com/brunojet/go-infra-backend/pkg/bootstrap/contracts"
)

type HttpServer = internalbootstrap.HttpServer

func NewHttpServerWithObservability(sm contracts.ShutdownManager) *HttpServer {
	return internalbootstrap.NewHttpServerWithObservability(sm)
}
