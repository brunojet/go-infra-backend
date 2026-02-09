package config

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
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

var hostnameRegex = regexp.MustCompile(
	`^[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?(?:\.[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?)*$`,
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

// IsValidHost returns true if host is a valid hostname, IPv4, or IPv6.
// It does not accept empty string.
func IsValidHost(host string) bool {
	if host == "" {
		return false
	}
	if len(host) >= 2 && host[0] == '[' && host[len(host)-1] == ']' {
		host = host[1 : len(host)-1]
	}
	if ip := net.ParseIP(host); ip != nil {
		return true
	}
	return hostnameRegex.MatchString(host)
}

// ValidateHost validates a host (hostname, IPv4, or IPv6).
func ValidateHost(host string) error {
	if host == "" {
		return fmt.Errorf("host is empty")
	}
	if !IsValidHost(host) {
		return fmt.Errorf("invalid host %q: expected hostname, IPv4, or IPv6", host)
	}
	return nil
}

// IsValidPort returns true if port is numeric and within 0..65535.
// It does not accept empty string.
func IsValidPort(port string) bool {
	if port == "" {
		return false
	}
	n, err := strconv.Atoi(port)
	if err != nil {
		return false
	}
	return n >= 0 && n <= 65535
}

// ValidatePort validates a TCP/UDP port in 0..65535.
func ValidatePort(port string) error {
	if port == "" {
		return fmt.Errorf("port is empty")
	}
	if !IsValidPort(port) {
		return fmt.Errorf("invalid port %q: expected 0..65535", port)
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

func GetEnvAsInt(key string, def int) int {
	v := GetEnv(key, strconv.Itoa(def))
	n, err := strconv.Atoi(v)
	if err == nil {
		return n
	}
	return def
}

func GetEnvAsBool(key string, def bool) bool {
	v := GetEnv(key, strconv.FormatBool(def))
	if v != "" {
		switch v {
		case "1", "true", "TRUE", "True":
			return true
		case "0", "false", "FALSE", "False":
			return false
		}
	}
	return def
}
