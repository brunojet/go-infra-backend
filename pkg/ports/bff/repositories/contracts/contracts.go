package contracts

import "context"

// BffEntity is implemented by upstream/legacy DTOs used as the "entity"
// in BFF-backed repositories. ResourceName returns the resource name used
// to derive the upstream route (e.g., "applications", "terminal-models").
//
// This is the BFF counterpart of repositories/contracts.Entity (TableName).
type BffEntity interface {
	ResourceName() string
}

// QueryParams carries free-form filter scopes that the mapper translates
// into URL query parameters before the BffRepository issues the request.
type QueryParams struct {
	Scopes map[string]any
}

// BffListParams carries pagination and filtering for BFF repository queries.
// Unlike the database ListParams, Size is explicit because upstream APIs express
// pagination as page+size, not as an internal database offset.
type BffListParams struct {
	QueryParams
	Page    int
	Size    int
	OrderBy string
	Order   string
}

// BffRepository is the port for upstream/BFF-backed adapters.
// It is the parallel counterpart of repositories/contracts.Repository[E Entity]
// but operates over a remote upstream API instead of a local database.
//
// Upstream APIs (e.g. ServiceNow) are asymmetric — request and response shapes
// differ per verb, so three distinct types are required:
//   - CE: upstream create/write DTO (POST body)
//   - RE: upstream read DTO (GET / List response)
//   - UE: upstream update DTO (PATCH body)
//
// All three must implement BffEntity so the adapter can derive the resource name.
// Inputs (upstream) are passed by value; outputs (downstream) are pointers.
// IDs are plain strings (URL path segments) rather than map[string]any DB keys.
type BffRepository[CE, RE, UE BffEntity] interface {
	Create(ctx context.Context, upstream CE, downstream *RE) error
	GetByID(ctx context.Context, id string, downstream *RE) error
	List(ctx context.Context, params BffListParams, downstream *[]RE) error
	Update(ctx context.Context, id string, upstream UE, downstream *RE) error
	Delete(ctx context.Context, id string) error
}

// BffNestedRepository extends BffRepository for resources that are addressed
// under a parent path (e.g. /applications/{parentID}/profiles).
type BffNestedRepository[CE, RE, UE BffEntity] interface {
	BffRepository[CE, RE, UE]
	CreateNested(ctx context.Context, parentID string, upstream CE, downstream *RE) error
	ListNested(ctx context.Context, parentID string, params BffListParams, downstream *[]RE) error
}
