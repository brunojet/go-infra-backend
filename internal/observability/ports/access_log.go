package ports

import (
	"fmt"
	"time"
)

// FormatAccessLog returns the canonical access log line used by infra.
func FormatAccessLog(rid string, method string, route string, status int, latency time.Duration) string {
	return fmt.Sprintf(
		"http request rid=%s method=%s route=%s status=%d latency=%s",
		rid,
		method,
		route,
		status,
		latency,
	)
}
