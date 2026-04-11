package adapters

import (
	"net"
	"net/url"
	"strconv"
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

func isHostPort(s string) bool {
	// First try net.SplitHostPort which handles host:port and [ipv6]:port
	if _, port, err := net.SplitHostPort(s); err == nil {
		if p, err := strconv.Atoi(port); err == nil && p > 0 && p <= 65535 {
			return true
		}
		return false
	}
	// Fallback: allow numeric-only port like "4317"
	if p, err := strconv.Atoi(s); err == nil && p > 0 && p <= 65535 {
		return true
	}
	return false
}

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
	// prefer endpoint-specific env var, then global OTLP env, then default
	endPoint := config.GetEnv(endpoint, "")
	if endPoint == "" {
		endPoint = config.GetEnv(OTLPEndpointEnv, otlpEndpointDefault)
	}

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
	} else if isHostPort(endPoint) {
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
