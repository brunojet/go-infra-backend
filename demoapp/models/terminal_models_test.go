package models

import (
	"fmt"
	"sync"
	"testing"

	"github.com/brunojet/go-infra-backend/pkg/testutil/dbtest"
	"github.com/stretchr/testify/require"
)

func TestTerminalModel_TableName(t *testing.T) {
	require.Equal(t, tableTerminalModel, (TerminalModel{}).TableName())
}

func TestTerminalModel_BeforeCreate_UpsertCollision(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &TerminalModel{})

	_ = createTerminalModelWithDescription(t, gdb, "tm-upsert", "v1")
	dup := buildTerminalModel("tm-upsert", "v2")
	require.NoError(t, gdb.Create(&dup).Error)

	var count int64
	require.NoError(t, gdb.Model(&TerminalModel{}).Where("name = ?", "tm-upsert").Count(&count).Error)
	require.EqualValues(t, 1, count)

	var got TerminalModel
	require.NoError(t, gdb.Where("name = ?", "tm-upsert").First(&got).Error)
	require.Equal(t, "v1", got.Description.String)
}

func TestTerminalModelConfiguration_TableName(t *testing.T) {
	require.Equal(t, tableTerminalModelConfiguration, (TerminalModelConfiguration{}).TableName())
}

func TestTerminalModelConfiguration_BeforeCreate_UpsertCollision(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &TerminalModel{}, &TerminalModelConfiguration{})

	tm := createTerminalModel(t, gdb, "tm-cfg")
	_ = createTerminalModelConfiguration(t, gdb, tm.TerminalModelId, 7)
	dup := buildTerminalModelConfiguration(tm.TerminalModelId, 7)
	require.NoError(t, gdb.Create(&dup).Error)

	var count int64
	require.NoError(t, gdb.Model(&TerminalModelConfiguration{}).
		Where("terminal_model_id = ? AND integration_type = ?", tm.TerminalModelId, 7).
		Count(&count).Error)
	require.EqualValues(t, 1, count)
}

func TestTerminalModel_ConcurrentCreateSameName_UpsertSingleRow(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &TerminalModel{})

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
				entity := buildTerminalModel("tm-concurrent", fmt.Sprintf("writer-%d", idx))
				return gdb.Create(&entity).Error
			})
			errs <- err
		}()
	}

	close(start)
	wg.Wait()
	close(errs)

	for err := range errs {
		if err == nil {
			continue
		}
	}

	var count int64
	require.NoError(t, gdb.Model(&TerminalModel{}).Where("name = ?", "tm-concurrent").Count(&count).Error)
	require.EqualValues(t, 1, count)
}

func TestTerminalModelConfiguration_ConcurrentCreateSamePair_UpsertSingleRow(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &TerminalModel{}, &TerminalModelConfiguration{})

	tm := createTerminalModel(t, gdb, "tm-concurrent-cfg")

	const workers = 16
	start := make(chan struct{})
	errs := make(chan error, workers)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			err := dbtest.RunWithSQLiteRetry(func() error {
				cfg := buildTerminalModelConfiguration(tm.TerminalModelId, 11)
				return gdb.Create(&cfg).Error
			})
			errs <- err
		}()
	}

	close(start)
	wg.Wait()
	close(errs)

	for err := range errs {
		if err == nil {
			continue
		}
	}

	var count int64
	require.NoError(t, gdb.Model(&TerminalModelConfiguration{}).
		Where("terminal_model_id = ? AND integration_type = ?", tm.TerminalModelId, 11).
		Count(&count).Error)
	require.EqualValues(t, 1, count)
}

func TestTerminalModelConfiguration_OneShotCreateWithNestedTerminalModel(t *testing.T) {
	gdb := dbtest.OpenMemoryDB(t, &TerminalModel{}, &TerminalModelConfiguration{})

	oneShot := createTerminalModelConfigurationOneShot(t, gdb, "tm-oneshot-config", 21)

	var loaded TerminalModelConfiguration
	require.NoError(t, gdb.Preload("TerminalModel").
		Where("terminal_model_configuration_id = ?", oneShot.TerminalModelConfigurationId).
		First(&loaded).Error)
	require.NotNil(t, loaded.TerminalModel)
	require.Equal(t, "tm-oneshot-config", loaded.TerminalModel.Name.String)
	require.EqualValues(t, 21, loaded.IntegrationType.Int16)
}
