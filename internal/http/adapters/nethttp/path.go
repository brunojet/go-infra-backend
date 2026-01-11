package nethttp

import "strings"

// translatePath converts a Gin-style path (":id") into a chi-style path ("{id}").
// Example: "/apps/:id" -> "/apps/{id}".
func translatePath(path string) string {
	if path == "" {
		return path
	}
	parts := strings.Split(path, "/")
	for i, p := range parts {
		if len(p) > 1 && p[0] == ':' {
			parts[i] = "{" + p[1:] + "}"
		}
	}
	return strings.Join(parts, "/")
}
