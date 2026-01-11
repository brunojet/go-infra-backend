package ginadapter

import (
	"context"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type dummyWriter struct{ gin.ResponseWriter }

func (d *dummyWriter) Header() http.Header         { return make(http.Header) }
func (d *dummyWriter) Write(b []byte) (int, error) { return len(b), nil }
func (d *dummyWriter) WriteHeader(statusCode int)  {}

func TestNewGinContext(t *testing.T) {
	c, _ := gin.CreateTestContext(&dummyWriter{})
	ctx := NewGinContext(c)
	assert.NotNil(t, ctx)
}

func TestGinContext_RequestContext(t *testing.T) {
	c, _ := gin.CreateTestContext(&dummyWriter{})
	// Inicializa c.Request se necessário
	if c.Request == nil {
		c.Request = &http.Request{Header: make(http.Header)}
	}
	c.Request = c.Request.WithContext(context.WithValue(context.Background(), "k", "v"))
	ctx := NewGinContext(c)
	assert.Equal(t, "v", ctx.RequestContext().Value("k"))
}

func TestGinContext_Param(t *testing.T) {
	c, _ := gin.CreateTestContext(&dummyWriter{})
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "42"})
	ctx := NewGinContext(c)
	assert.Equal(t, "42", ctx.Param("id"))
}

func TestGinContext_BindJSON(t *testing.T) {
	c, _ := gin.CreateTestContext(&dummyWriter{})
	ctx := NewGinContext(c)
	// ShouldBindJSON returns error if no body, just check error is not nil
	var dst struct{}
	err := ctx.BindJSON(&dst)
	assert.Error(t, err)
}

func TestGinContext_JSON_String_Status(t *testing.T) {
	c, _ := gin.CreateTestContext(&dummyWriter{})
	ctx := NewGinContext(c)
	// JSON
	ctx.JSON(200, gin.H{"ok": true})
	// String
	ctx.String(201, "hello")
	// Status
	ctx.Status(204)
	// No assertion: just ensure no panic
}
