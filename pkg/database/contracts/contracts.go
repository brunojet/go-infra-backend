package contracts

import internalcontracts "github.com/brunojet/go-infra-backend/internal/database/contracts"

// Minimal public database contracts re-exported from internal.
//
// Transitional strategy:
// - external apps import pkg/database/*
// - pkg/database/* delegates to internal/database/*
// - internal/database/* remains the implementation source for now

type DatabaseAdapter = internalcontracts.DatabaseAdapter

type DatabaseManager = internalcontracts.DatabaseManager
