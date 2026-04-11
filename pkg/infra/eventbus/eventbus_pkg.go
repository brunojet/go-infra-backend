package eventbus

// Package eventbus provides exported contracts used across the repository.
// This file exists to provide package-level documentation and to make the
// package discoverable by tooling.

// Re-export commonly used contract types from the contracts subpackage for
// convenient imports such as "github.com/brunojet/go-infra-backend/pkg/infra/eventbus".
// Keeping the minimal surface here prevents import cycles while providing a
// stable package entrypoint.

import (
	"github.com/brunojet/go-infra-backend/internal/infra/eventbus"
	"github.com/brunojet/go-infra-backend/pkg/infra/eventbus/contracts"
)

// Useful aliases for easier consumption by other packages.
type (
	Handler     = contracts.Handler
	HandlerName = contracts.HandlerName
	EventBus    = contracts.EventBus
)

// NewEventBus constructs a new EventBus implementation and returns it as the
// exported contract type.
func NewEventBus() EventBus {
	return eventbus.NewEventBus()
}
