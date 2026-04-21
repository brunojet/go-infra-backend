package repositories

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/brunojet/go-infra-backend/pkg/ports/bff/repositories/contracts"
)

type stubHTTPAdapter struct {
	do func(ctx context.Context, req *http.Request) (*http.Response, error)
}

func (s *stubHTTPAdapter) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	return s.do(ctx, req)
}

type sampleEntity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func TestRestRepository_Create_UsesPOSTAndDecodesBody(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody string

	adapter := &stubHTTPAdapter{do: func(_ context.Context, req *http.Request) (*http.Response, error) {
		gotMethod = req.Method
		gotPath = req.URL.Path
		if req.Body != nil {
			b, _ := io.ReadAll(req.Body)
			gotBody = string(b)
		}
		payload, _ := json.Marshal(sampleEntity{ID: "1", Name: "created"})
		return &http.Response{
			StatusCode: http.StatusCreated,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(string(payload))),
		}, nil
	}}

	repo := NewRestRepository[sampleEntity, map[string]any](adapter,
		WithBaseURL("https://example.com/api"),
		WithPathConfig("items"),
	)

	resp := &contracts.RestResponse[sampleEntity]{}
	err := repo.Create(context.Background(), contracts.RestRequest{Body: map[string]any{"name": "new-item"}}, resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Fatalf("expected method POST, got %s", gotMethod)
	}
	if gotPath != "/api/items" {
		t.Fatalf("expected path /api/items, got %s", gotPath)
	}
	if !strings.Contains(gotBody, "new-item") {
		t.Fatalf("expected request body to contain payload, got %q", gotBody)
	}
	if resp.Body.ID != "1" || resp.Body.Name != "created" {
		t.Fatalf("unexpected response body: %+v", resp.Body)
	}
}

func TestRestRepository_List_UsesGETAndDecodesBodies(t *testing.T) {
	var gotMethod, gotPath, gotQuery string
	adapter := &stubHTTPAdapter{do: func(_ context.Context, req *http.Request) (*http.Response, error) {
		gotMethod = req.Method
		gotPath = req.URL.Path
		gotQuery = req.URL.RawQuery
		payload, _ := json.Marshal([]sampleEntity{{ID: "1", Name: "one"}, {ID: "2", Name: "two"}})
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(string(payload))),
		}, nil
	}}

	repo := NewRestRepository[sampleEntity, map[string]any](adapter,
		WithBaseURL("https://example.com/api"),
		WithPathConfig("items"),
	)

	opts := &contracts.RestRequestOptions{Values: map[string][]string{"page": {"1"}}}
	resp := &contracts.RestResponse[sampleEntity]{}
	var meta map[string]any
	err := repo.List(context.Background(), opts, resp, &meta)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotMethod != http.MethodGet {
		t.Fatalf("expected method GET, got %s", gotMethod)
	}
	if gotPath != "/api/items" {
		t.Fatalf("expected path /api/items, got %s", gotPath)
	}
	if gotQuery != "page=1" {
		t.Fatalf("expected query page=1, got %s", gotQuery)
	}
	if len(resp.Bodies) != 2 {
		t.Fatalf("expected 2 bodies, got %d", len(resp.Bodies))
	}
}

func TestRestRepository_Delete_204ReturnsNoContentMessage(t *testing.T) {
	var gotMethod, gotPath string
	adapter := &stubHTTPAdapter{do: func(_ context.Context, req *http.Request) (*http.Response, error) {
		gotMethod = req.Method
		gotPath = req.URL.Path
		return &http.Response{
			StatusCode: http.StatusNoContent,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	}}

	repo := NewRestRepository[sampleEntity, map[string]any](adapter,
		WithBaseURL("https://example.com/api"),
		WithPathConfig("items"),
	)

	resp := &contracts.RestResponse[sampleEntity]{}
	err := repo.Delete(context.Background(), nil, resp, "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotMethod != http.MethodDelete {
		t.Fatalf("expected method DELETE, got %s", gotMethod)
	}
	if gotPath != "/api/items/123" {
		t.Fatalf("expected path /api/items/123, got %s", gotPath)
	}
	if resp.Opts.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", resp.Opts.StatusCode)
	}
}
