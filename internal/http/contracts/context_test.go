package contracts

import (
	"context"
	"testing"
)

type dummyContext struct{}

func (dummyContext) RequestContext() context.Context { return context.Background() }
func (dummyContext) Param(name string) string        { return "" }
func (dummyContext) BindJSON(dst any) error          { return nil }
func (dummyContext) JSON(status int, v any)          {}
func (dummyContext) String(status int, v string)     {}
func (dummyContext) Status(status int)               {}

func TestContextInterface(t *testing.T) {
	var _ Context = dummyContext{}
}

func TestHandlerFuncType(t *testing.T) {
	var h HandlerFunc = func(ctx Context) {}
	h(dummyContext{})
}
