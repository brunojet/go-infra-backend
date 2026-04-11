package repositories

import "github.com/brunojet/go-infra-backend/pkg/ports/bff/repositories/contracts"

// ---- Contracts ----

type (
	BffEntity                                           = contracts.BffEntity
	QueryParams                                         = contracts.QueryParams
	BffListParams                                       = contracts.BffListParams
	BffRepository[CE, RE, UE contracts.BffEntity]       = contracts.BffRepository[CE, RE, UE]
	BffNestedRepository[CE, RE, UE contracts.BffEntity] = contracts.BffNestedRepository[CE, RE, UE]
)
