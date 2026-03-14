package services

import "errors"

var (
	ErrScopeKeyNotFound  = errors.New("scope key not found or invalid")
	ErrScopeValueInvalid = errors.New("scope value invalid or below minimum")
)
