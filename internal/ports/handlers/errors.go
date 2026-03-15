package handlers

import "errors"

var (
	ErrNotSupported    = errors.New("operation not supported")
	ErrInvalidIDFormat = errors.New("invalid ID format")
	ErrInvalidParentID = errors.New("invalid parentID")
	ErrInvalidJSONBody = errors.New("invalid JSON body")
)
