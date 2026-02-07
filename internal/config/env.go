package config

import "os"

// Get returns the environment variable value for key if present and non-empty,
// otherwise returns def.
func Get(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}
