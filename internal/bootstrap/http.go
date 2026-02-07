package bootstrap

import (
	"errors"
	"log"
	"net/http"

	"github.com/brunojet/go-infra-backend/internal/bootstrap/contracts"
	"github.com/brunojet/go-infra-backend/internal/config"
	middlewares "github.com/brunojet/go-infra-backend/internal/observability/http_middlewares"
	"github.com/gin-gonic/gin"
)

const (
	// HTTPAddrEnv is the env var key to override the HTTP listen address.
	httpAddrEnv = "HTTP_ADDR"
	// defaultHTTPAddr is the fallback address used when env var is not set.
	defaultHTTPAddr = ":8080"
)

type HttpServer struct {
	Router *gin.Engine
	srv    *http.Server
	sm     contracts.ShutdownManager
}

func NewHttpServerWithObservability(sm contracts.ShutdownManager) *HttpServer {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middlewares.OtelGinMiddleware())
	router.Use(middlewares.OTLPErrorLogMiddleware(middlewares.WithMinStatus(http.StatusBadRequest)))
	addr := config.Get(httpAddrEnv, defaultHTTPAddr)
	srv := http.Server{
		Addr:    addr,
		Handler: router,
	}
	return &HttpServer{
		Router: router,
		srv:    &srv,
		sm:     sm,
	}
}

func (h *HttpServer) StartAndWaitTermination() error {
	ctx := h.sm.GetContext()
	serverErrCh := make(chan error, 1)
	h.sm.Register("http-server", h.srv)

	go func() {
		log.Println("http server listening on : ", h.srv.Addr)
		if err := h.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
			return
		}
		serverErrCh <- nil
	}()

	select {
	case <-ctx.Done():
		log.Println("shutdown signal received")
	case err := <-serverErrCh:
		if err != nil {
			log.Printf("http server error: %v", err)
		}
	}
	return nil
}
