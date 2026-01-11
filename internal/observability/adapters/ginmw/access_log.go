package ginmw

import (
	"log"
	"time"

	obsports "github.com/brunojet/go-infra-backend/internal/observability/ports"
	"github.com/brunojet/go-infra-backend/internal/observability/requestid"
	"github.com/gin-gonic/gin"
)

// AccessLog writes a minimal structured-ish log line per request.
//
// It is intentionally tiny and dependency-free (uses stdlib log).
func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		lat := time.Since(start)
		rid, _ := requestid.From(c.Request.Context())
		status := c.Writer.Status()
		method := c.Request.Method
		route := obsports.SelectRoute(c.FullPath(), c.Request.URL.Path)

		log.Print(obsports.FormatAccessLog(rid, method, route, status, lat))
	}
}
