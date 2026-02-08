package exporters

import internalexporters "github.com/brunojet/go-infra-backend/internal/observability/exporters"

// Re-export the env var key so external apps can configure OTLP without touching internal.
const OTLPEndpointEnv = internalexporters.OTLPEndpointEnv
