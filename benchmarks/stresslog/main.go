package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/brunojet/go-infra-backend/internal/observability"
	otellog "go.opentelemetry.io/otel/log"
	otellogglobal "go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

type noopLogExporter struct{}

func (n *noopLogExporter) Export(ctx context.Context, records []sdklog.Record) error { return nil }
func (n *noopLogExporter) Shutdown(ctx context.Context) error                        { return nil }
func (n *noopLogExporter) ForceFlush(ctx context.Context) error                      { return nil }

func main() {
	var (
		goroutines = flag.Int("goroutines", runtime.NumCPU(), "number of goroutines")
		perG       = flag.Int("per", 100_000, "messages per goroutine")
		msg        = flag.String("msg", "hello", "message payload")
	)
	flag.Parse()

	// Set up an in-process logger provider with a no-op exporter, so the stress test
	// mainly measures stdlib-log -> OTel pipeline overhead.
	prevProvider := otellogglobal.GetLoggerProvider()
	defer otellogglobal.SetLoggerProvider(prevProvider)

	exp := &noopLogExporter{}
	lp := sdklog.NewLoggerProvider(sdklog.WithProcessor(sdklog.NewSimpleProcessor(exp)))
	otellogglobal.SetLoggerProvider(lp)
	defer func() { _ = lp.Shutdown(context.Background()) }()

	restore := observability.RedirectStdLog("stdlib", otellog.SeverityInfo)
	defer restore()

	start := time.Now()
	var wg sync.WaitGroup
	for g := 0; g < *goroutines; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < *perG; i++ {
				log.Printf("g=%d i=%d %s", id, i, *msg)
			}
		}(g)
	}

	wg.Wait()
	elapsed := time.Since(start)

	total := int64(*goroutines) * int64(*perG)
	opsPerSec := float64(total) / elapsed.Seconds()

	fmt.Printf("goroutines=%d per=%d total=%d elapsed=%s ops/sec=%.0f\n", *goroutines, *perG, total, elapsed, opsPerSec)
}
