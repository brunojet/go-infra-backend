package repositories

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/brunojet/go-infra-backend/pkg/ports/bff/repositories/contracts"
)

func TestReadLimitedBody_NoTruncation(t *testing.T) {
	data, truncated, err := readLimitedBody(strings.NewReader("hello"), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if truncated {
		t.Fatalf("expected non-truncated body")
	}
	if string(data) != "hello" {
		t.Fatalf("unexpected data: %q", string(data))
	}
}

func TestReadLimitedBody_WithTruncation(t *testing.T) {
	input := strings.Repeat("a", 20)
	data, truncated, err := readLimitedBody(strings.NewReader(input), 8)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !truncated {
		t.Fatalf("expected truncation")
	}
	if len(data) != 8 {
		t.Fatalf("expected len(data)=8, got %d", len(data))
	}
}

func TestMessageExtractorLimited_SetsMessageAndTruncationHeader(t *testing.T) {
	httpResp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(strings.Repeat("x", defaultErrorLenLimit+50))),
	}
	resp := &contracts.RestResponse[struct{}]{
		Opts: contracts.RestResponseOptions{Header: httpResp.Header},
	}

	err := messageExtractorLimited(resp, httpResp, defaultErrorLenLimit)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Message) != defaultErrorLenLimit {
		t.Fatalf("expected message to be truncated to %d bytes, got %d", defaultErrorLenLimit, len(resp.Message))
	}
	if got := resp.Opts.Header.Get("X-Raw-Body-Truncated"); got != "1" {
		t.Fatalf("expected truncation header to be set, got %q", got)
	}
}

func TestMessageExtractorLimited_EmptyBodyFallsBackToStatusText(t *testing.T) {
	httpResp := &http.Response{
		StatusCode: http.StatusNotFound,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader("")),
	}
	resp := &contracts.RestResponse[struct{}]{
		Opts: contracts.RestResponseOptions{Header: httpResp.Header},
	}

	err := messageExtractorLimited(resp, httpResp, defaultErrorLenLimit)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Message != http.StatusText(http.StatusNotFound) {
		t.Fatalf("unexpected fallback message: %q", resp.Message)
	}
}
