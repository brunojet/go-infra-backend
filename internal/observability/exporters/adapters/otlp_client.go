package adapters

import (
	"context"
	"fmt"

	otlptrace "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// OTLPClient é um cliente genérico OTLP que mantém uma instância de
// otlptrace.Client que pode ser reutilizada por exporters de tracing/metrics/logging.
type OTLPClient struct {
	endpoint string
	client   *otlptrace.Client
}

// NewOTLPClient cria um client OTLP configurado com o endpoint (ex.: "localhost:4317").
func NewOTLPClient(ctx context.Context, endpoint string) (*OTLPClient, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("endpoint vazio")
	}
	client := otlptrace.NewClient(otlptrace.WithEndpoint(endpoint), otlptrace.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))
	return &OTLPClient{endpoint: endpoint, client: client}, nil
}

// Client retorna o cliente OTLP subjacente.
func (c *OTLPClient) Client() *otlptrace.Client { return c.client }

// Close fecha recursos do client se necessário. Atualmente não é necessário.
func (c *OTLPClient) Close(ctx context.Context) error { return nil }
