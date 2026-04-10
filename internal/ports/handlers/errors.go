package handlers

import "errors"

var (
	ErrNotSupported    = errors.New("operation not supported")
	ErrInvalidIDFormat = errors.New("invalid ID format")
	ErrInvalidJSONBody = errors.New("invalid JSON body")
)
