package streams

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	bffcts "github.com/brunojet/go-infra-backend/pkg/infra/bffclient/contracts"
)

// ---------------------------------------------------------------------------
// JsonRequestStream
// ---------------------------------------------------------------------------

// JsonRequestStream implements BffHttpRequestStream for JSON payloads.
//
// The body is serialised at construction time; Reader() returns a fresh
// bytes.Reader on each call, making the stream safe for retry loops.
// Use NewJsonRequest when the request carries a body (POST, PATCH, PUT).
// Use NewJsonNoBodyRequest for body-less requests (GET, DELETE).
type JsonRequestStream struct {
	httpRequestBase
	body []byte // nil = no body
}

// NewJsonRequest creates a JsonRequestStream that serialises body to JSON.
// method is typically http.MethodPost, http.MethodPatch, or http.MethodPut.
func NewJsonRequest[T any](method, path string, params map[string]string, body T) (*JsonRequestStream, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("bffclient/streams: marshal request body: %w", err)
	}
	return &JsonRequestStream{httpRequestBase: newHttpRequestBase(method, path, params), body: data}, nil
}

// NewJsonNoBodyRequest creates a JsonRequestStream without a body.
// Use for GET and DELETE requests.
func NewJsonNoBodyRequest(method, path string, params map[string]string) *JsonRequestStream {
	return &JsonRequestStream{httpRequestBase: newHttpRequestBase(method, path, params)}
}

// ContentType returns "application/json" when a body is present, empty otherwise.
func (s *JsonRequestStream) ContentType() string {
	if s.body == nil {
		return ""
	}
	return "application/json"
}

// Reader returns a fresh bytes.Reader on each call.
// Returns nil when there is no body (GET, DELETE).
func (s *JsonRequestStream) Reader() io.Reader {
	if s.body == nil {
		return nil
	}
	return bytes.NewReader(s.body)
}

// ---------------------------------------------------------------------------
// JsonResponseStream
// ---------------------------------------------------------------------------

// JsonResponseStream[T] implements BffHttpResponseStream for JSON payloads.
//
// After a successful Emit call, Value holds the deserialised response body.
// T may be a slice (e.g. JsonResponseStream[[]MyEntity]) for list responses.
type JsonResponseStream[T any] struct {
	httpResponseBase
	Value T
}

// Decode deserialises the JSON response body into Value.
// If the injected status code indicates an upstream error (>= 400),
// the body is read raw and returned as *BffUpstreamError — never decoded as T.
func (s *JsonResponseStream[T]) Decode(r io.Reader) error {
	if s.statusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(r)
		return &bffcts.BffUpstreamError{StatusCode: s.statusCode, Body: body}
	}
	if err := json.NewDecoder(r).Decode(&s.Value); err != nil {
		return fmt.Errorf("bffclient/streams: decode response: %w", err)
	}
	return nil
}
