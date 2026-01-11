package ports

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeContextBase struct {
	status    int
	jsonOut   any
	stringOut string
	errMsg    string
	paramVal  string
	bindErr   error
	ctx       context.Context
}

func (f *fakeContextBase) RequestContext() context.Context {
	if f.ctx != nil {
		return f.ctx
	}
	return context.Background()
}
func (f *fakeContextBase) Param(name string) string    { return f.paramVal }
func (f *fakeContextBase) BindJSON(dst any) error      { return f.bindErr }
func (f *fakeContextBase) JSON(status int, v any)      { f.status = status; f.jsonOut = v }
func (f *fakeContextBase) String(status int, v string) { f.status = status; f.stringOut = v }
func (f *fakeContextBase) Status(status int)           { f.status = status }

func TestBaseHandler_JSON(t *testing.T) {
	h := BaseHandler{}
	ctx := &fakeContextBase{}
	h.JSON(ctx, 201, map[string]string{"foo": "bar"})
	assert.Equal(t, 201, ctx.status)
	assert.Equal(t, map[string]string{"foo": "bar"}, ctx.jsonOut)
}

func TestBaseHandler_Error(t *testing.T) {
	h := BaseHandler{}
	ctx := &fakeContextBase{}
	h.Error(ctx, 400, "msg")
	assert.Equal(t, 400, ctx.status)
	assert.Equal(t, ErrorResponse{Error: "msg"}, ctx.jsonOut)
}

func TestBaseHandler_BadRequest(t *testing.T) {
	h := BaseHandler{}
	ctx := &fakeContextBase{}
	h.BadRequest(ctx, "bad")
	assert.Equal(t, 400, ctx.status)
	assert.Equal(t, ErrorResponse{Error: "bad"}, ctx.jsonOut)
}

func TestBaseHandler_NotFound(t *testing.T) {
	h := BaseHandler{}
	ctx := &fakeContextBase{}
	h.NotFound(ctx, "nf")
	assert.Equal(t, 404, ctx.status)
	assert.Equal(t, ErrorResponse{Error: "nf"}, ctx.jsonOut)
}

func TestBaseHandler_InternalError(t *testing.T) {
	h := BaseHandler{}
	ctx := &fakeContextBase{}
	h.InternalError(ctx)
	assert.Equal(t, 500, ctx.status)
	assert.Equal(t, ErrorResponse{Error: "internal error"}, ctx.jsonOut)
}

func TestBaseHandler_BindJSON_Success(t *testing.T) {
	h := BaseHandler{}
	ctx := &fakeContextBase{}
	ok := h.BindJSON(ctx, &struct{}{})
	assert.True(t, ok)
}

func TestBaseHandler_BindJSON_Error(t *testing.T) {
	h := BaseHandler{}
	ctx := &fakeContextBase{bindErr: errors.New("fail")}
	ok := h.BindJSON(ctx, &struct{}{})
	assert.False(t, ok)
	assert.Equal(t, 400, ctx.status)
	assert.Equal(t, ErrorResponse{Error: "invalid json"}, ctx.jsonOut)
}

func TestBaseHandler_ParseUintParam_Success(t *testing.T) {
	h := BaseHandler{}
	ctx := &fakeContextBase{paramVal: "42"}
	v, ok := h.ParseUintParam(ctx, "id")
	assert.True(t, ok)
	assert.Equal(t, uint(42), v)
}

func TestBaseHandler_ParseUintParam_Error(t *testing.T) {
	h := BaseHandler{}
	ctx := &fakeContextBase{paramVal: "notanumber"}
	v, ok := h.ParseUintParam(ctx, "id")
	assert.False(t, ok)
	assert.Equal(t, uint(0), v)
	assert.Equal(t, 400, ctx.status)
	assert.Equal(t, ErrorResponse{Error: "invalid id"}, ctx.jsonOut)
}
