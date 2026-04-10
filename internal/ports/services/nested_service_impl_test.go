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
	gorm "gorm.io/gorm"
)

type nestedUnitModel struct {
	ID       int64
	ParentID int64
	Name     sql.NullString
}

func (nestedUnitModel) TableName() string { return "nested_unit_models" }

func (m nestedUnitModel) WhereOnConflict(tx *gorm.DB) *gorm.DB {
	return tx
}

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

func (m nestedTestMapper) ToPostModel(dto nestedUnitDTO, model *nestedUnitModel) error {
	return m.ToModel(&dto, model)
}
func (m nestedTestMapper) ToPatchModel(dto nestedUnitDTO, model *nestedUnitModel) error {
	return m.ToModel(&dto, model)
}

func (m nestedTestMapper) ToModel(dto *nestedUnitDTO, model *nestedUnitModel) error {
	model.Name = utils.ToNullString(dto.Name)
	return nil
}

func (m nestedTestMapper) ToDTO(model *nestedUnitModel, dto *nestedUnitDTO) error {
	dto.ID = strconv.FormatInt(model.ID, 10)
	dto.ParentID = model.ParentID
	dto.Name = utils.FromNullString(model.Name)
	return nil
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
	ExpectCreateWithTx(repo, func(inOut *nestedUnitModel) {
		assert.Equal(t, int64(42), inOut.ParentID)
		assert.Equal(t, "child", inOut.Name.String)
		inOut.ID = 7
	})

	var out nestedUnitDTO
	err := svc.CreateNested(ctx, "42", dto, &out)
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
	var out nestedUnitDTO
	err := svc1.CreateNested(ctx, "42", dto, &out)
	assert.Error(t, err)

	repo2 := NewMockRepository[nestedUnitModel](ctrl)
	svc2 := NewNestedServiceImpl(repoContracts.Repository[nestedUnitModel](repo2), nestedTestMapper{})
	ExpectCreateWithTx(repo2, nil).Return(errors.New("repo error"))
	err = svc2.CreateNested(ctx, "42", dto, &out)
	assert.Error(t, err)
}

func TestNestedService_ListNested(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository[nestedUnitModel](ctrl)
	svc := NewNestedServiceImpl(repoContracts.Repository[nestedUnitModel](repo), nestedTestMapper{})
	ctx := context.Background()

	params := svcContracts.ListParams{QueryParams: svcContracts.QueryParams{Scopes: map[string]any{"status": "active"}}}

	ExpectListWithTx(repo, func(out *[]nestedUnitModel) error {
		*out = []nestedUnitModel{{ID: 1, ParentID: 99, Name: utils.ToNullString("one")}}
		return nil
	})
	outList := make([]nestedUnitDTO, 1)
	total, err := svc.ListNested(ctx, "99", params, &outList)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, outList, 1)
	assert.Equal(t, "1", outList[0].ID)
	assert.Equal(t, int64(99), outList[0].ParentID)
	assert.Equal(t, "one", outList[0].Name)
}

func TestNestedService_ListNested_Errors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	params := svcContracts.ListParams{QueryParams: svcContracts.QueryParams{Scopes: map[string]any{"status": "active"}}}

	repo1 := NewMockRepository[nestedUnitModel](ctrl)
	svc1 := NewNestedServiceImpl(repoContracts.Repository[nestedUnitModel](repo1), nestedTestMapper{applyParentQueryScopesErr: errors.New("parent query error")})
	outList1 := make([]nestedUnitDTO, 1)
	_, err := svc1.ListNested(ctx, "99", params, &outList1)
	assert.Error(t, err)

	repo2 := NewMockRepository[nestedUnitModel](ctrl)
	svc2 := NewNestedServiceImpl(repoContracts.Repository[nestedUnitModel](repo2), nestedTestMapper{})
	ExpectListWithTx(repo2, nil).Return(int64(0), errors.New("repo error"))
	outList2 := make([]nestedUnitDTO, 1)
	_, err = svc2.ListNested(ctx, "99", params, &outList2)
	assert.Error(t, err)
}

func TestNestedService_DelegatesBaseMethods(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository[nestedUnitModel](ctrl)
	svc := NewNestedServiceImpl(repoContracts.Repository[nestedUnitModel](repo), nestedTestMapper{})
	ctx := context.Background()

	ExpectGetByIDWithTx(repo, func(out *nestedUnitModel) {
		*out = nestedUnitModel{ID: 5, ParentID: 42, Name: utils.ToNullString("get")}
	})
	var got nestedUnitDTO
	err := svc.GetByID(ctx, "5", &got)
	assert.NoError(t, err)
	assert.Equal(t, "5", got.ID)
	assert.Equal(t, int64(42), got.ParentID)

	upd := nestedUnitDTO{Name: "upd"}
	ExpectUpdateWithTx(repo, func(inOut *nestedUnitModel) {
		inOut.ID = 5
		inOut.ParentID = 42
	})
	var updOut nestedUnitDTO
	err = svc.Update(ctx, "5", upd, &updOut)
	assert.NoError(t, err)
	assert.Equal(t, "5", updOut.ID)
	assert.Equal(t, int64(42), updOut.ParentID)

	ExpectDeleteWithTx(repo, nil)
	err = svc.Delete(ctx, "5")
	assert.NoError(t, err)
}
