package adapters

import (
	"context"
	"os"
	"time"

	"github.com/brunojet/go-infra-backend/internal/observability/exporters/contracts"
	"github.com/rs/zerolog"

	"google.golang.org/grpc"

	collectorlogs "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	commonv1 "go.opentelemetry.io/proto/otlp/common/v1"
	logsproto "go.opentelemetry.io/proto/otlp/logs/v1"
	resourcev1 "go.opentelemetry.io/proto/otlp/resource/v1"
)

// OTLPLoggingExporter implements contracts.LoggingExporter. Currently it uses zerolog
// to write logs locally. OTLP logging exporters are not yet wired here.
type OTLPLoggingExporter struct {
	logger  zerolog.Logger
	conn    *grpc.ClientConn
	client  collectorlogs.LogsServiceClient
	svcName string
}

// NewOTLPLoggingExporter creates a new logging exporter (writes to stdout by default).
func NewOTLPLoggingExporter() contracts.LoggingExporter {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &OTLPLoggingExporter{logger: logger, svcName: "go-infra-backend"}
}

// Start initializes the logging exporter.
func (o *OTLPLoggingExporter) Start(ctx context.Context) error {
	endpoint := EndpointFromEnv()
	cc, err := DialOTLP(ctx, endpoint)
	if err != nil {
		// keep zerolog as fallback
		return nil
	}
	o.conn = cc
	o.client = collectorlogs.NewLogsServiceClient(cc)
	return nil
}

// Shutdown shuts down the logging exporter.
func (o *OTLPLoggingExporter) Shutdown(ctx context.Context) error {
	if o.conn != nil {
		return o.conn.Close()
	}
	return nil
}

// Emit writes a log with given severity and body.
func (o *OTLPLoggingExporter) Emit(ctx context.Context, severity string, body string) error {
	// If OTLP client available, send Export request; otherwise fallback to zerolog.
	if o.client == nil {
		switch severity {
		case "debug":
			o.logger.Debug().Msg(body)
		case "info":
			o.logger.Info().Msg(body)
		case "warn", "warning":
			o.logger.Warn().Msg(body)
		case "error":
			o.logger.Error().Msg(body)
		default:
			o.logger.Info().Msg(body)
		}
		return nil
	}

	now := uint64(time.Now().UnixNano())

	lr := &logsproto.LogRecord{
		TimeUnixNano: now,
		SeverityText: severity,
		Body:         &commonv1.AnyValue{Value: &commonv1.AnyValue_StringValue{StringValue: body}},
	}

	// Build ScopeLogs (newer OTLP naming) and ResourceLogs
	scopeLogs := &logsproto.ScopeLogs{
		Scope:      &commonv1.InstrumentationScope{Name: o.svcName},
		LogRecords: []*logsproto.LogRecord{lr},
	}

	resLogs := &logsproto.ResourceLogs{
		Resource:  &resourcev1.Resource{Attributes: []*commonv1.KeyValue{{Key: "service.name", Value: &commonv1.AnyValue{Value: &commonv1.AnyValue_StringValue{StringValue: o.svcName}}}}},
		ScopeLogs: []*logsproto.ScopeLogs{scopeLogs},
	}

	req := &collectorlogs.ExportLogsServiceRequest{ResourceLogs: []*logsproto.ResourceLogs{resLogs}}

	_, err := o.client.Export(ctx, req)
	if err != nil {
		// if export fails, fallback to local logger
		o.logger.Info().Msgf("otlp export failed: %v; original: %s", err, body)
		return err
	}
	return nil
}

var _ contracts.LoggingExporter = (*OTLPLoggingExporter)(nil)
