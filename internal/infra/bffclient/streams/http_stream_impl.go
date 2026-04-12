package streams

import (
	"net/http"
	"net/url"
	"strings"
)

// httpRequestBase carries the HTTP routing metadata shared by all HTTP request
// streams: method, path, and query parameters.
//
// Embed this struct in concrete request stream types to inherit Method(), Path(),
// and Params() without repetition. Each embedding type only needs to override
// the format-specific methods: ContentType() and Reader().
type httpRequestBase struct {
	method    string
	path      string
	rawParams string
	headers   http.Header
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

func newHttpRequestBase(method, path string, params map[string]string) httpRequestBase {
	return httpRequestBase{method: method, path: path, rawParams: buildRawQuery(params)}
}

func (b httpRequestBase) Method() string {
	return b.method
}

func (b httpRequestBase) Path() string {
	return b.path
}

func (b httpRequestBase) RawQuery() string {
	return b.rawParams
}

func (b *httpRequestBase) SetHeader(key, value string) {
	b.headers.Set(key, value)
}

func (b *httpRequestBase) SetHeaders(h http.Header) {
	b.headers = make(http.Header, len(h))
	for k, v := range h {
		b.headers.Set(k, strings.Join(v, ", "))
	}
}

func (b *httpRequestBase) Headers() http.Header {
	return b.headers
}

// ---------------------------------------------------------------------------
// httpResponseBase
// ---------------------------------------------------------------------------

// httpResponseBase carries the HTTP response metadata injected by the adapter
// before Decode is called: status code and response headers.
//
// Embed this struct in concrete response stream types (JsonResponseStream,
// FileDownloadStream, NoBodyResponseStream) to inherit SetStatusCode, SetHeaders
// and Headers without repetition. Decode reads statusCode to decide
// between the success path and BffUpstreamError.
type httpResponseBase struct {
	statusCode int
	headers    http.Header
}

func (b *httpResponseBase) SetStatusCode(code int) {
	b.statusCode = code
}

func (b *httpResponseBase) SetHeader(key, value string) {
	b.headers.Set(key, value)
}

func (b *httpResponseBase) SetHeaders(h http.Header) {
	b.headers = make(http.Header, len(h))
	for k, v := range h {
		b.headers.Set(k, strings.Join(v, ", "))
	}
}
