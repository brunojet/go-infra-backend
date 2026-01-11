package ginadapter

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	httptypes "github.com/brunojet/go-infra-backend/internal/http/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORS_AllHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := httptypes.CORSConfig{
		AllowOrigins:     []string{"http://test"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"X-Test"},
		ExposeHeaders:    []string{"X-Expose"},
		AllowCredentials: true,
		MaxAge:           10 * time.Second,
	}
	r := gin.New()
	r.Use(CORS(cfg))
	r.GET("/", func(c *gin.Context) { c.String(200, "ok") })
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "http://test")
	r.ServeHTTP(w, req)

	assert.Equal(t, "http://test", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "X-Test", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "X-Expose", w.Header().Get("Access-Control-Expose-Headers"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	assert.NotEmpty(t, w.Header().Get("Access-Control-Max-Age"))
}

func TestCORS_OptionsPreflight(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := httptypes.CORSConfig{AllowOrigins: []string{"http://test"}}
	r := gin.New()
	r.Use(CORS(cfg))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/", nil)
	req.Header.Set("Origin", "http://test")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}
