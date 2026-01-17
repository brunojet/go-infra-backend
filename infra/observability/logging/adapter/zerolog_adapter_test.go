package adapter_test

import (
	"context"
	"io"
	"testing"

	ada "github.com/brunojet/go-infra-backend/infra/observability/logging/adapter"
	"github.com/rs/zerolog"
)

func TestZerologAdapter_NoPanicAndDerived(t *testing.T) {
	z := zerolog.New(io.Discard).With().Timestamp().Logger()
	ad := ada.NewZerologAdapter(z)
	// basic calls should not panic
	ad.Info(context.Background(), "info", nil)
	ad.Debug(context.Background(), "debug", map[string]interface{}{"k": "v"})
	ad.Warn(context.Background(), "warn", nil)
	ad.Error(context.Background(), "err", nil, map[string]interface{}{"e": 1})

	// derived should also work
	d := ad.WithFields(map[string]interface{}{"a": "1"}).WithContext(context.WithValue(context.Background(), "trace", "t1"))
	d.Info(context.Background(), "msg", map[string]interface{}{"b": "2"})
}
