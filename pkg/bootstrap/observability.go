package bootstrap

import (
	internalbootstrap "github.com/brunojet/go-infra-backend/internal/bootstrap"
	"github.com/brunojet/go-infra-backend/pkg/bootstrap/contracts"
)

func InitLogger(sm contracts.ShutdownManager) error { return internalbootstrap.InitLogger(sm) }

func InitMetrics(sm contracts.ShutdownManager) error { return internalbootstrap.InitMetrics(sm) }

func InitTracing(sm contracts.ShutdownManager) error { return internalbootstrap.InitTracing(sm) }

func InitObservability(sm contracts.ShutdownManager) error {
	return internalbootstrap.InitObservability(sm)
}
