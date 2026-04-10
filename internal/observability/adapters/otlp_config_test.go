package adapters

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test_getTraceExporterConfigFromEnv_Default: quando não há variáveis configuradas,
// deve-se usar o valor default definido em `otlpEndpointDefault`.
func Test_getTraceExporterConfigFromEnv_Default(t *testing.T) {
	t.Setenv(otlpTraceEndpointKey, "")
	t.Setenv(OTLPEndpointEnv, "")

	cfg := getTraceExporterConfigFromEnv()
	require.Equal(t, otlpEndpointDefault, cfg.Endpoint)
	require.True(t, cfg.IsInsecure)
}

// Test_getTraceExporterConfigFromEnv_HTTPScheme: aceita endpoints com scheme https://
func Test_getTraceExporterConfigFromEnv_HTTPScheme(t *testing.T) {
	t.Setenv(otlpTraceEndpointKey, "https://collector:4317")

	cfg := getTraceExporterConfigFromEnv()
	require.Equal(t, "collector:4317", cfg.Endpoint)
	require.False(t, cfg.IsInsecure)
}

// Test_getTraceExporterConfigFromEnv_HostPort: aceita endpoints do tipo :4317
func Test_getTraceExporterConfigFromEnv_HostPort(t *testing.T) {
	t.Setenv(otlpTraceEndpointKey, ":4317")

	cfg := getTraceExporterConfigFromEnv()
	require.Equal(t, ":4317", cfg.Endpoint)
}

// Test_getLoggerExporterConfigFromEnv_ConsoleFlag: reconhece flag de console
func Test_getLoggerExporterConfigFromEnv_ConsoleFlag(t *testing.T) {
	t.Setenv(otlpConsoleLogEnv, "1")

	cfg := getLoggerExporterConfigFromEnv()
	require.True(t, cfg.ConsoleLogEnable)
}
