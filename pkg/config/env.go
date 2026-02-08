package config

import internalconfig "github.com/brunojet/go-infra-backend/internal/config"

// GetEnv delegates to internal/config.GetEnv.
// Public surface is intentionally small.
func GetEnv(key, def string) string {
	return internalconfig.GetEnv(key, def)
}

// GetEnvAsBool delegates to internal/config.GetEnvAsBool.
func GetEnvAsBool(key string, def bool) bool {
	return internalconfig.GetEnvAsBool(key, def)
}

// IsValidEndpoint reports whether endpoint matches host:port or :port.
// Host can be a hostname, IPv4, or IPv6 in brackets.
func IsValidEndpoint(endpoint string) bool {
	return internalconfig.IsValidEndpoint(endpoint)
}

// ValidateEndpoint validates that endpoint matches host:port or :port.
func ValidateEndpoint(endpoint string) error {
	return internalconfig.ValidateEndpoint(endpoint)
}
