package streams

import (
	"net/http"
	"net/url"
	"strings"
)

// httpRequest carries the HTTP routing metadata shared by all HTTP request
// streams: method, path, and query parameters.
//
// Embed this struct in concrete request stream types to inherit Method(), Path(),
// and Params() without repetition. Each embedding type only needs to override
// the format-specific methods: ContentType() and Reader().
type httpRequest struct {
	http.Header
	method    string
	path      string
	rawParams string
}

func buildRawQuery(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}
	q := make(url.Values, len(params))
	for k, v := range params {
		q.Set(k, v)
	}
	return q.Encode()
}

func NewHttpRequest(method, path string, params map[string]string) httpRequest {
	return httpRequest{Header: make(http.Header), method: method, path: path, rawParams: buildRawQuery(params)}
}

func (b *httpRequest) SetHeader(key, value string) {
	b.Header.Set(key, value)
}

func (b *httpRequest) MergeHeader(h http.Header) {
	for k, v := range h {
		b.Header.Set(k, strings.Join(v, ", "))
	}
}

func (b httpRequest) Headers() http.Header {
	return b.Header
}

func (b httpRequest) Method() string {
	return b.method
}

func (b httpRequest) Path() string {
	return b.path
}

func (b httpRequest) RawQuery() string {
	return b.rawParams
}

// ---------------------------------------------------------------------------
// httpResponseBase
// ---------------------------------------------------------------------------

// httpResponse carries the HTTP response metadata injected by the adapter
// before Decode is called: status code and response headers.
//
// Embed this struct in concrete response stream types (JsonResponseStream,
// FileDownloadStream, NoBodyResponseStream) to inherit SetStatusCode, SetHeaders
// and Headers without repetition. Decode reads statusCode to decide
// between the success path and BffUpstreamError.
type httpResponse struct {
	http.Header
	statusCode int
}

func NewHttpResponse(code int) httpResponse {
	return httpResponse{Header: make(http.Header), statusCode: code}
}

func (b *httpResponse) StatusCode() int {
	return b.statusCode
}

func (b *httpResponse) SetStatusCode(code int) {
	b.statusCode = code
}

func (b *httpResponse) SetHeader(key, value string) {
	if b.Header == nil {
		b.Header = make(http.Header)
	}
	b.Header.Set(key, value)
}

func (b *httpResponse) MergeHeader(h http.Header) {
	if b.Header == nil {
		b.Header = make(http.Header)
	}
	for k, v := range h {
		b.Header.Set(k, strings.Join(v, ", "))
	}
}

func (b *httpResponse) Headers() http.Header {
	if b.Header == nil {
		b.Header = make(http.Header)
	}
	return b.Header
}
