package types

import "time"

// CORSConfig configures Cross-Origin Resource Sharing (CORS) behavior.
//
// This is a contract between the infra/http wiring layer and the HTTP framework
// adapters (gin, net/http).
//
// Guardrails (production-safe defaults):
//   - Prefer leaving CORS disabled unless explicitly enabled by configuration.
//   - Never combine AllowCredentials=true with AllowOrigins containing "*".
//     If that invalid combination is configured, the adapters will not apply CORS
//     headers for any request origin.
type CORSConfig struct {
	Enabled          bool
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           time.Duration
}
