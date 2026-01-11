package nethttp

import (
	"net/http"

	httpports "github.com/brunojet/go-infra-backend/internal/http/ports"
	httptypes "github.com/brunojet/go-infra-backend/internal/http/types"
)

// CORS returns a net/http middleware that applies CORS headers and handles preflight OPTIONS.
func CORS(cfg httptypes.CORSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if h, ok := httpports.BuildCORSHeaders(cfg, origin); ok {
				w.Header().Set("Access-Control-Allow-Origin", h.AllowOrigin)
				if h.VaryOrigin {
					w.Header().Set("Vary", "Origin")
				}
				w.Header().Set("Access-Control-Allow-Methods", h.AllowMethods)
				w.Header().Set("Access-Control-Allow-Headers", h.AllowHeaders)
				if h.ExposeHeaders != "" {
					w.Header().Set("Access-Control-Expose-Headers", h.ExposeHeaders)
				}
				if h.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
				if h.MaxAgeSeconds > 0 {
					w.Header().Set("Access-Control-Max-Age", h.MaxAgeString())
				}
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
