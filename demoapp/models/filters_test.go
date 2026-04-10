package models

import (
	"database/sql"
	"fmt"
	"sync"
	"testing"

	"github.com/brunojet/go-infra-backend/pkg/testutil/dbtest"
	"github.com/stretchr/testify/require"
)

func TestFilterType_TableName(t *testing.T) {
	require.Equal(t, tableFilterType, (FilterType{}).TableName())
}

func TestFilterType_BeforeCreate_UpsertCollision(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &FilterType{})

	_ = createFilterTypeWithDescription(t, gdb, "ft-upsert", "v1")
	dup := FilterType{Name: sql.NullString{String: "ft-upsert", Valid: true}, Description: sql.NullString{String: "v2", Valid: true}}
	require.NoError(t, gdb.Create(&dup).Error)

	var count int64
	require.NoError(t, gdb.Model(&FilterType{}).Where("name = ?", "ft-upsert").Count(&count).Error)
	require.EqualValues(t, 1, count)

	var got FilterType
	require.NoError(t, gdb.Where("name = ?", "ft-upsert").First(&got).Error)
	require.Equal(t, "v1", got.Description.String)
}

func TestFilter_TableName(t *testing.T) {
	require.Equal(t, tableFilter, (Filter{}).TableName())
}

func TestFilter_BeforeCreate_UpsertCollision(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &FilterType{}, &Filter{})

	ft := createFilterType(t, gdb, "ft-filter-upsert")
	_ = createFilterWithDescription(t, gdb, ft.FilterTypeId, "color", "v1")
	dup := Filter{FilterTypeId: ft.FilterTypeId, Name: sql.NullString{String: "color", Valid: true}, Description: sql.NullString{String: "v2", Valid: true}}
	require.NoError(t, gdb.Create(&dup).Error)

	var count int64
	require.NoError(t, gdb.Model(&Filter{}).Where("name = ? AND filter_type_id = ?", "color", ft.FilterTypeId).Count(&count).Error)
	require.EqualValues(t, 1, count)

	var got Filter
	require.NoError(t, gdb.Where("name = ? AND filter_type_id = ?", "color", ft.FilterTypeId).First(&got).Error)
	require.Equal(t, "v1", got.Description.String)
}

func TestFilter_ConcurrentCreateSameNameAndType_UpsertSingleRow(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &FilterType{}, &Filter{})
	ft := createFilterType(t, gdb, "ft-concurrent")

	const workers = 16
	start := make(chan struct{})
	errs := make(chan error, workers)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		idx := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			err := dbtest.RunWithSQLiteRetry(func() error {
				entity := Filter{
					FilterTypeId: ft.FilterTypeId,
					Name:         sql.NullString{String: "color-concurrent", Valid: true},
					Description:  sql.NullString{String: fmt.Sprintf("writer-%d", idx), Valid: true},
				}
				return gdb.Create(&entity).Error
			})
			errs <- err
		}()
	}

	close(start)
	wg.Wait()
	close(errs)

	for err := range errs {
		require.NoError(t, err)
	}

	var count int64
	require.NoError(t, gdb.Model(&Filter{}).Where("name = ? AND filter_type_id = ?", "color-concurrent", ft.FilterTypeId).Count(&count).Error)
	require.EqualValues(t, 1, count)
}

func TestFilter_OneShotCreateWithFilterType(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &FilterType{}, &Filter{})

	oneShot := createFilterOneShot(t, gdb, "ft-oneshot", "f-oneshot")

	var loaded Filter
	require.NoError(t, gdb.Preload("FilterType").
		Where("filter_id = ?", oneShot.FilterId).
		First(&loaded).Error)
	require.NotNil(t, loaded.FilterType)
	require.Equal(t, "ft-oneshot", loaded.FilterType.Name.String)
	require.Equal(t, "f-oneshot", loaded.Name.String)
}

func TestFilter_OneShotCreateWithNestedFilterTypeAssociation(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &FilterType{}, &Filter{})

	f := createFilterOneShotWithNestedFilterType(t, gdb, "nested-type", "nested-filter")

	var loaded Filter
	require.NoError(t, gdb.Preload("FilterType").Where("name = ?", "nested-filter").First(&loaded).Error)
	require.NotNil(t, loaded.FilterType)
	require.Equal(t, "nested-type", loaded.FilterType.Name.String)
	require.Equal(t, f.FilterId, loaded.FilterId)
}
