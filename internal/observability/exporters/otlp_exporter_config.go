package exporters

import (
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

type OTLPLoggerExporterConfig struct {
	Endpoint              string
	EndpointLogEnable     bool
	ConsoleLogEnable      bool
	ConsoleRedirectEnable bool
}

type OTLPMetricExporterConfig struct {
	Endpoint string
}

type OTLPTracerExporterConfig struct {
	Endpoint string
}

type OTLPExporterConfig struct {
	Endpoint     string
	loggerConfig OTLPLoggerExporterConfig
	metricConfig OTLPMetricExporterConfig
	tracerConfig OTLPTracerExporterConfig
}

func getLoggerExporterConfigFromEnv(globalEndpoint string) OTLPLoggerExporterConfig {
	endPoint := config.GetEnv(otlpLoggerEndpointKey, globalEndpoint)
	consoleRedirectEnable := config.GetEnvAsBool(consoleLoggerRedirectEnv, false)
	return OTLPLoggerExporterConfig{
		Endpoint:              endPoint,
		EndpointLogEnable:     endPoint != "",
		ConsoleLogEnable:      config.GetEnvAsBool(otlpConsoleLogEnv, endPoint == "" || consoleRedirectEnable),
		ConsoleRedirectEnable: consoleRedirectEnable,
	}
}

func getMetricExporterConfigFromEnv(globalEndpoint string) OTLPMetricExporterConfig {
	return OTLPMetricExporterConfig{
		Endpoint: config.GetEnv(otlpMetricEndpointKey, globalEndpoint),
	}
}

func getTraceExporterConfigFromEnv(globalEndpoint string) OTLPTracerExporterConfig {
	return OTLPTracerExporterConfig{
		Endpoint: config.GetEnv(otlpTraceEndpointKey, globalEndpoint),
	}
}

func GetExporterConfigFromEnv() OTLPExporterConfig {
	globalEndpoint := config.GetEnv(OTLPEndpointEnv, "")
	return OTLPExporterConfig{
		Endpoint:     globalEndpoint,
		loggerConfig: getLoggerExporterConfigFromEnv(globalEndpoint),
		metricConfig: getMetricExporterConfigFromEnv(globalEndpoint),
		tracerConfig: getTraceExporterConfigFromEnv(globalEndpoint),
	}
}
