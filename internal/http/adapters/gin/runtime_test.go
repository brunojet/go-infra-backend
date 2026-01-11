package ginadapter

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDefaultRuntimeOptions(t *testing.T) {
	o := defaultRuntimeOptions()
	assert.Nil(t, o.middlewares)
}

func TestWithMiddlewares(t *testing.T) {
	o := defaultRuntimeOptions()
	mw := func(c *gin.Context) {}
	WithMiddlewares(mw)(&o)
	// Use the same variable for both adding and checking
	assert.True(t, len(o.middlewares) == 1 && getFuncPtr(o.middlewares[0]) == getFuncPtr(mw), "middleware should be present")
}

// getFuncPtr returns the pointer value of a function for comparison
func getFuncPtr(f interface{}) uintptr {
	return reflect.ValueOf(f).Pointer()
}

func TestNewRuntime(t *testing.T) {
	router, handler := NewRuntime()
	assert.NotNil(t, router)
	assert.NotNil(t, handler)

	// Test with middleware
	mw := func(c *gin.Context) { c.Set("mw", true); c.Next() }
	router2, handler2 := NewRuntime(WithMiddlewares(mw))
	assert.NotNil(t, router2)
	assert.NotNil(t, handler2)

	// Test handler is http.Handler
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	handler2.ServeHTTP(w, req)
}
