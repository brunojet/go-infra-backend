package repositories

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/brunojet/go-infra-backend/pkg/infra/bffclient"
	bffcts "github.com/brunojet/go-infra-backend/pkg/infra/bffclient/contracts"
	bffrpocts "github.com/brunojet/go-infra-backend/pkg/ports/bff/repositories/contracts"
)

// ---------------------------------------------------------------------------
// test fixtures
// ---------------------------------------------------------------------------

// itemEntity is an upstream DTO that satisfies BffEntity.
type itemEntity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (itemEntity) ResourceName() string { return "items" }

// childEntity lives under a parent resource.
type childEntity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (childEntity) ResourceName() string { return "children" }

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func newClient(t *testing.T, serverURL string) httpClient {
	t.Helper()
	cfg := bffcts.BffClientConfig{BaseURL: serverURL}
	client, _, err := bffclient.NewNetHttpAdapter(cfg, nil)
	require.NoError(t, err)
	return client
}

func jsonHandler(t *testing.T, status int, body any) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != nil {
			require.NoError(t, json.NewEncoder(w).Encode(body))
		}
	}
}

// ---------------------------------------------------------------------------
// BffRepository — Create
// ---------------------------------------------------------------------------

func TestBffRepository_Create_Success(t *testing.T) {
	want := itemEntity{ID: "1", Name: "created"}
	var capturedMethod, capturedPath string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		capturedPath = r.URL.Path
		jsonHandler(t, http.StatusCreated, want)(w, r)
	}))
	t.Cleanup(srv.Close)

	repo := NewBffRepository[itemEntity, itemEntity, itemEntity](newClient(t, srv.URL))

	var got itemEntity
	err := repo.Create(context.Background(), itemEntity{Name: "created"}, &got)

	require.NoError(t, err)
	assert.Equal(t, http.MethodPost, capturedMethod)
	assert.Equal(t, "/items", capturedPath)
	assert.Equal(t, want, got)
}

// ---------------------------------------------------------------------------
// BffRepository — GetByID
// ---------------------------------------------------------------------------

func TestBffRepository_GetByID_Success(t *testing.T) {
	want := itemEntity{ID: "42", Name: "found"}
	var capturedPath string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		jsonHandler(t, http.StatusOK, want)(w, r)
	}))
	t.Cleanup(srv.Close)

	repo := NewBffRepository[itemEntity, itemEntity, itemEntity](newClient(t, srv.URL))

	var got itemEntity
	err := repo.GetByID(context.Background(), "42", &got)

	require.NoError(t, err)
	assert.Equal(t, "/items/42", capturedPath)
	assert.Equal(t, want, got)
}

func TestBffRepository_GetByID_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	t.Cleanup(srv.Close)

	repo := NewBffRepository[itemEntity, itemEntity, itemEntity](newClient(t, srv.URL))

	var got itemEntity
	err := repo.GetByID(context.Background(), "99", &got)

	require.Error(t, err)
	assert.True(t, bffclient.IsNotFound(err))
}

// ---------------------------------------------------------------------------
// BffRepository — List
// ---------------------------------------------------------------------------

func TestBffRepository_List_Success(t *testing.T) {
	items := []itemEntity{{ID: "1"}, {ID: "2"}}
	var capturedQuery string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		jsonHandler(t, http.StatusOK, items)(w, r)
	}))
	t.Cleanup(srv.Close)

	repo := NewBffRepository[itemEntity, itemEntity, itemEntity](newClient(t, srv.URL))

	params := bffrpocts.BffListParams{
		Page:    1,
		Size:    10,
		OrderBy: "created_at",
		Order:   "asc",
		QueryParams: bffrpocts.QueryParams{
			Scopes: map[string]any{"status": "active"},
		},
	}

	var got []itemEntity
	err := repo.List(context.Background(), params, &got)

	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Contains(t, capturedQuery, "page=1")
	assert.Contains(t, capturedQuery, "size=10")
	assert.Contains(t, capturedQuery, "orderBy=created_at")
	assert.Contains(t, capturedQuery, "order=asc")
	assert.Contains(t, capturedQuery, "status=active")
}

// ---------------------------------------------------------------------------
// BffRepository — Update
// ---------------------------------------------------------------------------

func TestBffRepository_Update_Success(t *testing.T) {
	updated := itemEntity{ID: "5", Name: "updated"}
	var capturedMethod, capturedPath string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		capturedPath = r.URL.Path
		jsonHandler(t, http.StatusOK, updated)(w, r)
	}))
	t.Cleanup(srv.Close)

	repo := NewBffRepository[itemEntity, itemEntity, itemEntity](newClient(t, srv.URL))

	var got itemEntity
	err := repo.Update(context.Background(), "5", itemEntity{Name: "updated"}, &got)

	require.NoError(t, err)
	assert.Equal(t, http.MethodPatch, capturedMethod)
	assert.Equal(t, "/items/5", capturedPath)
	assert.Equal(t, updated, got)
}

// ---------------------------------------------------------------------------
// BffRepository — Delete
// ---------------------------------------------------------------------------

func TestBffRepository_Delete_Success(t *testing.T) {
	var capturedPath, capturedMethod string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(srv.Close)

	repo := NewBffRepository[itemEntity, itemEntity, itemEntity](newClient(t, srv.URL))

	err := repo.Delete(context.Background(), "7")

	require.NoError(t, err)
	assert.Equal(t, http.MethodDelete, capturedMethod)
	assert.Equal(t, "/items/7", capturedPath)
}

// ---------------------------------------------------------------------------
// BffNestedRepository — CreateNested
// ---------------------------------------------------------------------------

func TestBffNestedRepository_CreateNested_Success(t *testing.T) {
	want := childEntity{ID: "c1", Name: "child"}
	var capturedPath string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		jsonHandler(t, http.StatusCreated, want)(w, r)
	}))
	t.Cleanup(srv.Close)

	repo := NewBffNestedRepository[childEntity, childEntity, childEntity](newClient(t, srv.URL), "parents")

	var got childEntity
	err := repo.CreateNested(context.Background(), "p42", childEntity{Name: "child"}, &got)

	require.NoError(t, err)
	assert.Equal(t, "/parents/p42/children", capturedPath)
	assert.Equal(t, want, got)
}

// ---------------------------------------------------------------------------
// BffNestedRepository — ListNested
// ---------------------------------------------------------------------------

func TestBffNestedRepository_ListNested_Success(t *testing.T) {
	children := []childEntity{{ID: "c1"}, {ID: "c2"}}
	var capturedPath string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		jsonHandler(t, http.StatusOK, children)(w, r)
	}))
	t.Cleanup(srv.Close)

	repo := NewBffNestedRepository[childEntity, childEntity, childEntity](newClient(t, srv.URL), "parents")

	var got []childEntity
	err := repo.ListNested(context.Background(), "p42", bffrpocts.BffListParams{Page: 1, Size: 5}, &got)

	require.NoError(t, err)
	assert.Equal(t, "/parents/p42/children", capturedPath)
	assert.Len(t, got, 2)
}

// ---------------------------------------------------------------------------
// BffNestedRepository — delegates flat ops to embedded bffRepositoryImpl
// ---------------------------------------------------------------------------

func TestBffNestedRepository_GetByID_DelegatesCorrectly(t *testing.T) {
	want := childEntity{ID: "c99", Name: "flat"}
	var capturedPath string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		jsonHandler(t, http.StatusOK, want)(w, r)
	}))
	t.Cleanup(srv.Close)

	repo := NewBffNestedRepository[childEntity, childEntity, childEntity](newClient(t, srv.URL), "parents")

	var got childEntity
	err := repo.GetByID(context.Background(), "c99", &got)

	require.NoError(t, err)
	assert.Equal(t, "/children/c99", capturedPath)
	assert.Equal(t, want, got)
}

// ---------------------------------------------------------------------------
// toQueryParams
// ---------------------------------------------------------------------------

func TestToQueryParams_EmptyParams(t *testing.T) {
	q := toQueryParams(bffrpocts.BffListParams{})
	assert.Empty(t, q)
}

func TestToQueryParams_AllFields(t *testing.T) {
	params := bffrpocts.BffListParams{
		Page:    2,
		Size:    20,
		OrderBy: "name",
		Order:   "desc",
		QueryParams: bffrpocts.QueryParams{
			Scopes: map[string]any{"active": true, "limit": 5},
		},
	}
	q := toQueryParams(params)

	assert.Equal(t, "2", q["page"])
	assert.Equal(t, "20", q["size"])
	assert.Equal(t, "name", q["orderBy"])
	assert.Equal(t, "desc", q["order"])
	assert.Equal(t, "true", q["active"])
	assert.Equal(t, "5", q["limit"])
}
