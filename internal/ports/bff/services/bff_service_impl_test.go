package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	bffrpocts "github.com/brunojet/go-infra-backend/pkg/ports/bff/repositories/contracts"
	bffsvccts "github.com/brunojet/go-infra-backend/pkg/ports/bff/services/contracts"
	svccts "github.com/brunojet/go-infra-backend/pkg/ports/services/contracts"
)

// ---------------------------------------------------------------------------
// test fixtures — domain types
// ---------------------------------------------------------------------------

type createDTO struct{ Name string }
type readDTO struct{ ID, Name string }
type updateDTO struct{ Name string }

// ---------------------------------------------------------------------------
// test fixtures — upstream entities
// ---------------------------------------------------------------------------

type upstreamEntity struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Total int64  `json:"total"` // meta field used by ExtractUpstreamTotal
}

func (upstreamEntity) ResourceName() string { return "items" }

// ---------------------------------------------------------------------------
// stub mapper (flat — same type for CE/RE/UE)
// ---------------------------------------------------------------------------

type stubMapper struct {
	toPostErr  error
	toPatchErr error
	toDTOErr   error
	getIDErr   error
	scopesErr  error
}

var _ bffsvccts.BffServiceMapper[createDTO, readDTO, updateDTO, upstreamEntity, upstreamEntity, upstreamEntity] = (*stubMapper)(nil)

func (m *stubMapper) ToUpstreamPost(dto createDTO, up *upstreamEntity) error {
	if m.toPostErr != nil {
		return m.toPostErr
	}
	up.Name = dto.Name
	return nil
}

func (m *stubMapper) ToUpstreamPatch(dto updateDTO, up *upstreamEntity) error {
	if m.toPatchErr != nil {
		return m.toPatchErr
	}
	up.Name = dto.Name
	return nil
}

func (m *stubMapper) ToDomainDTO(up *upstreamEntity, dto *readDTO) error {
	if m.toDTOErr != nil {
		return m.toDTOErr
	}
	dto.ID = up.ID
	dto.Name = up.Name
	return nil
}

func (m *stubMapper) GetUpstreamID(id string) (string, error) {
	if m.getIDErr != nil {
		return "", m.getIDErr
	}
	return id, nil
}

func (m *stubMapper) ApplyQueryScopes(scopes map[string]any) (map[string]any, error) {
	if m.scopesErr != nil {
		return nil, m.scopesErr
	}
	return scopes, nil
}

func (m *stubMapper) ExtractUpstreamTotal(upstream []upstreamEntity) int64 {
	if len(upstream) == 0 {
		return 0
	}
	return upstream[0].Total
}

// ---------------------------------------------------------------------------
// stub repository
// ---------------------------------------------------------------------------

type stubRepo struct {
	createErr   error
	getByIDErr  error
	listErr     error
	updateErr   error
	deleteErr   error
	createResp  upstreamEntity
	getByIDResp upstreamEntity
	listResp    []upstreamEntity
	updateResp  upstreamEntity
}

var _ bffrpocts.BffRepository[upstreamEntity, upstreamEntity, upstreamEntity] = (*stubRepo)(nil)

func (r *stubRepo) Create(_ context.Context, _ upstreamEntity, downstream *upstreamEntity) error {
	if r.createErr != nil {
		return r.createErr
	}
	*downstream = r.createResp
	return nil
}

func (r *stubRepo) GetByID(_ context.Context, _ string, downstream *upstreamEntity) error {
	if r.getByIDErr != nil {
		return r.getByIDErr
	}
	*downstream = r.getByIDResp
	return nil
}

func (r *stubRepo) List(_ context.Context, _ bffrpocts.BffListParams, downstream *[]upstreamEntity) error {
	if r.listErr != nil {
		return r.listErr
	}
	*downstream = r.listResp
	return nil
}

func (r *stubRepo) Update(_ context.Context, _ string, _ upstreamEntity, downstream *upstreamEntity) error {
	if r.updateErr != nil {
		return r.updateErr
	}
	*downstream = r.updateResp
	return nil
}

func (r *stubRepo) Delete(_ context.Context, _ string) error {
	return r.deleteErr
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func newSvc() (svccts.Service[createDTO, readDTO, updateDTO], *stubRepo, *stubMapper) {
	repo := &stubRepo{}
	mapper := &stubMapper{}
	svc := NewBffServiceImpl(repo, mapper)
	return svc, repo, mapper
}

// ---------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------

func TestBffService_Create_Success(t *testing.T) {
	svc, repo, _ := newSvc()
	repo.createResp = upstreamEntity{ID: "1", Name: "foo"}

	var resp readDTO
	err := svc.Create(context.Background(), createDTO{Name: "foo"}, &resp)

	require.NoError(t, err)
	assert.Equal(t, "1", resp.ID)
	assert.Equal(t, "foo", resp.Name)
}

func TestBffService_Create_MapperError(t *testing.T) {
	svc, _, mapper := newSvc()
	mapper.toPostErr = errors.New("map error")

	var resp readDTO
	err := svc.Create(context.Background(), createDTO{}, &resp)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "map error")
}

func TestBffService_Create_RepoError(t *testing.T) {
	svc, repo, _ := newSvc()
	repo.createErr = errors.New("repo error")

	var resp readDTO
	err := svc.Create(context.Background(), createDTO{}, &resp)
	require.Error(t, err)
}

// ---------------------------------------------------------------------------
// GetByID
// ---------------------------------------------------------------------------

func TestBffService_GetByID_Success(t *testing.T) {
	svc, repo, _ := newSvc()
	repo.getByIDResp = upstreamEntity{ID: "42", Name: "bar"}

	var resp readDTO
	err := svc.GetByID(context.Background(), "42", &resp)

	require.NoError(t, err)
	assert.Equal(t, "42", resp.ID)
	assert.Equal(t, "bar", resp.Name)
}

func TestBffService_GetByID_IDError(t *testing.T) {
	svc, _, mapper := newSvc()
	mapper.getIDErr = errors.New("bad id")

	var resp readDTO
	err := svc.GetByID(context.Background(), "bad", &resp)
	require.Error(t, err)
}

func TestBffService_GetByID_RepoError(t *testing.T) {
	svc, repo, _ := newSvc()
	repo.getByIDErr = errors.New("not found")

	var resp readDTO
	err := svc.GetByID(context.Background(), "99", &resp)
	require.Error(t, err)
}

// ---------------------------------------------------------------------------
// List
// ---------------------------------------------------------------------------

func TestBffService_List_Success(t *testing.T) {
	svc, repo, _ := newSvc()
	repo.listResp = []upstreamEntity{
		{ID: "1", Name: "a", Total: 42},
		{ID: "2", Name: "b", Total: 42},
	}

	var responses []readDTO
	total, err := svc.List(context.Background(), svccts.ListParams{Page: 1, Size: 10}, &responses)

	require.NoError(t, err)
	assert.EqualValues(t, 42, total)
	assert.Len(t, responses, 2)
	assert.Equal(t, "a", responses[0].Name)
}

func TestBffService_List_EmptyResult(t *testing.T) {
	svc, repo, _ := newSvc()
	repo.listResp = []upstreamEntity{}

	var responses []readDTO
	total, err := svc.List(context.Background(), svccts.ListParams{Page: 1, Size: 10}, &responses)

	require.NoError(t, err)
	assert.EqualValues(t, 0, total)
	assert.Empty(t, responses)
}

func TestBffService_List_ScopesError(t *testing.T) {
	svc, _, mapper := newSvc()
	mapper.scopesErr = errors.New("scope error")

	var responses []readDTO
	_, err := svc.List(context.Background(), svccts.ListParams{}, &responses)
	require.Error(t, err)
}

func TestBffService_List_RepoError(t *testing.T) {
	svc, repo, _ := newSvc()
	repo.listErr = errors.New("repo error")

	var responses []readDTO
	_, err := svc.List(context.Background(), svccts.ListParams{}, &responses)
	require.Error(t, err)
}

// ---------------------------------------------------------------------------
// Update
// ---------------------------------------------------------------------------

func TestBffService_Update_Success(t *testing.T) {
	svc, repo, _ := newSvc()
	repo.updateResp = upstreamEntity{ID: "5", Name: "updated"}

	var resp readDTO
	err := svc.Update(context.Background(), "5", updateDTO{Name: "updated"}, &resp)

	require.NoError(t, err)
	assert.Equal(t, "5", resp.ID)
	assert.Equal(t, "updated", resp.Name)
}

func TestBffService_Update_IDError(t *testing.T) {
	svc, _, mapper := newSvc()
	mapper.getIDErr = errors.New("bad id")

	var resp readDTO
	err := svc.Update(context.Background(), "bad", updateDTO{}, &resp)
	require.Error(t, err)
}

func TestBffService_Update_RepoError(t *testing.T) {
	svc, repo, _ := newSvc()
	repo.updateErr = errors.New("upstream error")

	var resp readDTO
	err := svc.Update(context.Background(), "5", updateDTO{}, &resp)
	require.Error(t, err)
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

func TestBffService_Delete_Success(t *testing.T) {
	svc, _, _ := newSvc()
	err := svc.Delete(context.Background(), "7")
	require.NoError(t, err)
}

func TestBffService_Delete_RepoError(t *testing.T) {
	svc, repo, _ := newSvc()
	repo.deleteErr = errors.New("delete failed")

	err := svc.Delete(context.Background(), "7")
	require.Error(t, err)
}

// ---------------------------------------------------------------------------
// toBffListParams — normalisation
// ---------------------------------------------------------------------------

func TestToBffListParams_DefaultsApplied(t *testing.T) {
	p := toBffListParams(svccts.ListParams{}, nil)

	assert.Equal(t, 1, p.Page)
	assert.Equal(t, 10, p.Size)
	assert.Equal(t, "created_at", p.OrderBy)
	assert.Equal(t, "asc", p.Order)
}

func TestToBffListParams_DescOrderPreserved(t *testing.T) {
	p := toBffListParams(svccts.ListParams{Order: "DESC"}, nil)
	assert.Equal(t, "desc", p.Order)
}

func TestToBffListParams_ScopesOverride(t *testing.T) {
	scopes := map[string]any{"x": 1}
	p := toBffListParams(svccts.ListParams{}, scopes)
	assert.Equal(t, scopes, p.Scopes)
}

// ---------------------------------------------------------------------------
// BffNestedService — stub nested mapper & repo
// ---------------------------------------------------------------------------

type stubNestedMapper struct {
	stubMapper
	getParentIDErr    error
	parentScopesErr   error
	parentUpstreamErr error
}

var _ bffsvccts.BffNestedServiceMapper[createDTO, readDTO, updateDTO, upstreamEntity, upstreamEntity, upstreamEntity] = (*stubNestedMapper)(nil)

func (m *stubNestedMapper) GetUpstreamParentID(parentID string) (string, error) {
	if m.getParentIDErr != nil {
		return "", m.getParentIDErr
	}
	return parentID, nil
}

func (m *stubNestedMapper) ApplyParentQueryScopes(parentID string, scopes map[string]any) (map[string]any, error) {
	if m.parentScopesErr != nil {
		return nil, m.parentScopesErr
	}
	return scopes, nil
}

func (m *stubNestedMapper) ApplyUpstreamParentScopes(parentID string, upstream *upstreamEntity) error {
	if m.parentUpstreamErr != nil {
		return m.parentUpstreamErr
	}
	return nil
}

type stubNestedRepo struct {
	stubRepo
	createNestedErr  error
	listNestedErr    error
	createNestedResp upstreamEntity
	listNestedResp   []upstreamEntity
}

var _ bffrpocts.BffNestedRepository[upstreamEntity, upstreamEntity, upstreamEntity] = (*stubNestedRepo)(nil)

func (r *stubNestedRepo) CreateNested(_ context.Context, _ string, _ upstreamEntity, downstream *upstreamEntity) error {
	if r.createNestedErr != nil {
		return r.createNestedErr
	}
	*downstream = r.createNestedResp
	return nil
}

func (r *stubNestedRepo) ListNested(_ context.Context, _ string, _ bffrpocts.BffListParams, downstream *[]upstreamEntity) error {
	if r.listNestedErr != nil {
		return r.listNestedErr
	}
	*downstream = r.listNestedResp
	return nil
}

func newNestedSvc() (svccts.NestedService[createDTO, readDTO, updateDTO], *stubNestedRepo, *stubNestedMapper) {
	repo := &stubNestedRepo{}
	mapper := &stubNestedMapper{}
	svc := NewBffNestedServiceImpl(repo, mapper)
	return svc, repo, mapper
}

func TestBffNestedService_CreateNested_Success(t *testing.T) {
	svc, repo, _ := newNestedSvc()
	repo.createNestedResp = upstreamEntity{ID: "c1", Name: "child"}

	var resp readDTO
	err := svc.CreateNested(context.Background(), "p1", createDTO{Name: "child"}, &resp)

	require.NoError(t, err)
	assert.Equal(t, "c1", resp.ID)
	assert.Equal(t, "child", resp.Name)
}

func TestBffNestedService_CreateNested_ParentIDError(t *testing.T) {
	svc, _, mapper := newNestedSvc()
	mapper.getParentIDErr = errors.New("bad parent")

	var resp readDTO
	err := svc.CreateNested(context.Background(), "bad", createDTO{}, &resp)
	require.Error(t, err)
}

func TestBffNestedService_ListNested_Success(t *testing.T) {
	svc, repo, _ := newNestedSvc()
	repo.listNestedResp = []upstreamEntity{
		{ID: "c1", Name: "child1", Total: 5},
		{ID: "c2", Name: "child2", Total: 5},
	}

	var responses []readDTO
	total, err := svc.ListNested(context.Background(), "p1", svccts.ListParams{Page: 1, Size: 10}, &responses)

	require.NoError(t, err)
	assert.EqualValues(t, 5, total)
	assert.Len(t, responses, 2)
}

func TestBffNestedService_ListNested_RepoError(t *testing.T) {
	svc, repo, _ := newNestedSvc()
	repo.listNestedErr = errors.New("upstream error")

	var responses []readDTO
	_, err := svc.ListNested(context.Background(), "p1", svccts.ListParams{}, &responses)
	require.Error(t, err)
}
