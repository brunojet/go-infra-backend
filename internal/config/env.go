package config

import "os"

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
