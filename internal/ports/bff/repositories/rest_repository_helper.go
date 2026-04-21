package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/brunojet/go-infra-backend/debugassert"
	"github.com/brunojet/go-infra-backend/internal/ports/errors"
	"github.com/brunojet/go-infra-backend/pkg/ports/bff/repositories/contracts"
)

func newEncoderError(message string, args ...any) error {
	return errors.NewRestError(http.StatusInternalServerError, fmt.Errorf(message, args...))
}

func jsonEncoder(body any) (io.Reader, int, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(body); err != nil {
		return nil, 0, newEncoderError("failed to encode request body: %w", err)
	}
	return &buf, buf.Len(), nil
}

func jsonDecoder(r io.Reader, target any) error {
	dec := json.NewDecoder(r)
	if err := dec.Decode(target); err != nil {
		return newEncoderError("failed to decode JSON response: %w", err)
	}
	return nil
}

func makeRestRequestURL(method string, url url.URL, opts *contracts.RestRequestOptions) http.Request {
	httpReq := http.Request{}
	httpReq.Method = method
	httpReq.URL = &url
	if opts != nil && opts.Header != nil {
		httpReq.Header = opts.Header
	}
	if httpReq.Header == nil {
		httpReq.Header = make(http.Header)
	}
	if httpReq.Header.Get("Content-Type") == "" {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	return httpReq
}

func setRestRequestOptions(httpReq *http.Request, opts *contracts.RestRequestOptions) {
	debugassert.Assert(httpReq != nil, "httpReq must not be nil")
	if opts == nil {
		return
	}
	if opts.Header != nil {
		httpReq.Header = opts.Header
	}
	if opts.Values != nil {
		httpReq.URL.RawQuery = opts.Values.Encode()
	}
	if httpReq.Header.Get("Content-Type") == "" {
		httpReq.Header.Set("Content-Type", "application/json")
	}
}

func setRestResponseOptions(opts *contracts.RestResponseOptions, httpResponse *http.Response) (hasBody bool) {
	opts.StatusCode = httpResponse.StatusCode
	opts.Header = httpResponse.Header
	hasBody = httpResponse.StatusCode != http.StatusNoContent
	return
}

func setRestRequestBody(httpReq *http.Request, body any) error {
	debugassert.Assert(httpReq != nil, "httpReq must not be nil")
	debugassert.Assert(body != nil, "body must not be nil")
	reader, bodyLen, err := jsonEncoder(body)
	if err != nil {
		return err
	}
	httpReq.Body = io.NopCloser(reader)
	httpReq.ContentLength = int64(bodyLen)
	return nil
}

func setFieldFromEnvelop(envelop map[string]any, fieldName string, target any) error {
	data, ok := envelop[fieldName]
	if !ok {
		return newEncoderError("field '%s' not found in response envelop", fieldName)
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return newEncoderError("failed to marshal field '%s': %w", fieldName, err)
	}
	if err := json.Unmarshal(dataBytes, target); err != nil {
		return newEncoderError("failed to unmarshal field '%s' into target struct: %w", fieldName, err)
	}
	return nil
}

// buildCollectionParentsFmt builds the collection parent path format from
// resource segment names, e.g. ["users", "orders"] -> "users/%s/orders/%s".
func buildCollectionParentsFmt(collectionParents ...string) (string, int) {
	debugassert.Assert(len(collectionParents) > 0, "collection parents must be provided")

	parts := make([]string, 0, len(collectionParents)*2)
	for _, parent := range collectionParents {
		parent = strings.TrimSpace(parent)
		if parent == "" {
			panic("collection parent must not be empty")
		}
		parts = append(parts, parent, "%s")
	}

	return strings.Join(parts, "/"), len(collectionParents)
}

// readLimitedBody reads up to limit bytes from r and detects if body was truncated.
// Returns (data, truncated, error).
// This helper is extracted for testability: allows unit tests to inject io.Reader
// and verify truncation detection without coupling to HTTP response handling.
func readLimitedBody(r io.Reader, limit int) ([]byte, bool, error) {
	debugassert.Assert(limit > 0, "limit must be positive")

	// Read up to limit+1 to detect truncation
	limitReader := io.LimitReader(r, int64(limit)+1)
	data, err := io.ReadAll(limitReader)
	if err != nil {
		return nil, false, err
	}

	truncated := len(data) > limit
	if truncated {
		// Keep only up to limit bytes
		data = data[:limit]
	}

	return data, truncated, nil
}

// messageExtractorLimited is the default extractor used when no custom
// message extractor callback is configured. It applies truncation, updates
// response.Message and truncation header, and provides a status-text fallback
// when upstream body is empty.
func messageExtractorLimited[DS any](response *contracts.RestResponse[DS], httpResponse *http.Response, limit int) error {
	debugassert.Assert(response != nil, "response must not be nil")
	debugassert.Assert(httpResponse != nil, "httpResponse must not be nil")

	data, truncated, err := readLimitedBody(httpResponse.Body, limit)
	if err != nil {
		return err
	}

	message := string(data)
	if message == "" {
		message = http.StatusText(httpResponse.StatusCode)
	}
	response.Message = message

	if truncated {
		response.Opts.Header.Set("X-Raw-Body-Truncated", "1")
	}

	return nil
}
