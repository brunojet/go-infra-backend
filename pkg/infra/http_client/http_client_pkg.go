package http_client

import "github.com/brunojet/go-infra-backend/pkg/infra/http_client/contracts"

// HttpClientAdapter is a type alias that re-exports the interface
// defined in the contracts subpackage so callers can use
// `http_client.HttpClientAdapter`.
type HttpClientAdapter = contracts.HttpClientAdapter
