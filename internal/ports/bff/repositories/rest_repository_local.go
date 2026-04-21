package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/brunojet/go-infra-backend/debugassert"
	"github.com/brunojet/go-infra-backend/internal/ports/errors"
	"github.com/brunojet/go-infra-backend/pkg/ports/bff/repositories/contracts"
)

const (
	defaultErrorLenLimit = 4 * 1024 // 4KB
)

func (r *restRepositoryImpl[DS, MS]) makeURL(instancePath string, parentsPath string, pathParents int, parentIds ...string) url.URL {
	debugassert.Assert(instancePath != "", "instancePath must not be empty")
	debugassert.Assert(len(parentIds) == pathParents, "number of parent IDs must match pathParents")
	pathFmt := path.Join(instancePath, parentsPath)
	args := make([]interface{}, len(parentIds))
	for i, id := range parentIds {
		args[i] = url.PathEscape(id)
	}
	rel := fmt.Sprintf(pathFmt, args...)
	u := r.URL
	u.Path = path.Join(u.Path, rel)
	return u
}

func (r *restRepositoryImpl[DS, MS]) makeCollectionURL(parentIds ...string) url.URL {
	return r.makeURL(
		r.pathConfig.instancePath,
		r.pathConfig.collectionParentsFmt,
		r.pathConfig.collectionParents,
		parentIds...)
}

func (r *restRepositoryImpl[DS, MS]) makeInstanceURL(id string) url.URL {
	return r.makeURL(r.pathConfig.instancePath+"/%s", "", 1, id)
}

func (r *restRepositoryImpl[DS, MS]) extractFromEnvelope(response *contracts.RestResponse[DS], meta *MS, httpResponse *http.Response) error {
	debugassert.Assert(r.envelopeSpec != nil, "EnvelopeSpec must be configured for repository to extract from envelop")
	debugassert.Assert(response != nil, "response must not be nil")
	debugassert.Assert(httpResponse != nil, "httpResponse must not be nil")
	spec := r.envelopeSpec.New()
	if err := jsonDecoder(httpResponse.Body, spec); err != nil {
		return err
	}
	if err := json.Unmarshal(spec.GetData(), &response.Body); err != nil {
		return newEncoderError("failed to unmarshal Data field into target struct: %w", err)
	}
	if meta != nil {
		if err := json.Unmarshal(spec.GetMeta(), meta); err != nil {
			return newEncoderError("failed to unmarshal Meta field into target struct: %w", err)
		}
	}
	return nil
}

func (r *restRepositoryImpl[DS, MS]) setRestResponse(response *contracts.RestResponse[DS], httpResponse *http.Response) error {
	if httpResponse.StatusCode == http.StatusNoContent {
		return nil
	}
	if r.envelopeSpec != nil {
		return r.extractFromEnvelope(response, nil, httpResponse)
	}
	if err := jsonDecoder(httpResponse.Body, &response.Body); err != nil {
		return err
	}
	return nil
}

func (r *restRepositoryImpl[DS, MS]) setRestResponses(response *contracts.RestResponse[DS], meta *MS, httpResponse *http.Response) error {
	if r.envelopeSpec != nil {
		return r.extractFromEnvelope(response, meta, httpResponse)
	}
	if err := jsonDecoder(httpResponse.Body, &response.Bodies); err != nil {
		return err
	}
	return nil
}

func (r *restRepositoryImpl[DS, MS]) setOtherResponse(response *contracts.RestResponse[DS], httpResponse *http.Response) error {
	if r.messageExtractor != nil {
		message, err := r.messageExtractor(httpResponse)
		if err != nil {
			return errors.NewRestError(
				http.StatusInternalServerError,
				fmt.Errorf("failed to extract message from upstream response: %w", err),
			)
		}
		if message != "" {
			response.Message = message
			return nil
		}
		response.Message = http.StatusText(httpResponse.StatusCode)
		return nil
	}

	if err := messageExtractorLimited(response, httpResponse, defaultErrorLenLimit); err != nil {
		return errors.NewRestError(
			http.StatusInternalServerError,
			fmt.Errorf("failed to read response body: %w", err),
		)
	}
	return nil
}

func (r *restRepositoryImpl[DS, MS]) handleRestResponse(response *contracts.RestResponse[DS], httpResponse *http.Response) error {
	if hasBody := setRestResponseOptions(&response.Opts, httpResponse); !hasBody {
		return nil
	}
	switch httpResponse.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		return r.setRestResponse(response, httpResponse)
	default:
		return r.setOtherResponse(response, httpResponse)
	}
}

func (r *restRepositoryImpl[DS, MS]) handleRestResponses(response *contracts.RestResponse[DS], meta *MS, httpResponse *http.Response) error {
	if hasBody := setRestResponseOptions(&response.Opts, httpResponse); !hasBody {
		return nil
	}
	switch httpResponse.StatusCode {
	case http.StatusOK, http.StatusPartialContent:
		return r.setRestResponses(response, meta, httpResponse)
	default:
		return r.setOtherResponse(response, httpResponse)
	}
}

func (r *restRepositoryImpl[DS, MS]) makeCollectionRestRequest(method string, opts *contracts.RestRequestOptions, parentIds ...string) http.Request {
	url := r.makeCollectionURL(parentIds...)
	httpReq := makeRestRequestURL(method, url)
	setRestRequestOptions(&httpReq, opts)
	return httpReq
}

func (r *restRepositoryImpl[DS, MS]) makeInstanceRestRequest(method string, opts *contracts.RestRequestOptions, id string) http.Request {
	url := r.makeInstanceURL(id)
	httpReq := makeRestRequestURL(method, url)
	setRestRequestOptions(&httpReq, opts)
	return httpReq
}

// doSingleRequest executa o padrão comum para operações que retornam um único item
// (GET instance, POST create, PUT update, PATCH save, DELETE)
func (r *restRepositoryImpl[DS, MS]) doSingleRequest(
	ctx context.Context,
	makeRequest func() http.Request,
	body any, // nil se não precisa de body
	response *contracts.RestResponse[DS],
) error {
	debugassert.Assert(response != nil, "response must not be nil")
	httpReq := makeRequest()
	if body != nil {
		if err := setRestRequestBody(&httpReq, body); err != nil {
			return err
		}
	}
	httpResp, err := r.adapter.Do(ctx, &httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()
	return r.handleRestResponse(response, httpResp)
}

// doListRequest executa o padrão comum para operações que retornam lista (GET collection)
func (r *restRepositoryImpl[DS, MS]) doListRequest(
	ctx context.Context,
	makeRequest func() http.Request,
	response *contracts.RestResponse[DS],
	meta *MS,
) error {
	debugassert.Assert(response != nil, "response must not be nil")
	httpReq := makeRequest()
	httpResp, err := r.adapter.Do(ctx, &httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()
	return r.handleRestResponses(response, meta, httpResp)
}
