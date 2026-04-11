package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ProblemDetails represents an HTTP Problem Details response body per RFC 9457.
type ProblemDetails struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

// ProblemDetailsMiddleware formats errors registered via c.Error(err) as RFC 9457
// Problem Details JSON responses.
//
// Handlers should call c.Error(err) to register the error and c.Status(code) to
// set the HTTP status code, then return without writing a response body. This
// middleware writes the JSON body after c.Next() completes.
//
// If no errors were registered, the middleware is a no-op.
//
// Content-Type is set to "application/problem+json" as required by RFC 9457.
func ProblemDetailsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		status := c.Writer.Status()
		if status < http.StatusBadRequest {
			status = http.StatusInternalServerError
		}

		c.JSON(status, ProblemDetails{
			Type:     "about:blank",
			Title:    http.StatusText(status),
			Status:   status,
			Detail:   c.Errors.Last().Err.Error(),
			Instance: c.Request.URL.Path,
		})
	}
}
