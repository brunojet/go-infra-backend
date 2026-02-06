package exporters

import "os"

const (
	OTLPEndpointEnv     = "OTEL_EXPORTER_OTLP_ENDPOINT"
	OTLPEndpointDefault = "localhost:4317"
)

type OTLPExporterConfig struct {
	Endpoint string
}

func GetExporterConfigFromEnv() OTLPExporterConfig {
	if endpoint, exists := os.LookupEnv(OTLPEndpointEnv); exists && endpoint != "" {
		return OTLPExporterConfig{Endpoint: endpoint}
	}

	return OTLPExporterConfig{Endpoint: ""}
}
