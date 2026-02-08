package contracts

import internalcontracts "github.com/brunojet/go-infra-backend/internal/ports/repositories/contracts"

type ListParams = internalcontracts.ListParams

type Entity = internalcontracts.Entity

type Repository[E Entity] = internalcontracts.Repository[E]
