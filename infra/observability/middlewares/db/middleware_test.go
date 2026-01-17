package db_test

import (
	"net/http/httptest"
	"testing"

	dbpkg "github.com/brunojet/go-infra-backend/infra/observability/middlewares/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestDBMiddleware_InsertsDBIntoContext(t *testing.T) {
	// use a dummy *gorm.DB
	var db *gorm.DB = &gorm.DB{}
	adapter := dbpkg.NewSimpleDBAdapter(db)

	r := gin.New()
	r.Use(dbpkg.DBMiddleware(adapter))
	got := false
	r.GET("/test", func(c *gin.Context) {
		db2 := dbpkg.WithDBFromContext(c.Request.Context())
		if db2 != nil {
			got = true
		}
		c.Status(204)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if !got {
		t.Fatalf("expected DB to be present in context")
	}
}
