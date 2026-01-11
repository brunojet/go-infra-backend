package requestid

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

// Header is the canonical HTTP header name used for request correlation.
const Header = "X-Request-ID"

type ctxKey struct{}

// New generates a reasonably unique, non-guessable request id.
func New() string {
	var b [16]byte
	_, err := rand.Read(b[:])
	if err != nil {
		// Extremely unlikely; fall back to something deterministic-ish.
		return hex.EncodeToString([]byte("fallback-request-id"))
	}
	return hex.EncodeToString(b[:])
}

// With stores request id into context.
func With(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKey{}, id)
}

// From extracts request id from context.
func From(ctx context.Context) (string, bool) {
	v := ctx.Value(ctxKey{})
	id, ok := v.(string)
	return id, ok && id != ""
}
