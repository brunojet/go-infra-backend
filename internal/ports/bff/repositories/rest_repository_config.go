package repositories

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type RestRepositoryOption func(*restRepositoryConfig)

type RestMessageExtractorFunc func(httpResponse *http.Response) (string, error)

type PathConfig struct {
	instancePath         string // resource path segment used by instance and collection operations, e.g. "orders"
	collectionParentsFmt string // internal fmt path generated from collection parent resource segments, e.g. "users/%s/projects/%s"
	collectionParents    int    // number of parent IDs expected for collection operations
}

func (c PathConfig) Validate() {
	if c.instancePath == "" {
		panic("instance path must be provided")
	}
}

type DefaultEnvelopeSpec struct {
	Data json.RawMessage `json:"data"`
	Meta json.RawMessage `json:"meta"`
}

func (e *DefaultEnvelopeSpec) GetData() json.RawMessage { return e.Data }
func (e *DefaultEnvelopeSpec) GetMeta() json.RawMessage { return e.Meta }
func (e *DefaultEnvelopeSpec) New() EnvelopeSpec        { return &DefaultEnvelopeSpec{} }

type EnvelopeSpec interface {
	New() EnvelopeSpec
	GetData() json.RawMessage
	GetMeta() json.RawMessage
}
type restRepositoryConfig struct {
	URL              url.URL
	pathConfig       PathConfig
	envelopeSpec     EnvelopeSpec
	messageExtractor RestMessageExtractorFunc
}

func (c *restRepositoryConfig) Validate() {
	if c.URL == (url.URL{}) {
		panic("base URL must be provided")
	}
	c.pathConfig.Validate()
}

func WithBaseURL(baseURL string) RestRepositoryOption {
	return func(cfg *restRepositoryConfig) {
		parsedURL, err := url.Parse(baseURL)
		if err != nil {
			panic(fmt.Sprintf("invalid base URL: %v", err))
		}
		cfg.URL = *parsedURL
	}
}

func WithPathConfig(instancePath string, collectionParents ...string) RestRepositoryOption {
	return func(cfg *restRepositoryConfig) {
		if instancePath == "" {
			panic("instance path must be provided")
		}
		cfg.pathConfig = PathConfig{
			instancePath: instancePath,
		}
		if len(collectionParents) > 0 {
			collectionParentsFmt, collectionParentsCount := buildCollectionParentsFmt(collectionParents...)
			cfg.pathConfig.collectionParentsFmt = collectionParentsFmt
			cfg.pathConfig.collectionParents = collectionParentsCount
		}
	}
}

// WithMessageExtractor configures a callback to extract a human-readable message
// from upstream non-success response bodies.
// The callback receives the full HTTP response and owns body consumption.
// When this callback is configured, the default truncation fallback is bypassed.
func WithMessageExtractor(extractor RestMessageExtractorFunc) RestRepositoryOption {
	return func(cfg *restRepositoryConfig) {
		cfg.messageExtractor = extractor
	}
}

// WithEnvelopeSpec allows configuring a custom envelope structure for responses. When configured, the
// decode will use this type to unmarshal and extract Data/Meta fields.
func WithEnvelopeSpec(spec EnvelopeSpec) RestRepositoryOption {
	return func(cfg *restRepositoryConfig) {
		cfg.envelopeSpec = spec
	}
}

func newRestConfig(opts ...RestRepositoryOption) restRepositoryConfig {
	cfg := restRepositoryConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}
	cfg.Validate()
	return cfg
}
