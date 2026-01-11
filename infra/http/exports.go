package http

import (
	"github.com/brunojet/go-infra-backend/internal/http"
	"github.com/brunojet/go-infra-backend/internal/http/contracts"
	"github.com/brunojet/go-infra-backend/internal/http/types"
)

type CORSConfig = types.CORSConfig

type Context = contracts.Context

type HandlerFunc = contracts.HandlerFunc

type Router = contracts.Router

type HTTPDriver = types.HTTPDriver

type HTTPParams = http.HTTPParams

const (
	HTTPDriverGin = types.HTTPDriverGin
	HTTPDriverChi = types.HTTPDriverChi
)

var NormalizeHTTPDriver = types.NormalizeHTTPDriver
var NewHTTPManager = http.NewHTTPManager
