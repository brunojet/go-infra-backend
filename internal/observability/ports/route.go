package ports

import "strings"

// RouteUnmatched is used when the router template/pattern is unavailable.
//
// This prevents high-cardinality observability labels caused by falling back to
// the raw URL path (which may contain IDs, hashes, etc.).
const RouteUnmatched = "unmatched"

// SelectRoute prefers a low-cardinality route template when available.
//
// For example:
//   - template: "/apps/:id" (Gin) or "/apps/{id}" (Chi)
//   - actual:   "/apps/123"
//
// When template is empty, it returns RouteUnmatched.
func SelectRoute(template string, _ string) string {
	if strings.TrimSpace(template) != "" {
		return template
	}
	return RouteUnmatched
}
