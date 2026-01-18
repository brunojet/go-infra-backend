package database

import (
	"context"
	"encoding/json"
	"testing"

	plugins "github.com/brunojet/go-infra-backend/internal/database/plugins"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

type testEntity struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex"`
}

// TestCreateRead_ExportsTrace creates and reads a GORM entity while the
// OpenTelemetry tracing plugin is registered and asserts that spans were
// exported to the in-memory exporter.
func TestCreateRead_ExportsTrace(t *testing.T) {
	// Setup in-memory exporter and tracer provider
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sdktrace.NewSimpleSpanProcessor(exporter)))
	otel.SetTracerProvider(tp)
	defer func() {
		_ = tp.Shutdown(context.Background())
	}()

	// Create the OTEL plugin and DB
	p, err := plugins.NewOtelGormPlugin()
	if err != nil {
		t.Fatalf("failed to create otel plugin: %v", err)
	}

	db, err := NewSQLiteDatabase(p)
	if err != nil {
		t.Fatalf("NewSQLiteDatabase returned error: %v", err)
	}
	defer func() {
		_ = db.Close()
	}()

	// Migrate schema
	if err := db.GormDB().AutoMigrate(&testEntity{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}

	// Create
	e := testEntity{Name: "trace-test"}
	if err := db.GormDB().Create(&e).Error; err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Read
	var got testEntity
	if err := db.GormDB().First(&got, "name = ?", "trace-test").Error; err != nil {
		t.Fatalf("read failed: %v", err)
	}

	// Ensure spans were exported. Flush provider then inspect exporter.
	if err := tp.ForceFlush(context.Background()); err != nil {
		t.Fatalf("force flush failed: %v", err)
	}

	spans := exporter.GetSpans()
	if len(spans) == 0 {
		t.Fatalf("expected exported spans, got 0")
	}

	// Basic assertion: at least one span exists. Convert spans to JSON and log.
	t.Logf("exported %d spans", len(spans))

	type spanOut struct {
		Name       string                 `json:"name"`
		TraceID    string                 `json:"traceID"`
		SpanID     string                 `json:"spanID"`
		Attributes map[string]interface{} `json:"attributes"`
	}

	var out []spanOut
	for _, s := range spans {
		attrs := make(map[string]interface{})
		for _, kv := range s.Attributes {
			attrs[string(kv.Key)] = kv.Value.AsInterface()
		}
		out = append(out, spanOut{
			Name:       s.Name,
			TraceID:    s.SpanContext.TraceID().String(),
			SpanID:     s.SpanContext.SpanID().String(),
			Attributes: attrs,
		})
	}

	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		t.Logf("failed to marshal spans to json: %v", err)
	} else {
		t.Logf("exported spans json:\n%s", string(b))
	}
}
