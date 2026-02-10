package exporters

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/brunojet/go-infra-backend/internal/config"
)

const (
	OTLPEndpointEnv          = "OTEL_EXPORTER_OTLP_ENDPOINT"
	otlpEndpointDefault      = "localhost:4317"
	otlpLoggerEndpointKey    = "OTEL_EXPORTER_OTLP_ENDPOINT_LOGGER"
	otlpMetricEndpointKey    = "OTEL_EXPORTER_OTLP_ENDPOINT_METRIC"
	otlpTraceEndpointKey     = "OTEL_EXPORTER_OTLP_ENDPOINT_TRACE"
	otlpConsoleLogEnv        = "OTEL_EXPORTER_OTLP_CONSOLE_LOGGER_ENABLE"
	consoleLoggerRedirectEnv = "CONSOLE_LOGGER_REDIRECT_ENABLE"
)

var reHostPort = regexp.MustCompile(
	`^(?:` +
		`(` +
		`localhost` +
		`|(?:[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?)(?:\.(?:[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?))*` + // hostname
		`|(?:\d{1,3}\.){3}\d{1,3}` + // IPv4 (range não validado aqui)
		`|\[[0-9A-Fa-f:]+\]` + // IPv6 entre colchetes (ex.: [::1])
		`)?` +
		`:` +
		`)?` +
		`(\d{1,5})$`,
)

type OTLPConfig struct {
	Endpoint   string
	IsInsecure bool
}

type OTLPLoggerConfig struct {
	OTLPConfig
	ConsoleLogEnable      bool
	ConsoleRedirectEnable bool
}

func getEndpointConfigFromEnv(endpoint string) OTLPConfig {
	endPoint := config.GetEnv(endpoint, OTLPEndpointEnv)

	hasHttps := strings.Contains(endPoint, "https://")
	hasHttp := strings.Contains(endPoint, "http://")
	IsInsecure := !hasHttps

	if hasHttps || hasHttp {
		if u, err := url.Parse(endPoint); err == nil {
			return OTLPConfig{
				Endpoint:   u.Host,
				IsInsecure: IsInsecure,
			}
		}
	} else if reHostPort.MatchString(endPoint) {
		return OTLPConfig{
			Endpoint:   endPoint,
			IsInsecure: IsInsecure,
		}
	}

	return OTLPConfig{}
}

func getLoggerExporterConfigFromEnv() OTLPLoggerConfig {
	exporterConfig := getEndpointConfigFromEnv(otlpLoggerEndpointKey)
	consoleLogEnable := config.GetEnvAsBool(otlpConsoleLogEnv, exporterConfig.Endpoint == "")
	consoleRedirectEnable := config.GetEnvAsBool(consoleLoggerRedirectEnv, false)
	return OTLPLoggerConfig{
		OTLPConfig:            exporterConfig,
		ConsoleLogEnable:      consoleLogEnable,
		ConsoleRedirectEnable: consoleLogEnable && consoleRedirectEnable,
	}
}

func getMetricExporterConfigFromEnv() OTLPConfig {
	return getEndpointConfigFromEnv(otlpMetricEndpointKey)
}

func getTraceExporterConfigFromEnv() OTLPConfig {
	return getEndpointConfigFromEnv(otlpTraceEndpointKey)
}
