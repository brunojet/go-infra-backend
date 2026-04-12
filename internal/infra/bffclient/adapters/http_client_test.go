package adapters

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	bffstreams "github.com/brunojet/go-infra-backend/internal/infra/bffclient/streams"
	bffcts "github.com/brunojet/go-infra-backend/pkg/infra/bffclient/contracts"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

type testPayload struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func newAdapter(t *testing.T, serverURL string, opts ...func(*bffcts.BffClientConfig)) *netHttpAdapter {
	t.Helper()
	cfg := bffcts.BffClientConfig{BaseURL: serverURL}
	for _, o := range opts {
		o(&cfg)
	}
	a, err := NewNetHttpAdapter(cfg, nil) // nil → http.DefaultTransport
	require.NoError(t, err)
	return a
}

func jsonHandler(t *testing.T, status int, body any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != nil {
			require.NoError(t, json.NewEncoder(w).Encode(body))
		}
	}
}

// ---------------------------------------------------------------------------
// NewNetHttpAdapter — validation
// ---------------------------------------------------------------------------

func TestNewNetHttpAdapter_EmptyBaseURL(t *testing.T) {
	_, err := NewNetHttpAdapter(bffcts.BffClientConfig{}, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "BaseURL")
}

// ---------------------------------------------------------------------------
// Emit — POST (with body)
// ---------------------------------------------------------------------------

func TestEmit_Post_Success(t *testing.T) {
	want := testPayload{ID: "1", Name: "test"}
	var capturedMethod, capturedContentType string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		capturedContentType = r.Header.Get("Content-Type")
		jsonHandler(t, http.StatusCreated, want)(w, r)
	}))
	t.Cleanup(srv.Close)

	a := newAdapter(t, srv.URL)
	req, err := bffstreams.NewJsonRequest(http.MethodPost, "/items", nil, testPayload{Name: "test"})
	require.NoError(t, err)
	resp := &bffstreams.JsonResponseStream[testPayload]{}

	require.NoError(t, a.Emit(context.Background(), req, resp))
	assert.Equal(t, http.MethodPost, capturedMethod)
	assert.Equal(t, "application/json", capturedContentType)
	assert.Equal(t, want, resp.Value)
}

func TestEmit_Post_UpstreamError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "conflict", http.StatusConflict)
	}))
	t.Cleanup(srv.Close)

	a := newAdapter(t, srv.URL)
	req, err := bffstreams.NewJsonRequest(http.MethodPost, "/items", nil, testPayload{})
	require.NoError(t, err)

	resp := bffstreams.NewNoBodyResponseStream()
	require.NoError(t, a.Emit(context.Background(), req, resp))
	assert.Equal(t, http.StatusConflict, resp.StatusCode())
}

// ---------------------------------------------------------------------------
// Emit — GET (no body, query params)
// ---------------------------------------------------------------------------

func TestEmit_Get_Success(t *testing.T) {
	want := testPayload{ID: "42", Name: "hello"}
	var capturedQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		jsonHandler(t, http.StatusOK, want)(w, r)
	}))
	t.Cleanup(srv.Close)

	a := newAdapter(t, srv.URL)
	req := bffstreams.NewJsonNoBodyRequest(http.MethodGet, "/items/42", map[string]string{"fields": "id,name"})
	resp := &bffstreams.JsonResponseStream[testPayload]{}

	require.NoError(t, a.Emit(context.Background(), req, resp))
	assert.Equal(t, want, resp.Value)
	assert.Contains(t, capturedQuery, "fields=id%2Cname")
}

func TestEmit_Get_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	t.Cleanup(srv.Close)

	a := newAdapter(t, srv.URL)
	req := bffstreams.NewJsonNoBodyRequest(http.MethodGet, "/items/99", nil)
	resp := bffstreams.NewNoBodyResponseStream()
	require.NoError(t, a.Emit(context.Background(), req, resp))
	assert.Equal(t, http.StatusNotFound, resp.StatusCode())
}

// ---------------------------------------------------------------------------
// Emit — GET list
// ---------------------------------------------------------------------------

func TestEmit_List_Success(t *testing.T) {
	items := []testPayload{{ID: "1"}, {ID: "2"}}
	srv := httptest.NewServer(jsonHandler(t, http.StatusOK, items))
	t.Cleanup(srv.Close)

	a := newAdapter(t, srv.URL)
	req := bffstreams.NewJsonNoBodyRequest(http.MethodGet, "/items", nil)
	resp := &bffstreams.JsonResponseStream[[]testPayload]{}

	require.NoError(t, a.Emit(context.Background(), req, resp))
	assert.Len(t, resp.Value, 2)
}

// ---------------------------------------------------------------------------
// Emit — PATCH
// ---------------------------------------------------------------------------

func TestEmit_Patch_Success(t *testing.T) {
	updated := testPayload{ID: "1", Name: "updated"}
	var capturedMethod, capturedPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		capturedPath = r.URL.Path
		jsonHandler(t, http.StatusOK, updated)(w, r)
	}))
	t.Cleanup(srv.Close)

	a := newAdapter(t, srv.URL)
	req, err := bffstreams.NewJsonRequest(http.MethodPatch, "/items/1", nil, testPayload{Name: "updated"})
	require.NoError(t, err)
	resp := &bffstreams.JsonResponseStream[testPayload]{}

	require.NoError(t, a.Emit(context.Background(), req, resp))
	assert.Equal(t, http.MethodPatch, capturedMethod)
	assert.Equal(t, "/items/1", capturedPath)
	assert.Equal(t, updated, resp.Value)
}

// ---------------------------------------------------------------------------
// Emit — DELETE
// ---------------------------------------------------------------------------

func TestEmit_Delete_Success(t *testing.T) {
	var capturedURL, capturedMethod string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedURL = r.URL.Path
		capturedMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(srv.Close)

	a := newAdapter(t, srv.URL)
	req := bffstreams.NewJsonNoBodyRequest(http.MethodDelete, "/items/7", nil)

	require.NoError(t, a.Emit(context.Background(), req, bffstreams.NewNoBodyResponseStream()))
	assert.Equal(t, http.MethodDelete, capturedMethod)
	assert.Equal(t, "/items/7", capturedURL)
}

// ---------------------------------------------------------------------------
// BffHealthChecker
// ---------------------------------------------------------------------------

func TestHealthChecker_NoBreakerAlwaysClosed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	t.Cleanup(srv.Close)

	a := newAdapter(t, srv.URL) // no circuit breaker configured
	assert.Equal(t, bffcts.BffCircuitClosed, a.CircuitState())
	assert.True(t, a.IsAvailable())
}

func TestHealthChecker_CircuitOpensAfterMaxFailures(t *testing.T) {
	// Use a server that is closed before requests — all calls get a transport
	// error (connection refused), which is the only kind that trips the breaker.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srvURL := srv.URL
	srv.Close()

	a := newAdapter(t, srvURL, func(cfg *bffcts.BffClientConfig) {
		cfg.CircuitBreaker = bffcts.BffCircuitBreakerConfig{
			Enabled:          true,
			MaxFailures:      3,
			ResetTimeout:     0,
			HalfOpenRequests: 1,
		}
	})

	req := bffstreams.NewJsonNoBodyRequest(http.MethodGet, "/fail", nil)
	for i := 0; i < 3; i++ {
		_ = a.Emit(context.Background(), req, bffstreams.NewNoBodyResponseStream())
	}

	assert.Equal(t, bffcts.BffCircuitOpen, a.CircuitState())
	assert.False(t, a.IsAvailable())
}

// ---------------------------------------------------------------------------
// Custom RoundTripper — verifies transport injection
// ---------------------------------------------------------------------------

type headerCapturingTransport struct {
	base    http.RoundTripper
	headers http.Header
}

func (t *headerCapturingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.headers = req.Header.Clone()
	return t.base.RoundTrip(req)
}

func TestCustomTransport_IsUsed(t *testing.T) {
	srv := httptest.NewServer(jsonHandler(t, http.StatusOK, testPayload{ID: "1"}))
	t.Cleanup(srv.Close)

	cap := &headerCapturingTransport{base: http.DefaultTransport}
	a, err := NewNetHttpAdapter(bffcts.BffClientConfig{BaseURL: srv.URL}, cap)
	require.NoError(t, err)

	req := bffstreams.NewJsonNoBodyRequest(http.MethodGet, "/items/1", nil)
	resp := &bffstreams.JsonResponseStream[testPayload]{}
	_ = a.Emit(context.Background(), req, resp)

	assert.NotNil(t, cap.headers)
}

// ---------------------------------------------------------------------------
// buildURL helper
// ---------------------------------------------------------------------------

func TestBuildURL(t *testing.T) {
	a := &netHttpAdapter{config: bffcts.BffClientConfig{BaseURL: "https://sn.example.com/"}}

	assert.Equal(t, "https://sn.example.com/incidents", a.buildURL("incidents"))
	assert.Equal(t, "https://sn.example.com/incidents/INC001", a.buildURL("incidents/INC001"))
	assert.Equal(t, "https://sn.example.com/incidents/INC001", a.buildURL("/incidents/INC001"))
}
