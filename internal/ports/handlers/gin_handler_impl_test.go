package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	rpo "github.com/brunojet/go-infra-backend/internal/ports/backend/repositories/errors"
	svccts "github.com/brunojet/go-infra-backend/pkg/ports/backend/services/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/handlers/contracts"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type SimpleEntity struct {
	ID   string
	Name string
}

func (e SimpleEntity) TableName() string {
	return "simple_entities"
}
func (e SimpleEntity) WhereOnConflict(tx *gorm.DB) *gorm.DB {
	return tx
}

type SimpleDTO struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

// helper to create via GinHandler
func createEntityForHandlerTest(t *testing.T, h contracts.GenericHandler[SimpleDTO, SimpleDTO, SimpleDTO], ms *MockService[SimpleDTO, SimpleDTO, SimpleDTO, SimpleEntity]) string {
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	body := `{"name":"bob"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	ExpectServiceCreate(ms, func(_ context.Context, dto SimpleDTO, out *SimpleDTO) error {
		d := dto
		d.ID = "created-id"
		if out != nil {
			*out = d
		}
		return nil
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
	h := NewGenericHandler(hp, ms)

	id := createEntityForHandlerTest(t, h, ms)
	require.Equal(t, "created-id", id)

	// error case
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	body := `{"name":"bob"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")
	ExpectServiceCreate(ms, func(_ context.Context, _ SimpleDTO, _ *SimpleDTO) error {
		return errors.New("boom")
	})
	h.Create(c)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetByID_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := NewMockService[SimpleDTO, SimpleDTO, SimpleDTO, SimpleEntity](ctrl)
	hp := HandlerParameters{IDValidationRule: Int64GtZero}
	h := NewGenericHandler(hp, ms)

	// success
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	ExpectServiceGetByID(ms, func(_ context.Context, id string, out *SimpleDTO) error {
		*out = SimpleDTO{ID: id}
		return nil
	})
	h.GetByID(c)
	require.Equal(t, http.StatusOK, rec.Code)

	// not found
	rec = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/2", nil)
	c.Params = gin.Params{{Key: "id", Value: "2"}}
	ExpectServiceGetByID(ms, func(_ context.Context, id string, _ *SimpleDTO) error {
		return rpo.ErrNotFound
	})
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
	h := NewGenericHandler(hp, ms)

	// success
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ExpectServiceList(ms, func(_ context.Context, _ svccts.ListParams, out *[]SimpleDTO) (int64, error) {
		if out != nil && len(*out) > 0 {
			(*out)[0] = SimpleDTO{ID: "1"}
		}
		return 1, nil
	})
	h.List(c)
	require.Equal(t, http.StatusPartialContent, rec.Code)

	// db unavailable
	rec = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ExpectServiceList(ms, func(_ context.Context, _ svccts.ListParams, _ *[]SimpleDTO) (int64, error) {
		return 0, rpo.ErrDBUnavailable
	})
	h.List(c)
	require.Equal(t, http.StatusServiceUnavailable, rec.Code)
}

func TestUpdate_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := NewMockService[SimpleDTO, SimpleDTO, SimpleDTO, SimpleEntity](ctrl)
	hp := HandlerParameters{IDValidationRule: Int64GtZero}
	h := NewGenericHandler(hp, ms)

	// success
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	upd := `{"name":"joe"}`
	c.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(upd))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "1"}} // id válido para Int64GtZero
	ExpectServiceUpdate(ms, func(_ context.Context, id string, _ SimpleDTO, _ *SimpleDTO) error {
		if id == "1" {
			return nil
		}
		return errors.New("unexpected id")
	})
	h.Update(c)
	require.Equal(t, http.StatusOK, rec.Code)

	// not found
	rec = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(upd))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "2"}} // outro id válido
	ExpectServiceUpdate(ms, func(_ context.Context, id string, _ SimpleDTO, _ *SimpleDTO) error {
		if id == "2" {
			return rpo.ErrNotFound
		}
		return nil
	})
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
	h := NewGenericHandler(hp, ms)

	// success
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodDelete, "/", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	ExpectServiceDelete(ms, func(_ context.Context, id string) error {
		if id == "1" {
			return nil
		}
		return errors.New("unexpected id")
	})
	h.Delete(c)
	require.Equal(t, http.StatusNoContent, rec.Code)

	// error
	rec = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodDelete, "/", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	ExpectServiceDelete(ms, func(_ context.Context, id string) error {
		return errors.New("boom")
	})
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
	h := NewGenericHandler(hp, ms)

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
	h := NewGenericHandler(hp, ms)

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
	h := NewGenericHandler(hp, ms)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString("not-json"))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "1"}} // id válido para Int64GtZero

	h.Update(c)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}
