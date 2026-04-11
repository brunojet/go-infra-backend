package httpmiddlewares

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMiddlewares_ConstructorsNotNil(t *testing.T) {
	h1 := OtelGinMiddleware()
	require.NotNil(t, h1)

	h2 := OTLPErrorLogMiddleware()
	require.NotNil(t, h2)

	h3 := OTLPErrorLogMiddleware(WithMinStatus(http.StatusBadRequest))
	require.NotNil(t, h3)
}
