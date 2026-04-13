package repositories

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/brunojet/go-infra-backend/debugassert"
	bffstreams "github.com/brunojet/go-infra-backend/internal/infra/bffclient/streams"
	porterrors "github.com/brunojet/go-infra-backend/internal/infra/errors"
	bffcts "github.com/brunojet/go-infra-backend/pkg/infra/bffclient/contracts"
	bffrpocts "github.com/brunojet/go-infra-backend/pkg/ports/bff/repositories/contracts"
)

// httpClient is the concrete BffClient instantiation used by all BFF repositories.
// Repositories are always backed by the HTTP adapter; protocol selection happens
// at the BffRepository boundary, not inside the transport.
type httpClient = bffcts.BffClient[bffcts.BffHttpRequestStream, bffcts.BffHttpResponseStream]

// Compile-time interface satisfaction checks.
var (
	_ bffrpocts.BffRepository[bffrpocts.BffEntity, bffrpocts.BffEntity, bffrpocts.BffEntity] = (*bffRepositoryImpl[bffrpocts.BffEntity, bffrpocts.BffEntity, bffrpocts.BffEntity])(nil)

	_ bffrpocts.BffNestedRepository[bffrpocts.BffEntity, bffrpocts.BffEntity, bffrpocts.BffEntity] = (*bffNestedRepositoryImpl[bffrpocts.BffEntity, bffrpocts.BffEntity, bffrpocts.BffEntity])(nil)
)

// bffRepositoryImpl is the concrete implementation of BffRepository backed by
// a BffClient (net/http + backoff + circuit-breaker).
//
// Inputs (upstream) are always passed by value; outputs (downstream) are
// always pointers — matching the BffRepository contract and the asymmetric
// shape of upstream APIs.
type bffRepositoryImpl[CE, RE, UE bffrpocts.BffEntity] struct {
	client httpClient
}

// NewBffRepository creates a BffRepository that delegates all HTTP calls to
// the provided BffClient.
//
// The resource path for each verb is derived at runtime from the upstream
// entity's ResourceName() method, keeping the repository free of hard-coded
// path constants.
func NewBffRepository[CE, RE, UE bffrpocts.BffEntity](client httpClient) bffrpocts.BffRepository[CE, RE, UE] {
	debugassert.Assert(client != nil, "NewBffRepository: client is nil")
	return &bffRepositoryImpl[CE, RE, UE]{client: client}
}

// resourcePath returns the URL path segment for RE (e.g. "terminal-models").
// It instantiates a zero-value RE and calls ResourceName(), which is safe
// because BffEntity is constrained to a value-receiver interface.
func resourcePath[E bffrpocts.BffEntity]() string {
	var zero E
	return zero.ResourceName()
}

func (r *bffRepositoryImpl[CE, RE, UE]) Create(ctx context.Context, upstream CE, downstream *RE) error {
	req, err := bffstreams.NewJsonRequest(http.MethodPost, upstream.ResourceName(), nil, upstream)
	if err != nil {
		return err
	}
	resp := &bffstreams.JsonResponseStream[RE]{}
	if err := r.client.Emit(ctx, req, resp); err != nil {
		return err
	}
	if err := checkStatus(resp); err != nil {
		return err
	}
	*downstream = resp.Value
	return nil
}

func (r *bffRepositoryImpl[CE, RE, UE]) GetByID(ctx context.Context, id string, downstream *RE) error {
	path := fmt.Sprintf("%s/%s", resourcePath[RE](), id)
	req := bffstreams.NewJsonNoBodyRequest(http.MethodGet, path, nil)
	resp := &bffstreams.JsonResponseStream[RE]{}
	if err := r.client.Emit(ctx, req, resp); err != nil {
		return err
	}
	if err := checkStatus(resp); err != nil {
		return err
	}
	*downstream = resp.Value
	return nil
}

func (r *bffRepositoryImpl[CE, RE, UE]) List(ctx context.Context, params bffrpocts.BffListParams, downstream *[]RE) error {
	req := bffstreams.NewJsonNoBodyRequest(http.MethodGet, resourcePath[RE](), toQueryParams(params))
	resp := &bffstreams.JsonResponseStream[[]RE]{}
	if err := r.client.Emit(ctx, req, resp); err != nil {
		return err
	}
	if err := checkStatus(resp); err != nil {
		return err
	}
	*downstream = resp.Value
	return nil
}

func (r *bffRepositoryImpl[CE, RE, UE]) Update(ctx context.Context, id string, upstream UE, downstream *RE) error {
	path := fmt.Sprintf("%s/%s", upstream.ResourceName(), id)
	req, err := bffstreams.NewJsonRequest(http.MethodPatch, path, nil, upstream)
	if err != nil {
		return err
	}
	resp := &bffstreams.JsonResponseStream[RE]{}
	if err := r.client.Emit(ctx, req, resp); err != nil {
		return err
	}
	if err := checkStatus(resp); err != nil {
		return err
	}
	*downstream = resp.Value
	return nil
}

func (r *bffRepositoryImpl[CE, RE, UE]) Delete(ctx context.Context, id string) error {
	path := fmt.Sprintf("%s/%s", resourcePath[RE](), id)
	resp := bffstreams.NewNoBodyResponseStream()
	if err := r.client.Emit(ctx, bffstreams.NewJsonNoBodyRequest(http.MethodDelete, path, nil), resp); err != nil {
		return err
	}
	return checkStatus(resp)
}

// ---------------------------------------------------------------------------
// BffNestedRepository
// ---------------------------------------------------------------------------

// bffNestedRepositoryImpl extends bffRepositoryImpl with parent-scoped
// Create and List operations for resources addressed under a parent path
// (e.g. /applications/{parentID}/profiles).
type bffNestedRepositoryImpl[CE, RE, UE bffrpocts.BffEntity] struct {
	bffRepositoryImpl[CE, RE, UE]
	parentResourceName string
}

// NewBffNestedRepository creates a BffNestedRepository for child resources
// that live under a parent path.
//
// parentResourceName is the URL segment of the parent resource
// (e.g. "applications" for /applications/{parentID}/profiles).
func NewBffNestedRepository[CE, RE, UE bffrpocts.BffEntity](
	client httpClient,
	parentResourceName string,
) bffrpocts.BffNestedRepository[CE, RE, UE] {
	debugassert.Assert(client != nil, "NewBffNestedRepository: client is nil")
	debugassert.Assert(parentResourceName != "", "NewBffNestedRepository: parentResourceName is empty")
	return &bffNestedRepositoryImpl[CE, RE, UE]{
		bffRepositoryImpl:  bffRepositoryImpl[CE, RE, UE]{client: client},
		parentResourceName: parentResourceName,
	}
}

// nestedBasePath builds the parent-scoped resource path.
// e.g. parentResourceName="applications", parentID="42", resourcePath[RE]()="profiles"
//
//	→ "applications/42/profiles"
func (r *bffNestedRepositoryImpl[CE, RE, UE]) nestedBasePath(parentID string) string {
	return fmt.Sprintf("%s/%s/%s", r.parentResourceName, parentID, resourcePath[RE]())
}

func (r *bffNestedRepositoryImpl[CE, RE, UE]) CreateNested(ctx context.Context, parentID string, upstream CE, downstream *RE) error {
	req, err := bffstreams.NewJsonRequest(http.MethodPost, r.nestedBasePath(parentID), nil, upstream)
	if err != nil {
		return err
	}
	resp := &bffstreams.JsonResponseStream[RE]{}
	if err := r.client.Emit(ctx, req, resp); err != nil {
		return err
	}
	if err := checkStatus(resp); err != nil {
		return err
	}
	*downstream = resp.Value
	return nil
}

func (r *bffNestedRepositoryImpl[CE, RE, UE]) ListNested(ctx context.Context, parentID string, params bffrpocts.BffListParams, downstream *[]RE) error {
	req := bffstreams.NewJsonNoBodyRequest(http.MethodGet, r.nestedBasePath(parentID), toQueryParams(params))
	resp := &bffstreams.JsonResponseStream[[]RE]{}
	if err := r.client.Emit(ctx, req, resp); err != nil {
		return err
	}
	if err := checkStatus(resp); err != nil {
		return err
	}
	*downstream = resp.Value
	return nil
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// checkStatus returns an ErrHTTP when the response stream carries a non-2xx
// status code. Emit only returns transport errors; HTTP-level errors are
// communicated exclusively through the response stream's StatusCode().
func checkStatus(resp bffcts.BffHttpResponseStream) error {
	if resp.StatusCode() >= http.StatusBadRequest {
		return porterrors.NewHTTPError(resp.StatusCode())
	}
	return nil
}

// toQueryParams converts BffListParams to the flat string map expected by
// BffClient. The Scopes map should already contain upstream-translated filter
// param names (applied by the service mapper via ApplyQueryScopes before the
// repository is called). Pagination params are added with standard names.
func toQueryParams(params bffrpocts.BffListParams) map[string]string {
	q := make(map[string]string, len(params.Scopes)+4)
	for k, v := range params.Scopes {
		q[k] = fmt.Sprintf("%v", v)
	}
	if params.Page > 0 {
		q["page"] = strconv.Itoa(params.Page)
	}
	if params.Size > 0 {
		q["size"] = strconv.Itoa(params.Size)
	}
	if params.OrderBy != "" {
		q["orderBy"] = params.OrderBy
	}
	if params.Order != "" {
		q["order"] = params.Order
	}
	return q
}
