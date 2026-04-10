package models

import "github.com/brunojet/go-infra-backend/pkg/ports/errors"

// Centralized model errors for demoapp/models

var (
	errInvalidStageTransition = errors.NewBusinessRuleError("invalid stage transition")
	errInvalidStage           = errors.NewBusinessRuleError("invalid stage")
)
