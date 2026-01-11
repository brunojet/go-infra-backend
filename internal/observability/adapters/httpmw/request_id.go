package httpmw

import (
	"net/http"
	"strings"

	"github.com/brunojet/go-infra-backend/internal/observability/requestid"
)

// RequestID ensures every request has a correlation id and exposes it in:
// - response header X-Request-ID
// - request context (for downstream handlers)
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.Header.Get(requestid.Header))
		if id == "" {
			id = requestid.New()
		}
		w.Header().Set(requestid.Header, id)
		ctx := requestid.With(r.Context(), id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
