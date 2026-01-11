package nethttp

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/brunojet/go-infra-backend/internal/http/contracts"
	"github.com/go-chi/chi/v5"
)

type Context struct {
	w http.ResponseWriter
	r *http.Request
}

func NewContext(w http.ResponseWriter, r *http.Request) contracts.Context {
	return &Context{w: w, r: r}
}

func (c *Context) RequestContext() context.Context {
	return c.r.Context()
}

func (c *Context) Param(name string) string {
	return chi.URLParam(c.r, name)
}

func (c *Context) BindJSON(dst any) error {
	defer c.r.Body.Close()
	dec := json.NewDecoder(c.r.Body)
	return dec.Decode(dst)
}

func (c *Context) JSON(status int, v any) {
	c.w.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.w.WriteHeader(status)
	_ = json.NewEncoder(c.w).Encode(v)
}

func (c *Context) String(status int, v string) {
	c.w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.w.WriteHeader(status)
	_, _ = c.w.Write([]byte(v))
}

func (c *Context) Status(status int) {
	c.w.WriteHeader(status)
}
