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
