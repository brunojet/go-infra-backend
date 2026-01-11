package ports

import (
	"fmt"
	"slices"
	"strings"

	httptypes "github.com/brunojet/go-infra-backend/internal/http/types"
)

// CORSHeaders represents the headers that should be applied when handling a CORS request.
//
// Adapters are responsible for applying these headers to their framework-specific
// response types.
type CORSHeaders struct {
	AllowOrigin      string
	VaryOrigin       bool
	AllowMethods     string
	AllowHeaders     string
	ExposeHeaders    string
	AllowCredentials bool
	MaxAgeSeconds    int
}

// BuildCORSHeaders computes the headers for a given Origin.
//
// ok is false when the Origin should not be allowed (or Origin is empty).
func BuildCORSHeaders(cfg httptypes.CORSConfig, origin string) (h CORSHeaders, ok bool) {
	if strings.TrimSpace(origin) == "" {
		return CORSHeaders{}, false
	}
	// Guardrail: credentials must never be enabled with wildcard origins.
	// In this invalid configuration, do not apply any CORS headers.
	if cfg.AllowCredentials && slices.Contains(cfg.AllowOrigins, "*") {
		return CORSHeaders{}, false
	}
	if !slices.Contains(cfg.AllowOrigins, "*") && !slices.Contains(cfg.AllowOrigins, origin) {
		return CORSHeaders{}, false
	}

	allowAll := len(cfg.AllowOrigins) == 1 && cfg.AllowOrigins[0] == "*" && !cfg.AllowCredentials
	if allowAll {
		h.AllowOrigin = "*"
		h.VaryOrigin = false
	} else {
		h.AllowOrigin = origin
		h.VaryOrigin = true
	}

	h.AllowMethods = strings.Join(cfg.AllowMethods, ", ")
	h.AllowHeaders = strings.Join(cfg.AllowHeaders, ", ")
	h.ExposeHeaders = strings.Join(cfg.ExposeHeaders, ", ")
	h.AllowCredentials = cfg.AllowCredentials

	maxAge := int(cfg.MaxAge.Seconds())
	if maxAge < 0 {
		maxAge = 0
	}
	h.MaxAgeSeconds = maxAge

	return h, true
}

func (h CORSHeaders) MaxAgeString() string {
	return fmt.Sprintf("%d", h.MaxAgeSeconds)
}
