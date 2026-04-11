package stats

import (
	"context"
	"sync"
	"time"

	"github.com/brunojet/go-infra-backend/pkg/infra/observability/contracts"
)

const (
	maxErrorsInArray      = 5
	minErrorsForDegraded  = 3
	minErrorsForUnhealthy = maxErrorsInArray
)

// ObservabilityStats is a lightweight stats/health recorder for an exporter.
// It maintains a small circular buffer of recent errors to drive health state
// using the maxErrorsInArray rule.
type ObservabilityStats struct {
	mu             sync.Mutex
	id             string
	implementation string
	endpoint       string
	isInsecure     bool
	configVersion  string

	errors    []contracts.ObservabilityAdapterError
	maxErrors int

	// simple stats
	onqueue int
	sent    int
	failed  int

	state contracts.HealthState
}

// NewObservabilityStats creates a configured ObservabilityStats.
func NewObservabilityStats(id, implementation, endpoint string, isInsecure bool, configVersion string, maxErrors int) *ObservabilityStats {
	if maxErrors <= maxErrorsInArray {
		maxErrors = maxErrorsInArray
	}

	return &ObservabilityStats{
		id:             id,
		implementation: implementation,
		endpoint:       endpoint,
		isInsecure:     isInsecure,
		configVersion:  configVersion,
		errors:         make([]contracts.ObservabilityAdapterError, 0, maxErrors),
		maxErrors:      maxErrors,
		state:          contracts.HealthStatePending,
	}
}

func (t *ObservabilityStats) reportErrorUnsafe(count int, err error) {
	t.failed += count
	errEntry := contracts.ObservabilityAdapterError{
		Timestamp: time.Now(),
		Error:     err.Error(),
	}
	if len(t.errors) >= t.maxErrors {
		t.removeOldestErrorUnsafe()
	}
	t.errors = append(t.errors, errEntry)
}

func (t *ObservabilityStats) reportSuccessUnsafe(count int) {
	t.sent += count
	if len(t.errors) > 0 {
		t.removeOldestErrorUnsafe()
	}
}

func (t *ObservabilityStats) removeOldestErrorUnsafe() {
	if len(t.errors) == 0 {
		return
	}
	oldestIdx := 0
	oldestTS := t.errors[0].Timestamp
	for i := 1; i < len(t.errors); i++ {
		if t.errors[i].Timestamp.Before(oldestTS) {
			oldestIdx = i
			oldestTS = t.errors[i].Timestamp
		}
	}
	t.errors = append(t.errors[:oldestIdx], t.errors[oldestIdx+1:]...)
}

func (t *ObservabilityStats) reconcileStateUnsafe() {
	errsLen := len(t.errors)

	if errsLen >= minErrorsForUnhealthy {
		t.state = contracts.HealthStateUnhealthy
	} else if errsLen >= minErrorsForDegraded {
		t.state = contracts.HealthStateDegraded
	} else {
		t.state = contracts.HealthStateHealthy
	}
}

func (t *ObservabilityStats) processQueueResults(count int, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.onqueue -= count
	if err != nil {
		t.reportErrorUnsafe(count, err)
	} else {
		t.reportSuccessUnsafe(count)
	}
	t.reconcileStateUnsafe()
}

func (t *ObservabilityStats) OnqueueExport(count int) func(err error) {
	if count <= 0 {
		return func(err error) {}
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.onqueue += count

	var once sync.Once
	return func(err error) {
		once.Do(func() {
			t.processQueueResults(count, err)
		})
	}
}

// GetAdapterStats implements contracts.ObservabilityAdapter
func (t *ObservabilityStats) GetAdapterStats() []contracts.ObservabilityAdapterStats {
	t.mu.Lock()
	defer t.mu.Unlock()

	stats := contracts.ObservabilityAdapterStats{
		Endpoint:   t.endpoint,
		IsInsecure: t.isInsecure,
		Onqueue:    t.onqueue,
		Sent:       t.sent,
		Failed:     t.failed,
	}

	return []contracts.ObservabilityAdapterStats{stats}
}

// Shutdown implements contracts.ObservabilityAdapter
func (t *ObservabilityStats) Shutdown(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.state = contracts.HealthStateDefault
	return nil
}

// GetHealth returns the current health snapshot including recent errors.
func (t *ObservabilityStats) GetHealth() contracts.ObservabilityAdapterHealth {
	t.mu.Lock()
	defer t.mu.Unlock()

	// copy last errors for safe exposure
	lastErrs := make([]contracts.ObservabilityAdapterError, len(t.errors))
	copy(lastErrs, t.errors)

	return contracts.ObservabilityAdapterHealth{
		HealthState: t.state,
		LastErrors:  lastErrs,
	}
}
