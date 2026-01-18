package providers

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	adapters "github.com/brunojet/go-infra-backend/internal/observability/exporters/adapters"
	otlptrace "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

// NewOTLPTracerProvider cria um TracerProvider que exporta spans para um OTLP
// collector via gRPC. Retorna o provider, uma função de shutdown e um erro.
// Endpoint deve ser no formato "host:port" (ex.: "localhost:4317").
func NewOTLPTracerProvider(ctx context.Context, endpoint string, setGlobal bool) (*sdktrace.TracerProvider, func(context.Context) error, error) {
	if endpoint == "" {
		return nil, nil, fmt.Errorf("endpoint vazio")
	}

	// Cria exporter OTLP gRPC via cliente comum
	client, err := adapters.NewOTLPClient(ctx, endpoint)
	if err != nil {
		return nil, nil, fmt.Errorf("criar otlp client: %w", err)
	}

	// cria exporter usando o OTLP client
	exporter, err := otlptrace.New(ctx, otlptrace.WithClient(client.Client()))
	if err != nil {
		return nil, nil, fmt.Errorf("criar otlp exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(bsp))

	if setGlobal {
		otel.SetTracerProvider(tp)
	}

	shutdown := func(ctx context.Context) error {
		// Primeiro desligar o tracer provider
		if err := tp.Shutdown(ctx); err != nil {
			return fmt.Errorf("shutdown tracer provider: %w", err)
		}
		// Fechar conexão OTLP client se aplicável
		if client != nil {
			_ = client.Close(ctx)
		}
		return nil
	}

	return tp, shutdown, nil
}
