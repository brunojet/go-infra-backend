package observability

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeAdapter struct {
	id     int
	record *[]int
	err    error
}

func (f *fakeAdapter) Shutdown(ctx context.Context) error {
	if f.record != nil {
		*f.record = append(*f.record, f.id)
	}
	return f.err
}

func Test_ObservabilityManager_ShutdownOrderAndErrorAggregation(t *testing.T) {
	mgr := NewObservabilityManager()
	var record []int

	a1 := &fakeAdapter{id: 1, record: &record}
	a2 := &fakeAdapter{id: 2, record: &record, err: errors.New("boom")}
	a3 := &fakeAdapter{id: 3, record: &record}

	mgr.RegisterAdapter(a1)
	mgr.RegisterAdapter(nil) // should be ignored
	mgr.RegisterAdapter(a2)
	mgr.RegisterAdapter(a3)

	err := mgr.Shutdown(context.Background())
	require.Error(t, err)

	// shutdown should be called in reverse registration order: a3, a2, a1
	require.Equal(t, []int{3, 2, 1}, record)
}
