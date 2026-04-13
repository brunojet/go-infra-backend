package repositories

import (
	"github.com/brunojet/go-infra-backend/internal/ports/bff/repositories"
	bffcts "github.com/brunojet/go-infra-backend/pkg/infra/bffclient/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/bff/repositories/contracts"
)

// ---- Contracts ----

type (
	BffEntity                                           = contracts.BffEntity
	QueryParams                                         = contracts.QueryParams
	BffListParams                                       = contracts.BffListParams
	BffRepository[CE, RE, UE contracts.BffEntity]       = contracts.BffRepository[CE, RE, UE]
	BffNestedRepository[CE, RE, UE contracts.BffEntity] = contracts.BffNestedRepository[CE, RE, UE]
)

// ---- Constructors ----

// NewBffRepository creates a BffRepository[CE,RE,UE] backed by the given
// BffClient. Each method derives the upstream resource path at runtime via
// ResourceName() on the upstream entity type.
func NewBffRepository[CE, RE, UE contracts.BffEntity](client bffcts.BffClient[bffcts.BffHttpRequestStream, bffcts.BffHttpResponseStream]) contracts.BffRepository[CE, RE, UE] {
	return repositories.NewBffRepository[CE, RE, UE](client)
}

// NewBffNestedRepository creates a BffNestedRepository[CE,RE,UE] for child
// resources that live under a parent path
// (e.g. /parentResourceName/{parentID}/child-resource).
func NewBffNestedRepository[CE, RE, UE contracts.BffEntity](client bffcts.BffClient[bffcts.BffHttpRequestStream, bffcts.BffHttpResponseStream], parentResourceName string) contracts.BffNestedRepository[CE, RE, UE] {
	return repositories.NewBffNestedRepository[CE, RE, UE](client, parentResourceName)
}
