package contracts

import (
	"context"
	"net/http"
	"net/url"
)

type RestRequestOptions struct {
	http.Header
	url.Values
}

type RestResponseOptions struct {
	http.Header
	StatusCode int
}

type RestRequest struct {
	Opts RestRequestOptions
	Body any
}

type RestResponse[DS any] struct {
	Opts    RestResponseOptions
	Body    DS
	Bodies  []DS
	Message string
}
type RestRepository[DS, MS any] interface {
	Create(ctx context.Context, request RestRequest, response *RestResponse[DS], parentIds ...string) error
	List(ctx context.Context, opts *RestRequestOptions, response *RestResponse[DS], metastream *MS, parentIds ...string) error
	GetById(ctx context.Context, opts *RestRequestOptions, response *RestResponse[DS], id string) error
	Update(ctx context.Context, request RestRequest, response *RestResponse[DS], id string) error
	Save(ctx context.Context, request RestRequest, response *RestResponse[DS], id string) error
	Delete(ctx context.Context, opts *RestRequestOptions, response *RestResponse[DS], id string) error
}
