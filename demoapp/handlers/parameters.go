package handlers

import (
	hnd "github.com/brunojet/go-infra-backend/pkg/ports/handlers"
)

var (
	IDInt64Parameters = &hnd.HandlerParameters{
		IDValidationRule: hnd.Int64GtZero,
	}
	IDBase64Parameters = &hnd.HandlerParameters{
		IDValidationRule: hnd.Base64UrlSafe,
	}
)
