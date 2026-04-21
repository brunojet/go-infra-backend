package repositories

import (
	"context"
	"net/http"

	"github.com/brunojet/go-infra-backend/pkg/infra/http_client"
	"github.com/brunojet/go-infra-backend/pkg/ports/bff/repositories/contracts"
)

type restRepositoryImpl[DS, MS any] struct {
	adapter http_client.HttpClientAdapter
	restRepositoryConfig
}

func NewRestRepository[DS, MS any](adapter http_client.HttpClientAdapter, opts ...RestRepositoryOption) contracts.RestRepository[DS, MS] {
	return &restRepositoryImpl[DS, MS]{
		adapter:              adapter,
		restRepositoryConfig: newRestConfig(opts...),
	}
}

// Create implements [contracts.RestRepository].
func (r *restRepositoryImpl[DS, MS]) Create(ctx context.Context, request contracts.RestRequest, response *contracts.RestResponse[DS], parentIds ...string) error {
	return r.doSingleRequest(ctx,
		func() http.Request {
			return r.makeCollectionRestRequest(http.MethodPost, &request.Opts, parentIds...)
		},
		request.Body,
		response,
	)
}

// List implements [contracts.RestRepository].
func (r *restRepositoryImpl[DS, MS]) List(ctx context.Context, opts *contracts.RestRequestOptions, response *contracts.RestResponse[DS], metastream *MS, parentIds ...string) error {
	return r.doListRequest(ctx,
		func() http.Request {
			return r.makeCollectionRestRequest(http.MethodGet, opts, parentIds...)
		},
		response,
		metastream,
	)
}

// GetById implements [contracts.RestRepository].
func (r *restRepositoryImpl[DS, MS]) GetById(ctx context.Context, opts *contracts.RestRequestOptions, response *contracts.RestResponse[DS], id string) error {
	return r.doSingleRequest(ctx,
		func() http.Request {
			return r.makeInstanceRestRequest(http.MethodGet, opts, id)
		},
		nil,
		response,
	)
}

// Update implements [contracts.RestRepository].
func (r *restRepositoryImpl[DS, MS]) Update(ctx context.Context, request contracts.RestRequest, response *contracts.RestResponse[DS], id string) error {
	return r.doSingleRequest(ctx,
		func() http.Request {
			return r.makeInstanceRestRequest(http.MethodPut, &request.Opts, id)
		},
		request.Body,
		response,
	)
}

// Save implements [contracts.RestRepository].
func (r *restRepositoryImpl[DS, MS]) Save(ctx context.Context, request contracts.RestRequest, response *contracts.RestResponse[DS], id string) error {
	return r.doSingleRequest(ctx,
		func() http.Request {
			return r.makeInstanceRestRequest(http.MethodPatch, &request.Opts, id)
		},
		request.Body,
		response,
	)
}

// Delete implements [contracts.RestRepository].
func (r *restRepositoryImpl[DS, MS]) Delete(ctx context.Context, opts *contracts.RestRequestOptions, response *contracts.RestResponse[DS], id string) error {
	return r.doSingleRequest(ctx,
		func() http.Request {
			return r.makeInstanceRestRequest(http.MethodDelete, opts, id)
		},
		nil,
		response,
	)
}
