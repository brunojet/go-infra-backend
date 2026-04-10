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

// IsValidHost reports whether host is a valid hostname, IPv4, or IPv6.
func IsValidHost(host string) bool {
	return internalconfig.IsValidHost(host)
}

// ValidateHost validates host (hostname, IPv4, or IPv6).
func ValidateHost(host string) error {
	return internalconfig.ValidateHost(host)
}

// IsValidPort reports whether port is numeric and within 0..65535.
func IsValidPort(port string) bool {
	return internalconfig.IsValidPort(port)
}

// ValidatePort validates port (0..65535).
func ValidatePort(port string) error {
	return internalconfig.ValidatePort(port)
}
