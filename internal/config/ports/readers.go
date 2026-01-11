package ports

import (
	"strconv"
	"strings"
	"time"

	"github.com/brunojet/go-infra-backend/internal/config/contracts"
)

// Trimmed returns the trimmed value for key, or empty string if not set.
func Trimmed(src contracts.Source, key string) string {
	v, _ := src.Lookup(key)
	return strings.TrimSpace(v)
}

// LowerTrimmed returns the trimmed, lowercased value for key, or empty string if not set.
func LowerTrimmed(src contracts.Source, key string) string {
	return strings.ToLower(Trimmed(src, key))
}

// Bool returns a boolean parsed from key with a default.
func Bool(src contracts.Source, key string, def bool) bool {
	v, ok := src.Lookup(key)
	if !ok {
		return def
	}
	vv := strings.TrimSpace(strings.ToLower(v))
	if vv == "" {
		return def
	}
	switch vv {
	case "1", "true", "t", "yes", "y", "on", "enabled", "enable":
		return true
	case "0", "false", "f", "no", "n", "off", "disabled", "disable":
		return false
	default:
		return def
	}
}

// Duration returns a duration parsed from key with a default.
//
// Supported formats:
//   - Go duration (e.g. "250ms", "5s", "1m")
//   - seconds as int (e.g. "5")
func Duration(src contracts.Source, key string, def time.Duration) time.Duration {
	v, ok := src.Lookup(key)
	if !ok {
		return def
	}
	vv := strings.TrimSpace(v)
	if vv == "" {
		return def
	}
	if d, err := time.ParseDuration(vv); err == nil {
		return d
	}
	if n, err := strconv.Atoi(vv); err == nil {
		return time.Duration(n) * time.Second
	}
	return def
}

// Int returns an int parsed from key with a default.
//
// If the key is not set, empty, or invalid, def is returned.
func Int(src contracts.Source, key string, def int) int {
	v, ok := src.Lookup(key)
	if !ok {
		return def
	}
	vv := strings.TrimSpace(v)
	if vv == "" {
		return def
	}
	n, err := strconv.Atoi(vv)
	if err != nil {
		return def
	}
	return n
}

// ParseKeyValueCSV parses a comma-separated list of key=value pairs.
//
// Examples:
//   - "k=v,k2=v2" => map[k]v map[k2]v2
//
// Invalid parts are ignored. Empty input returns nil.
func ParseKeyValueCSV(raw string) map[string]string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}

	out := map[string]string{}
	for _, part := range strings.Split(raw, ",") {
		p := strings.TrimSpace(part)
		if p == "" {
			continue
		}
		k, v, ok := strings.Cut(p, "=")
		if !ok {
			continue
		}
		kk := strings.TrimSpace(k)
		if kk == "" {
			continue
		}
		out[kk] = strings.TrimSpace(v)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// SplitCSV splits a comma-separated list into a slice.
//
// It trims whitespace, ignores empty items, and returns def when raw is empty
// (or when all items are empty after trimming).
func SplitCSV(raw string, def []string) []string {
	if strings.TrimSpace(raw) == "" {
		return def
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v == "" {
			continue
		}
		out = append(out, v)
	}
	if len(out) == 0 {
		return def
	}
	return out
}

// SplitCSVUpper is like SplitCSV but uppercases each returned item.
func SplitCSVUpper(raw string, def []string) []string {
	items := SplitCSV(raw, def)
	for i := range items {
		items[i] = strings.ToUpper(items[i])
	}
	return items
}
