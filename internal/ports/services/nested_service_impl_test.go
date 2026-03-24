package services

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"testing"

	"github.com/brunojet/go-infra-backend/internal/utils"
	repoContracts "github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	svcContracts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type nestedUnitModel struct {
	ID       int64
	ParentID int64
	Name     sql.NullString
}

func (nestedUnitModel) TableName() string { return "nested_unit_models" }

type nestedUnitDTO struct {
	ID       string
	ParentID int64
	Name     string
}

type nestedTestMapper struct {
	applyParentScopesErr      error
	applyParentQueryScopesErr error
	applyQueryScopesErr       error
	getModelKeyErr            error
}

func (m nestedTestMapper) ToPostModel(dto nestedUnitDTO, model *nestedUnitModel) {
	m.ToModel(&dto, model)
}
func (m nestedTestMapper) ToPatchModel(dto nestedUnitDTO, model *nestedUnitModel) {
	m.ToModel(&dto, model)
}

func (m nestedTestMapper) ToModel(dto *nestedUnitDTO, model *nestedUnitModel) {
	model.Name = utils.ToNullString(dto.Name)
}

func (m nestedTestMapper) ToDTO(model *nestedUnitModel, dto *nestedUnitDTO) {
	dto.ID = strconv.FormatInt(model.ID, 10)
	dto.ParentID = model.ParentID
	dto.Name = utils.FromNullString(model.Name)
}

func (m nestedTestMapper) GetModelKey(id string) (map[string]any, error) {
	if m.getModelKeyErr != nil {
		return nil, m.getModelKeyErr
	}
	parsedID, err := strconv.ParseInt(id, 10, 64)
	if err != nil || parsedID <= 0 {
		return nil, errors.New("invalid id")
	}
	return map[string]any{"id": parsedID}, nil
}

func (m nestedTestMapper) ApplyQueryScopes(queryScopes map[string]any) (map[string]any, error) {
	if m.applyQueryScopesErr != nil {
		return nil, m.applyQueryScopesErr
	}
	out := make(map[string]any, len(queryScopes))
	for key, value := range queryScopes {
		out[key] = value
	}
	return out, nil
}

func (m nestedTestMapper) ApplyParentQueryScopes(parentID string, queryScopes map[string]any) (map[string]any, error) {
	if m.applyParentQueryScopesErr != nil {
		return nil, m.applyParentQueryScopesErr
	}
	mappedScopes, err := m.ApplyQueryScopes(queryScopes)
	if err != nil {
		return nil, err
	}
	mappedScopes["parent_id"] = parentID
	return mappedScopes, nil
}

func (m nestedTestMapper) ApplyParentScopes(parentID string, model *nestedUnitModel) error {
	if m.applyParentScopesErr != nil {
		return m.applyParentScopesErr
	}
	parsedID, err := strconv.ParseInt(parentID, 10, 64)
	if err != nil || parsedID <= 0 {
		return errors.New("invalid parent id")
	}
	model.ParentID = parsedID
	return nil
}

func TestNestedService_CreateNested(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository[nestedUnitModel](ctrl)
	svc := NewNestedServiceImpl(repoContracts.Repository[nestedUnitModel](repo), nestedTestMapper{})
	ctx := context.Background()

	dto := nestedUnitDTO{Name: "child"}
	repo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, inOut *nestedUnitModel) error {
		assert.Equal(t, int64(42), inOut.ParentID)
		assert.Equal(t, "child", inOut.Name.String)
		inOut.ID = 7
		return nil
	})

	out, err := svc.CreateNested(ctx, "42", dto)
	assert.NoError(t, err)
	assert.Equal(t, "7", out.ID)
	assert.Equal(t, int64(42), out.ParentID)
	assert.Equal(t, "child", out.Name)
}

func TestNestedService_CreateNested_Errors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	dto := nestedUnitDTO{Name: "child"}

	repo1 := NewMockRepository[nestedUnitModel](ctrl)
	svc1 := NewNestedServiceImpl(repoContracts.Repository[nestedUnitModel](repo1), nestedTestMapper{applyParentScopesErr: errors.New("parent error")})
	_, err := svc1.CreateNested(ctx, "42", dto)
	assert.Error(t, err)

	repo2 := NewMockRepository[nestedUnitModel](ctrl)
	svc2 := NewNestedServiceImpl(repoContracts.Repository[nestedUnitModel](repo2), nestedTestMapper{})
	repo2.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("repo error"))
	_, err = svc2.CreateNested(ctx, "42", dto)
	assert.Error(t, err)
}

func TestNestedService_ListNested(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository[nestedUnitModel](ctrl)
	svc := NewNestedServiceImpl(repoContracts.Repository[nestedUnitModel](repo), nestedTestMapper{})
	ctx := context.Background()

	params := svcContracts.ListParams{QueryParams: svcContracts.QueryParams{Scopes: map[string]any{"status": "active"}}}

	repo.EXPECT().List(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, listParams repoContracts.ListParams) ([]nestedUnitModel, int64, error) {
		assert.Equal(t, 1, listParams.Page)
		assert.Equal(t, 10, listParams.Size)
		assert.Equal(t, "created_at", listParams.OrderBy)
		assert.Equal(t, "asc", listParams.Order)
		assert.Equal(t, "active", listParams.QueryParams.Scopes["status"])
		assert.Equal(t, "99", listParams.QueryParams.Scopes["parent_id"])

		return []nestedUnitModel{{ID: 1, ParentID: 99, Name: utils.ToNullString("one")}}, 1, nil
	})

	list, total, err := svc.ListNested(ctx, "99", params)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, list, 1)
	assert.Equal(t, "1", list[0].ID)
	assert.Equal(t, int64(99), list[0].ParentID)
	assert.Equal(t, "one", list[0].Name)
}

func TestNestedService_ListNested_Errors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	params := svcContracts.ListParams{QueryParams: svcContracts.QueryParams{Scopes: map[string]any{"status": "active"}}}

	repo1 := NewMockRepository[nestedUnitModel](ctrl)
	svc1 := NewNestedServiceImpl(repoContracts.Repository[nestedUnitModel](repo1), nestedTestMapper{applyParentQueryScopesErr: errors.New("parent query error")})
	_, _, err := svc1.ListNested(ctx, "99", params)
	assert.Error(t, err)

	repo2 := NewMockRepository[nestedUnitModel](ctrl)
	svc2 := NewNestedServiceImpl(repoContracts.Repository[nestedUnitModel](repo2), nestedTestMapper{})
	repo2.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, int64(0), errors.New("repo error"))
	_, _, err = svc2.ListNested(ctx, "99", params)
	assert.Error(t, err)
}

func TestNestedService_DelegatesBaseMethods(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository[nestedUnitModel](ctrl)
	svc := NewNestedServiceImpl(repoContracts.Repository[nestedUnitModel](repo), nestedTestMapper{})
	ctx := context.Background()

	repo.EXPECT().GetByID(gomock.Any(), map[string]any{"id": int64(5)}).Return(nestedUnitModel{ID: 5, ParentID: 42, Name: utils.ToNullString("get")}, nil)
	got, err := svc.GetByID(ctx, "5")
	assert.NoError(t, err)
	assert.Equal(t, "5", got.ID)
	assert.Equal(t, int64(42), got.ParentID)

	upd := nestedUnitDTO{Name: "upd"}
	repo.EXPECT().Update(gomock.Any(), map[string]any{"id": int64(5)}, gomock.Any()).DoAndReturn(func(_ context.Context, _ map[string]any, inOut *nestedUnitModel) error {
		inOut.ID = 5
		inOut.ParentID = 42
		return nil
	})
	updOut, err := svc.Update(ctx, "5", upd)
	assert.NoError(t, err)
	assert.Equal(t, "5", updOut.ID)
	assert.Equal(t, int64(42), updOut.ParentID)

	repo.EXPECT().Delete(gomock.Any(), map[string]any{"id": int64(5)}).Return(nil)
	err = svc.Delete(ctx, "5")
	assert.NoError(t, err)
}
