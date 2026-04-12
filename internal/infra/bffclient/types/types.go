package types

import "sync"

type BffRequestHeadersKey struct{}

type BffResponseHeadersKey struct{}

type ResponseHeadersBag struct {
	mu sync.Mutex
	m  map[string]string
}

func NewResponseHeadersBag(keys []string) *ResponseHeadersBag {
	m := make(map[string]string, len(keys))
	for _, k := range keys {
		m[k] = ""
	}
	return &ResponseHeadersBag{m: m}
}

func (b *ResponseHeadersBag) Set(k, v string) {
	b.mu.Lock()
	b.m[k] = v
	b.mu.Unlock()
}

// SetKnown updates k only if it was pre-registered via NewResponseHeadersBag.
// This lets the bag act as its own allow-list: callers pass all upstream headers
// and SetKnown silently ignores keys that were not configured.
func (b *ResponseHeadersBag) SetKnown(k, v string) {
	b.mu.Lock()
	if _, exists := b.m[k]; exists {
		b.m[k] = v
	}
	b.mu.Unlock()
}

func (b *ResponseHeadersBag) All() map[string]string {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make(map[string]string, len(b.m))
	for k, v := range b.m {
		out[k] = v
	}
	return out
}
