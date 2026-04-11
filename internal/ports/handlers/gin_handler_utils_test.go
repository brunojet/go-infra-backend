package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	repoerrs "github.com/brunojet/go-infra-backend/internal/ports/backend/repositories"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMapErrorToStatus(t *testing.T) {
	a := assert.New(t)
	a.Equal(http.StatusNotFound, MapErrorToStatus(repoerrs.ErrNotFound))
	a.Equal(http.StatusServiceUnavailable, MapErrorToStatus(repoerrs.ErrDBUnavailable))
}

func TestSetResponseFromError_WritesJSON(t *testing.T) {
	a := assert.New(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	SetResponseFromError(c, repoerrs.ErrNotFound)
	a.Equal(http.StatusNotFound, w.Code)
	a.Contains(w.Body.String(), "error")
}

type sampleDTO struct {
	Name string `json:"name"`
}

func TestBindJSONToDTOPtr_SuccessAndFailure(t *testing.T) {
	a := assert.New(t)
	// success
	b := bytes.NewBufferString(`{"name":"x"}`)
	req := httptest.NewRequest(http.MethodPost, "/", b)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	var dto sampleDTO
	ok := BindJSONToDTO(c, &dto)
	a.EqualValues(ok, true)

	// failure
	rb := io.NopCloser(bytes.NewBufferString(`{`))
	badreq := httptest.NewRequest(http.MethodPost, "/", rb)
	badreq.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = badreq
	var dto2 sampleDTO
	ok2 := BindJSONToDTO(c2, &dto2)
	a.EqualValues(ok2, false)
	a.Equal(http.StatusBadRequest, w2.Code)
}

func TestBuildListParamsFromRequest_UsesCanonicalKeys(t *testing.T) {
	a := assert.New(t)

	req := httptest.NewRequest(http.MethodGet, "/?page=2&size=25&orderBy=created_at&order=desc&application_id=10", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	params := BuildListParamsFromRequest(c)
	a.Equal(2, params.Page)
	a.Equal(25, params.Size)
	a.Equal("created_at", params.OrderBy)
	a.Equal("desc", params.Order)
	a.Equal("10", params.QueryParams.Scopes["application_id"])
}

func TestBuildListParamsFromRequest_LeavesUnknownQueryAsScope(t *testing.T) {
	a := assert.New(t)

	req := httptest.NewRequest(http.MethodGet, "/?order_by=legacy&status=active", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	params := BuildListParamsFromRequest(c)
	a.Equal("", params.OrderBy)
	a.Equal("legacy", params.QueryParams.Scopes["order_by"])
	a.Equal("active", params.QueryParams.Scopes["status"])
}
