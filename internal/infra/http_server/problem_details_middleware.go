package http_server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ProblemDetails represents an HTTP Problem Details response body per RFC 9457.
// Errors is an extension member (RFC 9457 §3.2) listing individual error messages
// when multiple errors are present. Detail holds a stable summary in that case.
type ProblemDetails struct {
	Type     string   `json:"type"`
	Title    string   `json:"title"`
	Status   int      `json:"status"`
	Detail   string   `json:"detail,omitempty"`
	Errors   []string `json:"errors,omitempty"`
	Instance string   `json:"instance,omitempty"`
}

func newProblemDetails(errs []*gin.Error, status int) ProblemDetails {
	var errorMessages []string
	detail := errs[len(errs)-1].Error()
	if len(errs) > 1 {
		errorMessages = make([]string, len(errs))
		for i, e := range errs {
			errorMessages[i] = e.Error()
		}
		detail = "one or more errors occurred"
	}
	return ProblemDetails{
		Type:   "about:blank",
		Title:  http.StatusText(status),
		Status: status,
		Detail: detail,
		Errors: errorMessages,
	}
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

		pd := newProblemDetails(c.Errors, status)
		c.Writer.Header().Set("Content-Type", "application/problem+json")
		c.JSON(status, pd)
	}
}
