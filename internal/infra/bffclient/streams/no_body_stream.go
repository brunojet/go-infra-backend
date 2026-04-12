package streams

import "io"

// ---------------------------------------------------------------------------
// NoBodyResponseStream
// ---------------------------------------------------------------------------

// NoBodyResponseStream implements BffHttpResponseStream for operations that
// produce no meaningful response body (e.g. DELETE returning 204 No Content).
//
// Decode discards any body the upstream sends; HTTP-level errors are detected
// by the adapter through the status code check after Decode returns.
type NoBodyResponseStream struct {
	httpResponse
}

// NewNoBodyResponseStream creates a NoBodyResponseStream.
func NewNoBodyResponseStream() *NoBodyResponseStream {
	return &NoBodyResponseStream{httpResponse: NewHttpResponse(0)}
}

// Decode discards the response body and always returns nil.
func (s *NoBodyResponseStream) Decode(r io.Reader) error {
	return nil
}
