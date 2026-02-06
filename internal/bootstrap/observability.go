package bootstrap

import (
	"log"

	bootcontracts "github.com/brunojet/go-infra-backend/internal/bootstrap/contracts"
	obs "github.com/brunojet/go-infra-backend/internal/observability"
	obsexporters "github.com/brunojet/go-infra-backend/internal/observability/exporters"
	obsproviders "github.com/brunojet/go-infra-backend/internal/observability/providers"
	otellog "go.opentelemetry.io/otel/log"
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

	otlpLoggerExporter, err := obsexporters.NewOTLPLoggerExporter(ctx)
	if err != nil {
		log.Fatalf("failed to create otlp logger exporter: %v", err)
	}
	consoleLoggerExporter, err := obsexporters.NewConsoleLoggerExporter()
	if err != nil {
		log.Fatalf("failed to create console logger exporter: %v", err)
	}
	_, loggerShutdown, err := obsproviders.NewOTLPLoggerProvider(ctx, otlpLoggerExporter, consoleLoggerExporter)
	if err != nil {
		log.Fatalf("failed to create otlp logger provider: %v", err)
	}
	// Redirect all stdlib log.Printf/log.Println output to the OTel logger pipeline.
	// This makes packages that still use `import "log"` automatically emit via OTel.
	obs.RedirectStdLog("stdlib", otellog.SeverityInfo)
	sm.RegisterFunc("otel-logger", loggerShutdown)
}
