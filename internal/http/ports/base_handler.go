package ports

import (
	"net/http"
	"strconv"

	"github.com/brunojet/go-infra-backend/internal/http/contracts"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type BaseHandler struct{}

func (h BaseHandler) JSON(c contracts.Context, status int, v any) {
	c.JSON(status, v)
}

func (h BaseHandler) Error(c contracts.Context, status int, message string) {
	h.JSON(c, status, ErrorResponse{Error: message})
}

func (h BaseHandler) BadRequest(c contracts.Context, message string) {
	h.Error(c, http.StatusBadRequest, message)
}

func (h BaseHandler) NotFound(c contracts.Context, message string) {
	h.Error(c, http.StatusNotFound, message)
}

func (h BaseHandler) InternalError(c contracts.Context) {
	h.Error(c, http.StatusInternalServerError, "internal error")
}

func (h BaseHandler) BindJSON(c contracts.Context, dst any) bool {
	if err := c.BindJSON(dst); err != nil {
		h.BadRequest(c, "invalid json")
		return false
	}
	return true
}

func (h BaseHandler) ParseUintParam(c contracts.Context, name string) (uint, bool) {
	raw := c.Param(name)
	v, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		h.BadRequest(c, "invalid "+name)
		return 0, false
	}
	return uint(v), true
}
