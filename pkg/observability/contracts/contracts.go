package contracts

import internalcontracts "github.com/brunojet/go-infra-backend/internal/observability/contracts"

// Minimal public contracts re-exported from internal.
//
// Transitional strategy:
// - External apps import pkg/*
// - pkg/* delegates to internal/*
// - internal/* remains the single implementation source for now

type ShutdownFunc = internalcontracts.ShutdownFunc

type ShutdownRegister = internalcontracts.ShutdownRegister
