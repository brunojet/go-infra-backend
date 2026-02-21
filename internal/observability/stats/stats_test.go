package stats

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOnqueueExport_Idempotent(t *testing.T) {
	s := NewObservabilityStats("id", "impl", "ep", false, "v1", 3)

	done := s.OnqueueExport(5)
	stats := s.GetAdapterStats()[0]
	require.Equal(t, 5, stats.Onqueue, "expected onqueue 5 after OnqueueExport")

	done(nil)
	stats = s.GetAdapterStats()[0]
	require.Equal(t, 0, stats.Onqueue, "expected onqueue 0 after done")

	// call done a second time - should not decrement again (idempotent)
	done(nil)
	stats = s.GetAdapterStats()[0]
	require.Equal(t, 0, stats.Onqueue, "expected onqueue remain 0 after second done")
}

func TestOnqueueExport_MultipleClosures_NoUnderflow(t *testing.T) {
	s := NewObservabilityStats("id", "impl", "ep", false, "v1", 3)

	doneA := s.OnqueueExport(8)
	doneB := s.OnqueueExport(8)

	stats := s.GetAdapterStats()[0]
	require.Equal(t, 16, stats.Onqueue, "expected onqueue 16 after two OnqueueExport")

	doneA(nil)
	stats = s.GetAdapterStats()[0]
	require.Equal(t, 8, stats.Onqueue, "expected onqueue 8 after doneA")

	doneB(nil)
	stats = s.GetAdapterStats()[0]
	require.Equal(t, 0, stats.Onqueue, "expected onqueue 0 after doneB")
}

func TestShutdownDoesNotPanic(t *testing.T) {
	s := NewObservabilityStats("id", "impl", "ep", false, "v1", 3)
	err := s.Shutdown(context.Background())
	require.NoError(t, err)
}

func TestOnqueueExport_ConcurrentRandomized(t *testing.T) {
	s := NewObservabilityStats("id", "impl", "ep", false, "v1", 3)

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	const N = 500

	type job struct {
		results []bool          // true => error, false => success; results[0] is the first call
		delays  []time.Duration // per-call delays
	}

	jobs := make([]job, 0, N)
	var expectedOnqueue, expectedFailed, expectedSent int

	// pre-generate scenarios and delays to avoid concurrent use of rand
	for i := 0; i < N; i++ {
		calls := rnd.Intn(5) + 1
		res := make([]bool, calls)
		delays := make([]time.Duration, calls)
		for j := 0; j < calls; j++ {
			res[j] = rnd.Intn(2) == 0
			delays[j] = time.Duration(100+rnd.Intn(401)) * time.Millisecond // 100-500ms
		}
		jobs = append(jobs, job{results: res, delays: delays})
		expectedOnqueue += calls
		if res[0] {
			expectedFailed += calls
		} else {
			expectedSent += calls
		}
	}

	var wg sync.WaitGroup
	wg.Add(len(jobs))
	for _, jb := range jobs {
		jb := jb
		go func() {
			defer wg.Done()
			done := s.OnqueueExport(len(jb.results))
			// call done sequentially according to pre-generated results with delays
			for idx, isErr := range jb.results {
				time.Sleep(jb.delays[idx])
				if isErr {
					done(errors.New("boom"))
				} else {
					done(nil)
				}
			}
		}()
	}
	wg.Wait()

	stats := s.GetAdapterStats()[0]
	assert.GreaterOrEqual(t, stats.Onqueue, 0, "underflow detected: onqueue negative")
	require.Equal(t, 0, stats.Onqueue, "expected onqueue 0 after all dones")

	require.Equal(t, expectedFailed, stats.Failed, "failed count should match expected based on first-call errors")
	require.Equal(t, expectedSent, stats.Sent, "sent count should match expected based on first-call successes")

	health := s.GetHealth()
	// recorded errors may be removed by successful reports; ensure it's not above capacity
	require.LessOrEqual(t, len(health.LastErrors), 3, "recorded errors should be capped by maxErrors")
}

func TestOnqueueExport_IgnoreNonPositiveCount(t *testing.T) {
	s := NewObservabilityStats("id", "impl", "ep", false, "v1", 3)

	// count <= 0 should be ignored and return a no-op
	done := s.OnqueueExport(0)
	done(nil)
	done(nil)

	stats := s.GetAdapterStats()[0]
	require.Equal(t, 0, stats.Onqueue, "expected onqueue 0 for non-positive count")
}

func TestOnqueueExport_ErrorReportedButCountZeroed(t *testing.T) {
	s := NewObservabilityStats("id", "impl", "ep", false, "v1", 3)

	done := s.OnqueueExport(5)
	// call done with an error; implementation zeroes the count before reporting
	done(errors.New("boom"))

	stats := s.GetAdapterStats()[0]
	// failed should increment by the original count (implementation does not zero count)
	require.Equal(t, 5, stats.Failed, "expected failed 5 after error report")

	health := s.GetHealth()
	require.Len(t, health.LastErrors, 1, "expected 1 last error recorded")
}
