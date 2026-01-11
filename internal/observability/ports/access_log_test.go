package ports

import (
	"strings"
	"testing"
	"time"
)

func TestFormatAccessLog_UsesRouteKey(t *testing.T) {
	line := FormatAccessLog("rid-123", "GET", "/apps/{id}", 200, 12*time.Millisecond)

	if !strings.Contains(line, "rid=rid-123") {
		t.Fatalf("log line missing rid: %q", line)
	}
	if !strings.Contains(line, "method=GET") {
		t.Fatalf("log line missing method: %q", line)
	}
	if !strings.Contains(line, "route=/apps/{id}") {
		t.Fatalf("log line missing route=: %q", line)
	}
	if strings.Contains(line, " path=") {
		t.Fatalf("log line should not use path=: %q", line)
	}
}
