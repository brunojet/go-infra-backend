package streams

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
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
	httpRequest
	body []byte // nil = no body
}

// NewJsonRequest creates a JsonRequestStream that serialises body to JSON.
// method is typically http.MethodPost, http.MethodPatch, or http.MethodPut.
func NewJsonRequest[T any](method, path string, params map[string]string, body T) (*JsonRequestStream, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, wrapErr(errMarshalRequestBody, err)
	}
	req := &JsonRequestStream{httpRequest: NewHttpRequest(method, path, params), body: data}
	req.SetHeader(headerContentType, contentTypeJSON)
	return req, nil
}

// NewJsonNoBodyRequest creates a JsonRequestStream without a body.
// Use for GET and DELETE requests.
func NewJsonNoBodyRequest(method, path string, params map[string]string) *JsonRequestStream {
	return &JsonRequestStream{httpRequest: NewHttpRequest(method, path, params)}
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
// RawBody is populated when the upstream responds with a non-JSON Content-Type
// on a 2xx response; Value will be zero in that case.
type JsonResponseStream[T any] struct {
	httpResponse
	Value   T
	RawBody []byte
}

// Decode deserialises the JSON response body into Value.
// Returns errUnexpectedNoContent if statusCode is 204 — use NoBodyResponseStream instead.
// If the upstream responds with a non-JSON Content-Type on a 2xx, the raw body
// is stored in RawBody and Value is left as its zero value.
func (s *JsonResponseStream[T]) Decode(r io.Reader) error {
	if s.statusCode == http.StatusNoContent {
		return errUnexpectedNoContent
	}
	if ct := s.Headers().Get(headerContentType); !strings.HasPrefix(ct, contentTypeJSON) {
		s.RawBody, _ = io.ReadAll(r)
		return nil
	}
	if err := json.NewDecoder(r).Decode(&s.Value); err != nil {
		return wrapErr(errDecodeResponse, err)
	}
	return nil
}
