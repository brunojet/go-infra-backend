package httpmw

import (
	"fmt"
	"net/http"
	"time"

	nooptelemetry "github.com/brunojet/go-infra-backend/internal/observability/adapters/nooptelemetry"
	"github.com/brunojet/go-infra-backend/internal/observability/contracts"
	obsports "github.com/brunojet/go-infra-backend/internal/observability/ports"
	obstypes "github.com/brunojet/go-infra-backend/internal/observability/types"
	"github.com/go-chi/chi/v5"
)

// Telemetry provides baseline tracing/metrics/logs for every HTTP request.
func Telemetry(p contracts.Provider) func(http.Handler) http.Handler {
	if p == nil {
		p = nooptelemetry.NoopProvider{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sr := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			start := time.Now()

			route := obsports.SelectRoute(chi.RouteContext(r.Context()).RoutePattern(), r.URL.Path)

			ctx, span := p.Tracer().Start(
				r.Context(),
				"http "+r.Method+" "+route,
				obstypes.Field{Key: "method", Value: r.Method},
				obstypes.Field{Key: "route", Value: route},
			)

			next.ServeHTTP(sr, r.WithContext(ctx))

			lat := time.Since(start)
			status := sr.status

			p.Metrics().Inc(
				"http.server.requests",
				1,
				obstypes.Field{Key: "method", Value: r.Method},
				obstypes.Field{Key: "route", Value: route},
				obstypes.Field{Key: "status", Value: status},
			)
			p.Metrics().ObserveDuration(
				"http.server.duration",
				lat,
				obstypes.Field{Key: "method", Value: r.Method},
				obstypes.Field{Key: "route", Value: route},
				obstypes.Field{Key: "status", Value: status},
			)

			var err error
			if status >= http.StatusInternalServerError {
				err = fmt.Errorf("http status %d", status)
			}
			span.End(err, obstypes.Field{Key: "status", Value: status}, obstypes.Field{Key: "latency", Value: lat.String()})
		})
	}
}
