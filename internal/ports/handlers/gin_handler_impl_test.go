package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	repoerrs "github.com/brunojet/go-infra-backend/internal/ports/repositories"
	hndcontracts "github.com/brunojet/go-infra-backend/pkg/ports/handlers/contracts"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type SimpleEntity struct {
	ID   string
	Name string
}

func (e SimpleEntity) TableName() string {
	return "simple_entities"
}

type SimpleDTO struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

// helper to create via GinHandler
func createEntityForHandlerTest(t *testing.T, h hndcontracts.GenericHandler[SimpleDTO, SimpleDTO, SimpleDTO, SimpleEntity], ms *MockService[SimpleDTO, SimpleDTO, SimpleDTO, SimpleEntity]) string {
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	body := `{"name":"bob"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	ms.EXPECT().Create(gomock.Any(), gomock.AssignableToTypeOf(SimpleDTO{})).DoAndReturn(func(ctx context.Context, dto any) (SimpleDTO, error) {
		d := dto.(SimpleDTO)
		d.ID = "created-id"
		return d, nil
	})

	h.Create(c)
	require.Equal(t, http.StatusCreated, rec.Code)
	var dto SimpleDTO
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &dto))
	return dto.ID
}

func TestCreate_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := NewMockService[SimpleDTO, SimpleDTO, SimpleDTO, SimpleEntity](ctrl)
	hp := HandlerParameters{IDValidationRule: Int64GtZero}
	h := NewGenericHandler[SimpleDTO, SimpleDTO, SimpleDTO, SimpleEntity](hp, ms)

	id := createEntityForHandlerTest(t, h, ms)
	require.Equal(t, "created-id", id)

	// error case
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	body := `{"name":"bob"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")
	ms.EXPECT().Create(gomock.Any(), gomock.AssignableToTypeOf(SimpleDTO{})).Return(SimpleDTO{}, errors.New("boom"))
	h.Create(c)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetByID_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := NewMockService[SimpleDTO, SimpleDTO, SimpleDTO, SimpleEntity](ctrl)
	hp := HandlerParameters{IDValidationRule: Int64GtZero}
	h := NewGenericHandler[SimpleDTO, SimpleDTO, SimpleDTO, SimpleEntity](hp, ms)

	// success
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	ms.EXPECT().GetByID(gomock.Any(), "1").Return(SimpleDTO{ID: "1"}, nil)
	h.GetByID(c)
	require.Equal(t, http.StatusOK, rec.Code)

	// not found
	rec = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/2", nil)
	c.Params = gin.Params{{Key: "id", Value: "2"}}
	ms.EXPECT().GetByID(gomock.Any(), "2").Return(SimpleDTO{}, repoerrs.ErrNotFound)
	h.GetByID(c)
	require.Equal(t, http.StatusNotFound, rec.Code)

	// id inválido
	rec = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/0", nil)
	c.Params = gin.Params{{Key: "id", Value: "0"}} // inválido para Int64GtZero
	h.GetByID(c)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestList_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := NewMockService[SimpleDTO, SimpleDTO, SimpleDTO, SimpleEntity](ctrl)
	hp := HandlerParameters{IDValidationRule: Int64GtZero}
	h := NewGenericHandler[SimpleDTO, SimpleDTO, SimpleDTO, SimpleEntity](hp, ms)

	// success
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ms.EXPECT().List(gomock.Any(), gomock.Any()).Return([]SimpleDTO{{ID: "1"}}, int64(1), nil)
	h.List(c)
	require.Equal(t, http.StatusOK, rec.Code)

	// db unavailable
	rec = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ms.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, int64(0), repoerrs.ErrDBUnavailable)
	h.List(c)
	require.Equal(t, http.StatusServiceUnavailable, rec.Code)
}

func TestUpdate_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := NewMockService[SimpleDTO, SimpleDTO, SimpleDTO, SimpleEntity](ctrl)
	hp := HandlerParameters{IDValidationRule: Int64GtZero}
	h := NewGenericHandler[SimpleDTO, SimpleDTO, SimpleDTO, SimpleEntity](hp, ms)

	// success
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	upd := `{"name":"joe"}`
	c.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(upd))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "1"}} // id válido para Int64GtZero
	ms.EXPECT().Update(gomock.Any(), "1", gomock.AssignableToTypeOf(SimpleDTO{})).Return(SimpleDTO{}, nil)
	h.Update(c)
	require.Equal(t, http.StatusOK, rec.Code)

	// not found
	rec = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(upd))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "2"}} // outro id válido
	ms.EXPECT().Update(gomock.Any(), "2", gomock.AssignableToTypeOf(SimpleDTO{})).Return(SimpleDTO{}, repoerrs.ErrNotFound)
	h.Update(c)
	require.Equal(t, http.StatusNotFound, rec.Code)

	// id inválido
	rec = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(upd))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "0"}} // inválido
	h.Update(c)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDelete_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := NewMockService[SimpleDTO, SimpleDTO, SimpleDTO, SimpleEntity](ctrl)
	hp := HandlerParameters{IDValidationRule: Int64GtZero}
	h := NewGenericHandler[SimpleDTO, SimpleDTO, SimpleDTO, SimpleEntity](hp, ms)

	// success
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodDelete, "/", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	ms.EXPECT().Delete(gomock.Any(), "1").Return(nil)
	h.Delete(c)
	require.Equal(t, http.StatusNoContent, rec.Code)

	// error
	rec = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodDelete, "/", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	ms.EXPECT().Delete(gomock.Any(), "1").Return(errors.New("boom"))
	h.Delete(c)
	require.Equal(t, http.StatusInternalServerError, rec.Code)

	// id inválido
	rec = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodDelete, "/", nil)
	c.Params = gin.Params{{Key: "id", Value: "0"}} // inválido
	h.Delete(c)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRegister_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := NewMockService[SimpleDTO, SimpleDTO, SimpleDTO, SimpleEntity](ctrl)
	hp := HandlerParameters{HandlerPath: "/ping"}
	h := NewGenericHandler[SimpleDTO, SimpleDTO, SimpleDTO, SimpleEntity](hp, ms)

	// create real gin engine and group
	engine := gin.New()
	rg := engine.Group("/api")

	// register a simple handler using lowercase method to exercise ToUpper
	h.RegisterCollection(rg, "get", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"pong": true})
	})

	// perform request
	req := httptest.NewRequest(http.MethodGet, "/api/ping", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	var body map[string]bool
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.Equal(t, true, body["pong"])
}

func TestCreate_Handler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := NewMockService[SimpleDTO, SimpleDTO, SimpleDTO, SimpleEntity](ctrl)
	hp := HandlerParameters{IDValidationRule: Int64GtZero}
	h := NewGenericHandler[SimpleDTO, SimpleDTO, SimpleDTO, SimpleEntity](hp, ms)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("not-json"))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Create(c)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdate_Handler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := NewMockService[SimpleDTO, SimpleDTO, SimpleDTO, SimpleEntity](ctrl)
	hp := HandlerParameters{IDValidationRule: Int64GtZero}
	h := NewGenericHandler[SimpleDTO, SimpleDTO, SimpleDTO, SimpleEntity](hp, ms)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString("not-json"))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "1"}} // id válido para Int64GtZero

	h.Update(c)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}
