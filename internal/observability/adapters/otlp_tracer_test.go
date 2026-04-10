package adapters

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewOTLPTracer_ExporterError(t *testing.T) {
	orig := otlptracegrpcNew
	// create a replacement function with the exact same signature using reflection
	tpe := reflect.TypeOf(otlptracegrpcNew)
	fn := reflect.MakeFunc(tpe, func(args []reflect.Value) []reflect.Value {
		numOut := tpe.NumOut()
		outs := make([]reflect.Value, numOut)
		for i := 0; i < numOut-1; i++ {
			outs[i] = reflect.Zero(tpe.Out(i))
		}
		outs[numOut-1] = reflect.ValueOf(errors.New("exporter failed"))
		return outs
	})
	reflect.ValueOf(&otlptracegrpcNew).Elem().Set(fn)
	defer reflect.ValueOf(&otlptracegrpcNew).Elem().Set(reflect.ValueOf(orig))

	cfg := OTLPConfig{Endpoint: "localhost:4317", IsInsecure: true}
	s, err := NewOTLPTracer(context.Background(), cfg)
	require.Error(t, err)
	require.Nil(t, s)
}
