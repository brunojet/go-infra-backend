package ports

import (
	"context"
	"errors"
	"net/http"

	"github.com/brunojet/go-infra-backend/internal/http/contracts"
)

type CRUDService[T any, CreateReq any, PatchReq any] interface {
	List(ctx context.Context) ([]T, error)
	Get(ctx context.Context, id uint) (*T, error)
	Create(ctx context.Context, req CreateReq) (*T, error)
	Patch(ctx context.Context, id uint, req PatchReq) (*T, error)
	Delete(ctx context.Context, id uint) error
}

type ToResponse[T any, Resp any] func(T) Resp

// CRUDHandler implements common HTTP plumbing for CRUD endpoints.
// It is adapter-level code (framework-agnostic) and should be reusable across resources.
//
// Error mapping:
// - if errors.Is(err, ValidationError) => 400
// - if errors.Is(err, NotFoundError) => 404
// - otherwise => 500
//
// Delete is assumed idempotent at the service level; HTTP returns 204 on success.
//
// This handler does not depend on Gin; the framework-specific adapter must provide a contracts.Context.
type CRUDHandler[T any, CreateReq any, PatchReq any, Resp any] struct {
	BaseHandler
	Service         CRUDService[T, CreateReq, PatchReq]
	ToResponse      ToResponse[T, Resp]
	NotFoundError   error
	ValidationError error
}

func (h *CRUDHandler[T, CreateReq, PatchReq, Resp]) List(c contracts.Context) {
	items, err := h.Service.List(c.RequestContext())
	if err != nil {
		h.writeError(c, err)
		return
	}

	resp := make([]Resp, 0, len(items))
	for _, it := range items {
		resp = append(resp, h.ToResponse(it))
	}
	h.JSON(c, http.StatusOK, resp)
}

func (h *CRUDHandler[T, CreateReq, PatchReq, Resp]) Get(c contracts.Context) {
	id, ok := h.ParseUintParam(c, "id")
	if !ok {
		return
	}

	item, err := h.Service.Get(c.RequestContext(), id)
	if err != nil {
		h.writeError(c, err)
		return
	}

	h.JSON(c, http.StatusOK, h.ToResponse(*item))
}

func (h *CRUDHandler[T, CreateReq, PatchReq, Resp]) Create(c contracts.Context) {
	var req CreateReq
	if !h.BindJSON(c, &req) {
		return
	}

	item, err := h.Service.Create(c.RequestContext(), req)
	if err != nil {
		h.writeError(c, err)
		return
	}

	h.JSON(c, http.StatusCreated, h.ToResponse(*item))
}

func (h *CRUDHandler[T, CreateReq, PatchReq, Resp]) Patch(c contracts.Context) {
	id, ok := h.ParseUintParam(c, "id")
	if !ok {
		return
	}

	var req PatchReq
	if !h.BindJSON(c, &req) {
		return
	}

	item, err := h.Service.Patch(c.RequestContext(), id, req)
	if err != nil {
		h.writeError(c, err)
		return
	}

	h.JSON(c, http.StatusOK, h.ToResponse(*item))
}

func (h *CRUDHandler[T, CreateReq, PatchReq, Resp]) Delete(c contracts.Context) {
	id, ok := h.ParseUintParam(c, "id")
	if !ok {
		return
	}

	if err := h.Service.Delete(c.RequestContext(), id); err != nil {
		h.writeError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *CRUDHandler[T, CreateReq, PatchReq, Resp]) writeError(c contracts.Context, err error) {
	if h.ValidationError != nil && errors.Is(err, h.ValidationError) {
		h.Error(c, http.StatusBadRequest, "validation")
		return
	}
	if h.NotFoundError != nil && errors.Is(err, h.NotFoundError) {
		h.NotFound(c, "not found")
		return
	}
	h.InternalError(c)
}
