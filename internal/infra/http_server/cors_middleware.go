package http_server

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func allowedOriginsFromEnv() []string {
	raw := os.Getenv("CORS_ALLOWED_ORIGINS")
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		origin := strings.TrimSpace(part)
		if origin != "" {
			origins = append(origins, origin)
		}
	}

	return origins
}

func resolveAllowedOrigin(requestOrigin string, allowedOrigins []string) (string, bool) {
	if len(allowedOrigins) == 0 {
		return "*", false
	}

	if requestOrigin == "" {
		return "", false
	}

	for _, allowedOrigin := range allowedOrigins {
		if allowedOrigin == requestOrigin {
			return requestOrigin, true
		}
	}

	return "", false
}

// CORSMiddleware adds CORS headers and handles OPTIONS requests globally.
// Allowed origins are read from the CORS_ALLOWED_ORIGINS env var (comma-separated).
// When the env var is absent, all origins are allowed ("*").
func CORSMiddleware() gin.HandlerFunc {
	allowedOrigins := allowedOriginsFromEnv()

	return func(c *gin.Context) {
		origin, allowCredentials := resolveAllowedOrigin(c.GetHeader("Origin"), allowedOrigins)
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			if allowCredentials {
				c.Writer.Header().Set("Vary", "Origin")
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			}
		}
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
