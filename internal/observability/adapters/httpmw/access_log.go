package httpmw

import (
	"log"
	"net/http"
	"time"

	obsports "github.com/brunojet/go-infra-backend/internal/observability/ports"
	"github.com/brunojet/go-infra-backend/internal/observability/requestid"
	"github.com/go-chi/chi/v5"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

// AccessLog writes a minimal log line per request.
func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sr := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		start := time.Now()
		next.ServeHTTP(sr, r)

		lat := time.Since(start)
		rid, _ := requestid.From(r.Context())
		route := obsports.SelectRoute(chi.RouteContext(r.Context()).RoutePattern(), r.URL.Path)

		log.Print(obsports.FormatAccessLog(rid, r.Method, route, sr.status, lat))
	})
}
