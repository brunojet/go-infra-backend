package httpserver

import (
	internalhttpserver "github.com/brunojet/go-infra-backend/internal/infra/http_server"
	"github.com/gin-gonic/gin"
)

// ProblemDetails is the RFC 9457 response body type.
type ProblemDetails = internalhttpserver.ProblemDetails

// CORSMiddleware returns a Gin middleware that sets CORS headers.
// Allowed origins are read from the CORS_ALLOWED_ORIGINS environment variable
// (comma-separated). Falls back to echoing the request Origin (permissive).
func CORSMiddleware() gin.HandlerFunc {
	return internalhttpserver.CORSMiddleware()
}

// ProblemDetailsMiddleware returns a Gin middleware that formats errors registered
// via c.Error(err) as RFC 9457 Problem Details JSON (application/problem+json).
func ProblemDetailsMiddleware() gin.HandlerFunc {
	return internalhttpserver.ProblemDetailsMiddleware()
}
