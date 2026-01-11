package server

import (
	"time"
)

// ServerParams holds HTTP server configuration.
//
// Env:
//   - PORT (default: 8080)
//   - READ_HEADER_TIMEOUT (default: 5s)
type ServerParams struct {
	Addr              string
	ReadHeaderTimeout time.Duration
}
