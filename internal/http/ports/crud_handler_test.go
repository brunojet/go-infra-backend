package ports

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeContext struct {
	params  map[string]string
	jsonIn  any
	jsonOut any
	status  int
	ctx     context.Context
}

func (f *fakeContext) RequestContext() context.Context {
	if f.ctx != nil {
		return f.ctx
	}
	return context.Background()
}
func (f *fakeContext) Param(name string) string { return f.params[name] }
func (f *fakeContext) BindJSON(dst any) error {
	if f.jsonIn == nil {
		return errors.New("bind error") // simula erro real de bind
	}
	switch d := dst.(type) {
	case *string:
		*d = f.jsonIn.(string)
	}
	return nil
}
func (f *fakeContext) JSON(status int, v any)      { f.status = status; f.jsonOut = v }
func (f *fakeContext) String(status int, v string) { f.status = status }
func (f *fakeContext) Status(status int)           { f.status = status }

type fakeService struct {
	list   func(ctx context.Context) ([]string, error)
	get    func(ctx context.Context, id uint) (*string, error)
	create func(ctx context.Context, req string) (*string, error)
	patch  func(ctx context.Context, id uint, req string) (*string, error)
	delete func(ctx context.Context, id uint) error
}

func (f *fakeService) List(ctx context.Context) ([]string, error)        { return f.list(ctx) }
func (f *fakeService) Get(ctx context.Context, id uint) (*string, error) { return f.get(ctx, id) }
func (f *fakeService) Create(ctx context.Context, req string) (*string, error) {
	return f.create(ctx, req)
}
func (f *fakeService) Patch(ctx context.Context, id uint, req string) (*string, error) {
	return f.patch(ctx, id, req)
}
func (f *fakeService) Delete(ctx context.Context, id uint) error { return f.delete(ctx, id) }

func toResp(s string) string { return s + "-resp" }

func TestCRUDHandler_List_Integration(t *testing.T) {
	h := &CRUDHandler[string, string, string, string]{
		Service: &fakeService{
			list: func(ctx context.Context) ([]string, error) { return []string{"a", "b"}, nil },
		},
		ToResponse: toResp,
	}
	ctx := &fakeContext{ctx: context.Background()}
	h.List(ctx)
	assert.Equal(t, []string{"a-resp", "b-resp"}, ctx.jsonOut)
}

func TestCRUDHandler_Get_Integration(t *testing.T) {
	h := &CRUDHandler[string, string, string, string]{
		Service: &fakeService{
			get: func(ctx context.Context, id uint) (*string, error) { v := "foo"; return &v, nil },
		},
		ToResponse: toResp,
	}
	ctx := &fakeContext{ctx: context.Background(), params: map[string]string{"id": "1"}}
	h.Get(ctx)
	assert.Equal(t, "foo-resp", ctx.jsonOut)
}

func TestCRUDHandler_Create_Integration(t *testing.T) {
	h := &CRUDHandler[string, string, string, string]{
		Service: &fakeService{
			create: func(ctx context.Context, req string) (*string, error) { v := req + "-created"; return &v, nil },
		},
		ToResponse: toResp,
	}
	ctx := &fakeContext{ctx: context.Background(), jsonIn: "foo"}
	h.Create(ctx)
	assert.Equal(t, "foo-created-resp", ctx.jsonOut)
}

func TestCRUDHandler_Patch_Integration(t *testing.T) {
	h := &CRUDHandler[string, string, string, string]{
		Service: &fakeService{
			patch: func(ctx context.Context, id uint, req string) (*string, error) { v := req + "-patched"; return &v, nil },
		},
		ToResponse: toResp,
	}
	ctx := &fakeContext{ctx: context.Background(), params: map[string]string{"id": "2"}, jsonIn: "bar"}
	h.Patch(ctx)
	assert.Equal(t, "bar-patched-resp", ctx.jsonOut)
}

func TestCRUDHandler_Delete_Integration(t *testing.T) {
	h := &CRUDHandler[string, string, string, string]{
		Service: &fakeService{
			delete: func(ctx context.Context, id uint) error { return nil },
		},
		ToResponse: toResp,
	}
	ctx := &fakeContext{ctx: context.Background(), params: map[string]string{"id": "3"}}
	h.Delete(ctx)
	assert.Equal(t, 204, ctx.status)
}

func TestCRUDHandler_List_Error(t *testing.T) {
	h := &CRUDHandler[string, string, string, string]{
		Service: &fakeService{
			list: func(ctx context.Context) ([]string, error) { return nil, errors.New("fail") },
		},
		ToResponse: toResp,
	}
	ctx := &fakeContext{ctx: context.Background()}
	h.List(ctx)
	assert.NotEqual(t, 200, ctx.status)
}

func TestCRUDHandler_Get_InvalidParam(t *testing.T) {
	h := &CRUDHandler[string, string, string, string]{
		Service:    &fakeService{},
		ToResponse: toResp,
	}
	ctx := &fakeContext{ctx: context.Background(), params: map[string]string{"id": "notanumber"}}
	h.Get(ctx)
	assert.NotEqual(t, 200, ctx.status)
}

func TestCRUDHandler_Get_Error(t *testing.T) {
	h := &CRUDHandler[string, string, string, string]{
		Service: &fakeService{
			get: func(ctx context.Context, id uint) (*string, error) { return nil, errors.New("fail") },
		},
		ToResponse: toResp,
	}
	ctx := &fakeContext{ctx: context.Background(), params: map[string]string{"id": "1"}}
	h.Get(ctx)
	assert.NotEqual(t, 200, ctx.status)
}

func TestCRUDHandler_Create_BindError(t *testing.T) {
	h := &CRUDHandler[string, string, string, string]{
		Service: &fakeService{
			create: func(ctx context.Context, req string) (*string, error) {
				t.Fatalf("Create should not be called when BindJSON falha")
				return nil, nil
			},
		},
		ToResponse: toResp,
	}
	ctx := &fakeContext{ctx: context.Background(), jsonIn: nil}
	h.Create(ctx)
	assert.NotEqual(t, 201, ctx.status)
}

func TestCRUDHandler_Create_Error(t *testing.T) {
	h := &CRUDHandler[string, string, string, string]{
		Service: &fakeService{
			create: func(ctx context.Context, req string) (*string, error) { return nil, errors.New("fail") },
		},
		ToResponse: toResp,
	}
	ctx := &fakeContext{ctx: context.Background(), jsonIn: "foo"}
	h.Create(ctx)
	assert.NotEqual(t, 201, ctx.status)
}

func TestCRUDHandler_Patch_InvalidParam(t *testing.T) {
	h := &CRUDHandler[string, string, string, string]{
		Service:    &fakeService{},
		ToResponse: toResp,
	}
	ctx := &fakeContext{ctx: context.Background(), params: map[string]string{"id": "notanumber"}, jsonIn: "bar"}
	h.Patch(ctx)
	assert.NotEqual(t, 200, ctx.status)
}

func TestCRUDHandler_Patch_BindError(t *testing.T) {
	h := &CRUDHandler[string, string, string, string]{
		Service: &fakeService{
			patch: func(ctx context.Context, id uint, req string) (*string, error) {
				t.Errorf("Patch should not be called when BindJSON falha")
				return nil, nil
			},
		},
		ToResponse: toResp,
	}
	ctx := &fakeContext{ctx: context.Background(), params: map[string]string{"id": "2"}, jsonIn: nil}
	h.Patch(ctx)
	assert.NotEqual(t, 200, ctx.status)
}

func TestCRUDHandler_Patch_Error(t *testing.T) {
	h := &CRUDHandler[string, string, string, string]{
		Service: &fakeService{
			patch: func(ctx context.Context, id uint, req string) (*string, error) { return nil, errors.New("fail") },
		},
		ToResponse: toResp,
	}
	ctx := &fakeContext{ctx: context.Background(), params: map[string]string{"id": "2"}, jsonIn: "bar"}
	h.Patch(ctx)
	assert.NotEqual(t, 200, ctx.status)
}

func TestCRUDHandler_Delete_InvalidParam(t *testing.T) {
	h := &CRUDHandler[string, string, string, string]{
		Service:    &fakeService{},
		ToResponse: toResp,
	}
	ctx := &fakeContext{ctx: context.Background(), params: map[string]string{"id": "notanumber"}}
	h.Delete(ctx)
	assert.NotEqual(t, 204, ctx.status)
}

func TestCRUDHandler_Delete_Error(t *testing.T) {
	h := &CRUDHandler[string, string, string, string]{
		Service: &fakeService{
			delete: func(ctx context.Context, id uint) error { return errors.New("fail") },
		},
		ToResponse: toResp,
	}
	ctx := &fakeContext{ctx: context.Background(), params: map[string]string{"id": "3"}}
	h.Delete(ctx)
	assert.NotEqual(t, 204, ctx.status)
}

func TestCRUDHandler_writeError_Validation(t *testing.T) {
	h := &CRUDHandler[string, string, string, string]{
		ValidationError: errors.New("validation error"),
	}
	ctx := &fakeContext{ctx: context.Background()}
	h.writeError(ctx, h.ValidationError)
	assert.NotEqual(t, 500, ctx.status)
}

func TestCRUDHandler_writeError_NotFound(t *testing.T) {
	h := &CRUDHandler[string, string, string, string]{
		NotFoundError: errors.New("not found error"),
	}
	ctx := &fakeContext{ctx: context.Background()}
	h.writeError(ctx, h.NotFoundError)
	assert.NotEqual(t, 500, ctx.status)
}

func TestCRUDHandler_writeError_Internal(t *testing.T) {
	h := &CRUDHandler[string, string, string, string]{}
	ctx := &fakeContext{ctx: context.Background()}
	h.writeError(ctx, errors.New("other error"))
	assert.Equal(t, 500, ctx.status)
}
