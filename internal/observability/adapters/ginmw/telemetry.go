package ginmw

import (
	"fmt"
	"net/http"
	"time"

	noopadapter "github.com/brunojet/go-infra-backend/internal/observability/adapters/nooptelemetry"
	obscontracts "github.com/brunojet/go-infra-backend/internal/observability/contracts"
	obsports "github.com/brunojet/go-infra-backend/internal/observability/ports"
	obstypes "github.com/brunojet/go-infra-backend/internal/observability/types"
	"github.com/gin-gonic/gin"
)

// Telemetry provides baseline tracing/metrics/logs for every HTTP request.
//
// It is middleware-level "wide" observability: measure everything, then later
// add deeper spans/metrics in specific services/repositories when needed.
func Telemetry(p obscontracts.Provider) gin.HandlerFunc {
	if p == nil {
		p = noopadapter.NoopProvider{}
	}
	return func(c *gin.Context) {
		start := time.Now()

		route := obsports.SelectRoute(c.FullPath(), c.Request.URL.Path)

		ctx, span := p.Tracer().Start(
			c.Request.Context(),
			"http "+c.Request.Method+" "+route,
			obstypes.Field{Key: "method", Value: c.Request.Method},
			obstypes.Field{Key: "route", Value: route},
		)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		status := c.Writer.Status()
		lat := time.Since(start)

		p.Metrics().Inc(
			"http.server.requests",
			1,
			obstypes.Field{Key: "method", Value: c.Request.Method},
			obstypes.Field{Key: "route", Value: route},
			obstypes.Field{Key: "status", Value: status},
		)
		p.Metrics().ObserveDuration(
			"http.server.duration",
			lat,
			obstypes.Field{Key: "method", Value: c.Request.Method},
			obstypes.Field{Key: "route", Value: route},
			obstypes.Field{Key: "status", Value: status},
		)

		var err error
		if status >= http.StatusInternalServerError {
			err = fmt.Errorf("http status %d", status)
		}
		span.End(err, obstypes.Field{Key: "status", Value: status}, obstypes.Field{Key: "latency", Value: lat.String()})
	}
}
