package bootstrap

import (
	"log"

	bootcontracts "github.com/brunojet/go-infra-backend/internal/bootstrap/contracts"
	obsexporters "github.com/brunojet/go-infra-backend/internal/observability/exporters"
	obsproviders "github.com/brunojet/go-infra-backend/internal/observability/providers"
)

func InitObservability(sm bootcontracts.ShutdownManager) {
	ctx := sm.GetContext()
	spanExporter, err := obsexporters.NewOTLPTracerExporter(ctx)
	if err != nil {
		log.Fatalf("failed to create otlp tracer exporter: %v", err)
	}
	_, tracerShutdown, err := obsproviders.NewOTLPTracerProvider(ctx, spanExporter)
	if err != nil {
		log.Fatalf("failed to create otlp tracer provider: %v", err)
	}
	sm.RegisterFunc("otel-tracer", tracerShutdown)

	metricExporter, err := obsexporters.NewOTLPMetricExporter(ctx)
	if err != nil {
		log.Fatalf("failed to create otlp metric exporter: %v", err)
	}
	_, metricShutdown, err := obsproviders.NewOTLPMetricProvider(ctx, metricExporter)
	if err != nil {
		log.Fatalf("failed to create otlp metric provider: %v", err)
	}
	sm.RegisterFunc("otel-metric", metricShutdown)

	loggerExporter, err := obsexporters.NewOTLPLoggerExporter(ctx)
	if err != nil {
		log.Fatalf("failed to create otlp logger exporter: %v", err)
	}
	_, loggerShutdown, err := obsproviders.NewOTLPLoggerProvider(ctx, loggerExporter)
	if err != nil {
		log.Fatalf("failed to create otlp logger provider: %v", err)
	}
	sm.RegisterFunc("otel-logger", loggerShutdown)
}
