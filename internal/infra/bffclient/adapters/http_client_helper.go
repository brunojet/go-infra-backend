package adapters

import (
	"github.com/brunojet/go-infra-backend/debugassert"
	"github.com/brunojet/go-infra-backend/pkg/infra/bffclient/contracts"
)

func SetBffClientHeaders(clientHeaders map[string]string, requestStream contracts.BffHttpRequestStream) {
	debugassert.Assert(clientHeaders != nil, "SetBffClientHeaders: clientHeaders must not be nil")
	debugassert.Assert(requestStream != nil, "SetBffClientHeaders: requestStream must not be nil")
	for key, value := range clientHeaders {
		requestStream.SetHeader(key, value)
	}
}
