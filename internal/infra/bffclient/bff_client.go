package bffclient

import (
	"context"
	"sync"
)

// ---------------------------------------------------------------------------
// Header propagation context helpers
// ---------------------------------------------------------------------------

type bffReqHeadersKey struct{}
type bffRespHeadersKey struct{}

// responseHeadersBag is a goroutine-safe mutable store for captured response headers.
// A pointer to it is placed in ctx by InitResponseCapture so that both the BFF adapter
// (writer) and BffHeadersMiddleware (reader) share the same instance.
type responseHeadersBag struct {
	mu sync.Mutex
	m  map[string]string
}

func (b *responseHeadersBag) set(k, v string) {
	b.mu.Lock()
	b.m[k] = v
	b.mu.Unlock()
}

func (b *responseHeadersBag) all() map[string]string {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make(map[string]string, len(b.m))
	for k, v := range b.m {
		out[k] = v
	}
	return out
}

// WithRequestHeaders stores header key→value pairs in ctx so that downstream
// BFF adapters can forward them to the outgoing upstream request.
// Typically called by BffHeadersMiddleware which reads the incoming HTTP request headers.
func WithRequestHeaders(ctx context.Context, headers map[string]string) context.Context {
	return context.WithValue(ctx, bffReqHeadersKey{}, headers)
}

// RequestHeadersFromCtx returns the headers stored by WithRequestHeaders, or nil.
func RequestHeadersFromCtx(ctx context.Context) map[string]string {
	m, _ := ctx.Value(bffReqHeadersKey{}).(map[string]string)
	return m
}

// InitResponseCapture installs a mutable response-header bag in ctx.
// Must be called by BffHeadersMiddleware before c.Next() so that BFF adapters
// can write captured upstream headers via CaptureResponseHeader.
func InitResponseCapture(ctx context.Context) context.Context {
	return context.WithValue(ctx, bffRespHeadersKey{}, &responseHeadersBag{m: make(map[string]string)})
}

// CaptureResponseHeader writes key→value into the response bag stored in ctx.
// Called by the BFF adapter after receiving an upstream response.
// No-op when the ctx has no bag (i.e. InitResponseCapture was not called).
func CaptureResponseHeader(ctx context.Context, key, value string) {
	if bag, ok := ctx.Value(bffRespHeadersKey{}).(*responseHeadersBag); ok {
		bag.set(key, value)
	}
}

// ResponseHeadersFromCtx returns all headers accumulated by CaptureResponseHeader,
// or nil when no bag is present in ctx.
// Typically called by BffHeadersMiddleware after c.Next() to write headers to the response.
func ResponseHeadersFromCtx(ctx context.Context) map[string]string {
	if bag, ok := ctx.Value(bffRespHeadersKey{}).(*responseHeadersBag); ok {
		return bag.all()
	}
	return nil
}
