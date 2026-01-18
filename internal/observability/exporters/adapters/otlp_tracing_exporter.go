package adapters

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otlptrace "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"github.com/brunojet/go-infra-backend/internal/observability/exporters/contracts"
)

// OTLPTracingExporter é um adapter que implementa contracts.TracingExporter
// usando o provider OTLP.
type OTLPTracingExporter struct {
	endpoint  string
	setGlobal bool
	tp        interface{} // manter tipo opaco para evitar dependência direta no SDK aqui
	shutdown  func(context.Context) error
	client    *OTLPClient
}

// NewOTLPTracingExporter cria um adapter configurado.
// NewOTLPTracingExporter cria um adapter que, por conveniência, instancia
// um `OTLPClient` internamente. Para compartilhar o mesmo client entre
// múltiplos exporters, use `NewOTLPTracingExporterWithClient`.
// NewOTLPTracingExporter cria um adapter e conecta internamente ao endpoint OTLP.
func NewOTLPTracingExporter(ctx context.Context, endpoint string, setGlobal bool) (*OTLPTracingExporter, error) {
	client, err := NewOTLPClient(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	return &OTLPTracingExporter{endpoint: endpoint, setGlobal: setGlobal, client: client}, nil
}

// NewOTLPTracingExporterWithClient cria um adapter usando a instância de client fornecida.
func NewOTLPTracingExporterWithClient(client *OTLPClient, setGlobal bool) *OTLPTracingExporter {
	return &OTLPTracingExporter{endpoint: client.endpoint, setGlobal: setGlobal, client: client}
}

// Start inicializa o TracerProvider via providers.NewOTLPTracerProvider
func (o *OTLPTracingExporter) Start(ctx context.Context) error {
	// Cria o exporter OTLP via cliente local (usa a otlptrace.Client compartilhado)
	exp, err := otlptrace.New(ctx, otlptrace.WithClient(o.client.Client()))
	if err != nil {
		return fmt.Errorf("iniciar OTLP tracing exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(exp)
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(bsp))

	if o.setGlobal {
		otel.SetTracerProvider(tp)
	}

	o.tp = tp
	o.shutdown = func(ctx context.Context) error { return tp.Shutdown(ctx) }
	return nil
}

// Shutdown encerra o provider/flush
func (o *OTLPTracingExporter) Shutdown(ctx context.Context) error {
	if o.shutdown == nil {
		return nil
	}
	return o.shutdown(ctx)
}

// StartSpan cria um span usando o TracerProvider interno.
func (o *OTLPTracingExporter) StartSpan(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, func()) {
	if o.tp == nil {
		// fallback: usar tracer global
		tr := otel.Tracer("otlp-exporter/fallback")
		ctx2, span := tr.Start(ctx, name, trace.WithAttributes(attrs...))
		return ctx2, func() { span.End() }
	}

	// tp é *sdktrace.TracerProvider, mas mantemos interface{} para reduzir acoplamento
	if tp, ok := o.tp.(interface {
		Tracer(string, ...trace.TracerOption) trace.Tracer
	}); ok {
		tr := tp.Tracer("otlp-exporter")
		ctx2, span := tr.Start(ctx, name, trace.WithAttributes(attrs...))
		return ctx2, func() { span.End() }
	}

	tr := otel.Tracer("otlp-exporter/fallback")
	ctx2, span := tr.Start(ctx, name, trace.WithAttributes(attrs...))
	return ctx2, func() { span.End() }
}

// Ensure OTLPTracingExporter implements the contracts.TracingExporter interface
var _ contracts.TracingExporter = (*OTLPTracingExporter)(nil)
