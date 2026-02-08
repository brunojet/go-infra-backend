package config

import (
	"fmt"
	"os"
	"regexp"
)

// endpointRegex matches:
// - ":5432" (empty host)
// - "localhost:5432" (hostname)
// - "127.0.0.1:5432" (IPv4)
// - "[::1]:5432" (IPv6 in brackets)
// Port range is validated to 0..65535.
//
// Notes:
// - This is intentionally strict (no schemes, no paths).
// - For full URL endpoints, prefer net/url and validate url.URL.
var endpointRegex = regexp.MustCompile(
	`^(?:(?:` +
		// hostname (RFC 1123-ish): labels separated by dots
		`[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?(?:\.[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?)*` +
		`)|(?:` +
		// IPv4
		`(?:25[0-5]|2[0-4]\d|1?\d?\d)(?:\.(?:25[0-5]|2[0-4]\d|1?\d?\d)){3}` +
		`)|(?:` +
		// IPv6 (basic) in brackets
		`\[[0-9A-Fa-f:]+\]` +
		`))?:(0|6553[0-5]|655[0-2]\d|65[0-4]\d{2}|6[0-4]\d{3}|[1-5]\d{4}|[1-9]\d{0,3})$`,
)

// IsValidEndpoint returns true if endpoint is in the form host:port or :port.
// Host can be a hostname, IPv4, or IPv6 in brackets.
func IsValidEndpoint(endpoint string) bool {
	return endpointRegex.MatchString(endpoint)
}

// ValidateEndpoint validates the endpoint format.
// It returns an error describing the expected format when invalid.
func ValidateEndpoint(endpoint string) error {
	if endpoint == "" {
		return fmt.Errorf("endpoint is empty")
	}
	if !IsValidEndpoint(endpoint) {
		return fmt.Errorf("invalid endpoint %q: expected host:port or :port (host can be hostname, IPv4, or [IPv6])", endpoint)
	}
	return nil
}

// Get returns the environment variable value for key if present and non-empty,
// otherwise returns def.
func GetEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func GetEnvAsBool(key string, def bool) bool {
	{
		if v, ok := os.LookupEnv(key); ok && v != "" {
			switch v {
			case "1", "true", "TRUE", "True":
				return true
			case "0", "false", "FALSE", "False":
				return false
			}
		}
		return def
	}
}
