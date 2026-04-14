package contracts

import (
	"context"
	"net/http"
)

type HttpClientAdapter interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}
