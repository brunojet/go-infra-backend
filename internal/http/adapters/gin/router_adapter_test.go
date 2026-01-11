package ginadapter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brunojet/go-infra-backend/internal/http/contracts"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGinRouter_BasicRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	e := gin.New()
	r := NewGinRouter(e.Group(""))

	var called string
	r.GET("/get", contracts.HandlerFunc(func(ctx contracts.Context) {
		called = "get"
	}))
	r.POST("/post", contracts.HandlerFunc(func(ctx contracts.Context) {
		called = "post"
	}))
	r.PATCH("/patch", contracts.HandlerFunc(func(ctx contracts.Context) {
		called = "patch"
	}))
	r.DELETE("/delete", contracts.HandlerFunc(func(ctx contracts.Context) {
		called = "delete"
	}))

	// Versão correta, igual nethttp:
	// h := contracts.HandlerFunc(func(ctx contracts.Context) { called = "get" })
	// r.GET("/get", h)
	// ...

	// Cast handlers to contracts.HandlerFunc (func(contracts.Context))
	// Para evitar erro de tipo, use:
	// r.GET("/get", contracts.HandlerFunc(func(ctx contracts.Context) { ... }))

	for _, route := range []struct {
		method, path, expect string
	}{
		{"GET", "/get", "get"},
		{"POST", "/post", "post"},
		{"PATCH", "/patch", "patch"},
		{"DELETE", "/delete", "delete"},
	} {
		called = ""
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(route.method, route.path, nil)
		e.ServeHTTP(w, req)
		assert.Equal(t, route.expect, called)
	}
}

func TestGinRouter_Group(t *testing.T) {
	gin.SetMode(gin.TestMode)
	e := gin.New()
	r := NewGinRouter(e.Group(""))
	sub := r.Group("/sub")
	assert.NotNil(t, sub)
}
