package contracts

import internalcontracts "github.com/brunojet/go-infra-backend/internal/ports/services/contracts"

type ServiceMapper[D any, E any] = internalcontracts.ServiceMapper[D, E]

type Service[D any, E any] = internalcontracts.Service[D, E]
