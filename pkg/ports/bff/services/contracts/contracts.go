package contracts

import svccts "github.com/brunojet/go-infra-backend/pkg/ports/backend/services/contracts"

// BffServiceMapper translates between domain DTOs (C/R/U) and the asymmetric
// upstream API DTOs (CE/RE/UE). Upstream APIs such as ServiceNow have different
// shapes per verb — CE for create, RE for read, UE for update — so all three
// are explicit type parameters.
//
// This is the BFF counterpart of ServiceMapper[C,R,U,E Entity].
type BffServiceMapper[C, R, U, CE, RE, UE any] interface {
	ToUpstreamPost(dto C, upstream *CE) error
	ToUpstreamPatch(dto U, upstream *UE) error
	ToDomainDTO(upstream *RE, dto *R) error
	// GetUpstreamID maps a domain resource ID to the path segment used
	// by the upstream API (e.g., numeric string → UUID, or direct passthrough).
	GetUpstreamID(id string) (string, error)
	ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error)
	// ExtractUpstreamTotal reads the total item count from a slice of upstream
	// read DTOs. Only the mapper knows the RE shape and where the count lives
	// (e.g. a metadata field in the first element), so extraction is delegated here
	// rather than handled by the repository or transport layers.
	ExtractUpstreamTotal(upstream []RE) int64
}

// BffNestedServiceMapper extends BffServiceMapper with parent-awareness,
// mirroring NestedServiceMapper for the HTTP/BFF adapter path.
type BffNestedServiceMapper[C, R, U, CE, RE, UE any] interface {
	BffServiceMapper[C, R, U, CE, RE, UE]
	// GetUpstreamParentID maps a domain parent ID to the path segment
	// used by the upstream API for the nested route (e.g., /parent/{id}/children).
	GetUpstreamParentID(parentID string) (string, error)
	ApplyParentQueryScopes(parentID string, queryScopes map[string]any) (map[string]any, error)
	// ApplyUpstreamParentScopes sets the parent reference on the upstream CE DTO
	// before a nested Create, for upstream APIs that require it in the request body.
	ApplyUpstreamParentScopes(parentID string, upstream *CE) error
}

// BffService is the service contract exposed to handlers for BFF-backed resources.
// Handlers remain identical regardless of whether the underlying adapter hits a
// database (serviceImpl) or an upstream HTTP API (bffServiceImpl).
type BffService[C, R, U any] interface {
	svccts.Service[C, R, U]
}

// BffNestedService is the nested-resource variant of BffService.
// Embeds NestedService directly — CreateNested and ListNested are inherited.
type BffNestedService[C, R, U any] interface {
	svccts.NestedService[C, R, U]
}
