package repositories

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	porterrors "github.com/brunojet/go-infra-backend/internal/ports/errors"
	"github.com/brunojet/go-infra-backend/pkg/ports/bff/repositories/contracts"
)

func TestSetOtherResponse_UsesCustomExtractor(t *testing.T) {
	repo := &restRepositoryImpl[struct{}, struct{}]{
		restRepositoryConfig: restRepositoryConfig{
			messageExtractor: func(resp *http.Response) (string, error) {
				data, err := io.ReadAll(resp.Body)
				if err != nil {
					return "", err
				}
				return "custom:" + string(data), nil
			},
		},
	}

	httpResp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader("payload")),
	}
	resp := &contracts.RestResponse[struct{}]{Opts: contracts.RestResponseOptions{Header: httpResp.Header}}

	err := repo.setOtherResponse(resp, httpResp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Message != "custom:payload" {
		t.Fatalf("unexpected message: %q", resp.Message)
	}
}

func TestSetOtherResponse_ExtractorErrorIsPropagatedAsRestError(t *testing.T) {
	repo := &restRepositoryImpl[struct{}, struct{}]{
		restRepositoryConfig: restRepositoryConfig{
			messageExtractor: func(_ *http.Response) (string, error) {
				return "", errors.New("parse failed")
			},
		},
	}

	httpResp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader("payload")),
	}
	resp := &contracts.RestResponse[struct{}]{Opts: contracts.RestResponseOptions{Header: httpResp.Header}}

	err := repo.setOtherResponse(resp, httpResp)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !porterrors.IsRestError(err) {
		t.Fatalf("expected RestError, got %T", err)
	}
	if status, ok := porterrors.GetRestErrorStatusCode(err); !ok || status != http.StatusInternalServerError {
		t.Fatalf("unexpected rest error status: %d, ok=%v", status, ok)
	}
}

func TestSetOtherResponse_FallbackLimitedExtractor(t *testing.T) {
	repo := &restRepositoryImpl[struct{}, struct{}]{}
	httpResp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(strings.Repeat("z", defaultErrorLenLimit+1))),
	}
	resp := &contracts.RestResponse[struct{}]{Opts: contracts.RestResponseOptions{Header: httpResp.Header}}

	err := repo.setOtherResponse(resp, httpResp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Message) != defaultErrorLenLimit {
		t.Fatalf("expected truncated message length %d, got %d", defaultErrorLenLimit, len(resp.Message))
	}
	if got := resp.Opts.Header.Get("X-Raw-Body-Truncated"); got != "1" {
		t.Fatalf("expected truncation header, got %q", got)
	}
}
