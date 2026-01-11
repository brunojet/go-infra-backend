package types

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestFieldStruct(t *testing.T) {
	f := Field{Key: "foo", Value: 123}
	if f.Key != "foo" || f.Value != 123 {
		t.Errorf("Field struct not working as expected: %+v", f)
	}
}

func TestMiddlewareConfigStruct(t *testing.T) {
	cfg := MiddlewareConfig{RequestID: true, AccessLog: true, Telemetry: false, Recovery: true}
	if !cfg.RequestID || !cfg.AccessLog || cfg.Telemetry || !cfg.Recovery {
		t.Errorf("MiddlewareConfig struct not working as expected: %+v", cfg)
	}
}

func TestMiddlewaresStruct(t *testing.T) {
	ginMw := func(*gin.Context) {}
	nethttpMw := func(next http.Handler) http.Handler { return next }
	m := Middlewares{
		Gin:     []GinMiddleware{ginMw},
		NetHTTP: []NetHTTPMiddleware{nethttpMw},
	}
	if len(m.Gin) != 1 || len(m.NetHTTP) != 1 {
		t.Errorf("Middlewares struct not working as expected: %+v", m)
	}
}
