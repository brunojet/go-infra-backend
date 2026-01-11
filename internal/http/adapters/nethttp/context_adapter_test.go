package nethttp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewContextAndRequestContext(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	ctx := NewContext(w, r)
	if ctx.RequestContext() != r.Context() {
		t.Error("RequestContext should return the request's context")
	}
}

func TestContext_Param(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/users/123", nil)
	r = r.WithContext(context.WithValue(r.Context(), "chiRouteContext", map[string]string{"id": "123"}))
	w := httptest.NewRecorder()
	ctx := &Context{w: w, r: r}
	// chi.URLParam will return "" since chiRouteContext is not a real chi context, but we check no panic
	_ = ctx.Param("id")
}

func TestContext_BindJSON(t *testing.T) {
	type payload struct{ Foo string }
	body := bytes.NewBufferString(`{"Foo":"bar"}`)
	r := httptest.NewRequest(http.MethodPost, "/", body)
	w := httptest.NewRecorder()
	ctx := &Context{w: w, r: r}
	var p payload
	err := ctx.BindJSON(&p)
	if err != nil || p.Foo != "bar" {
		t.Errorf("BindJSON failed: %v, %+v", err, p)
	}
}

func TestContext_JSON(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := &Context{w: w, r: r}
	ctx.JSON(201, map[string]string{"foo": "bar"})
	resp := w.Result()
	if resp.StatusCode != 201 {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Errorf("expected content-type json, got %s", ct)
	}
	var m map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&m)
	if m["foo"] != "bar" {
		t.Errorf("expected foo=bar, got %+v", m)
	}
}

func TestContext_String(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := &Context{w: w, r: r}
	ctx.String(202, "hello")
	resp := w.Result()
	if resp.StatusCode != 202 {
		t.Errorf("expected status 202, got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "text/plain; charset=utf-8" {
		t.Errorf("expected content-type text/plain, got %s", ct)
	}
	b, _ := io.ReadAll(resp.Body)
	if string(b) != "hello" {
		t.Errorf("expected body 'hello', got %q", string(b))
	}
}

func TestContext_Status(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := &Context{w: w, r: r}
	ctx.Status(204)
	resp := w.Result()
	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}
