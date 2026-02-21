package adapters

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	contracts "github.com/brunojet/go-infra-backend/pkg/events/contracts"
)

// FilesystemAdapter is a simple polling-based adapter that watches a directory
// tree and emits create/modify/delete events. It's intentionally small and
// dependency-free so it can be used immediately. For production usage replace
// with an evented watcher (fsnotify) or platform-specific implementation.
type FilesystemAdapter struct {
	dir         string
	interval    time.Duration
	knownFiles  map[string]os.FileInfo
	mu          sync.Mutex
	stoppedOnce sync.Once
}

// NewFilesystemAdapter creates a new adapter watching the provided directory.
// interval defines the polling interval used to scan the directory.
func NewFilesystemAdapter(dir string, interval time.Duration) *FilesystemAdapter {
	if interval <= 0 {
		interval = 1 * time.Second
	}
	return &FilesystemAdapter{
		dir:        dir,
		interval:   interval,
		knownFiles: make(map[string]os.FileInfo),
	}
}

func (f *FilesystemAdapter) Start(ctx context.Context) error {
	if f.dir == "" {
		return errors.New("directory not set")
	}
	files, err := f.listFiles()
	if err != nil {
		return err
	}
	f.replaceKnownFiles(files)
	return nil
}

func (f *FilesystemAdapter) Stop(ctx context.Context) error {
	return nil
}

// listFiles walks the directory and returns a map of relative path -> FileInfo.
// It centralizes the filesystem scanning logic used by Start, Poll and
// WaitForMessages.
func (f *FilesystemAdapter) listFiles() (map[string]os.FileInfo, error) {
	current := make(map[string]os.FileInfo)
	walkErr := filepath.Walk(f.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // ignore individual entry errors
		}
		if info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(f.dir, path)
		current[rel] = info
		return nil
	})
	if walkErr != nil {
		return nil, walkErr
	}
	return current, nil
}

// fsMessage is a lightweight message wrapper for filesystem events. Ack/Nack
// are no-ops because filesystem doesn't require explicit acknowledgement.
type fsMessage struct {
	evt      contracts.Event
	adapter  *FilesystemAdapter
	snapshot map[string]os.FileInfo // the 'current' map observed when creating this message
}

func (m *fsMessage) Event() contracts.Event { return m.evt }

// Ack updates the adapter knownFiles to the snapshot observed when the
// message was produced. This ensures at-least-once semantics: if the
// message is Nacked or the manager fails to process it, the adapter will
// re-emit the event on the next scan.
func (m *fsMessage) Ack(ctx context.Context) error {
	if m.adapter == nil {
		return fmt.Errorf("invalid adapter reference in message")
	}
	m.adapter.replaceKnownFiles(m.snapshot)
	return nil
}

// Nack signals failed processing. For filesystem adapter we simply do
// nothing here so the next WaitForMessages can re-detect the event.
func (m *fsMessage) Nack(ctx context.Context, delay time.Duration) error { return nil }

// ExtendVisibility is not applicable to filesystem polling adapter; no-op.
func (m *fsMessage) ExtendVisibility(ctx context.Context, d time.Duration) error { return nil }

// WaitForMessages blocks until at least one message is available or ctx is done.
// It returns all messages detected in the scan that woke it up.
func (f *FilesystemAdapter) WaitForMessages(ctx context.Context) ([]contracts.Message, error) {
	for {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		current, err := f.listFiles()
		if err != nil {
			return nil, err
		}

		oldCopy := f.snapshotKnownFiles()

		events := diffFiles(oldCopy, current)

		if len(events) > 0 {
			return eventsToMessages(events, current, f), nil
		}

		if err := f.waitIntervalOrDone(ctx); err != nil {
			return nil, err
		}
	}
}

// snapshotKnownFiles returns a copy of f.knownFiles under lock.
func (f *FilesystemAdapter) snapshotKnownFiles() map[string]os.FileInfo {
	f.mu.Lock()
	defer f.mu.Unlock()
	copy := make(map[string]os.FileInfo, len(f.knownFiles))
	for k, v := range f.knownFiles {
		copy[k] = v
	}
	return copy
}

// replaceKnownFiles replaces the knownFiles map under lock.
func (f *FilesystemAdapter) replaceKnownFiles(current map[string]os.FileInfo) {
	if current == nil {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.knownFiles = current
}

// diffFiles produces filesystem events comparing old and current maps.
func diffFiles(old, current map[string]os.FileInfo) []contracts.Event {
	events := make([]contracts.Event, 0)
	for p, info := range current {
		if oldInfo, ok := old[p]; !ok {
			events = append(events, contracts.Event{Source: contracts.SourceFilesystem, Type: "file.created", Timestamp: time.Now(), Metadata: map[string]any{"path": p}, Payload: map[string]any{"size": info.Size()}})
		} else if info.ModTime().After(oldInfo.ModTime()) || info.Size() != oldInfo.Size() {
			events = append(events, contracts.Event{Source: contracts.SourceFilesystem, Type: "file.modified", Timestamp: time.Now(), Metadata: map[string]any{"path": p}, Payload: map[string]any{"size": info.Size()}})
		}
	}
	for p := range old {
		if _, ok := current[p]; !ok {
			events = append(events, contracts.Event{Source: contracts.SourceFilesystem, Type: "file.deleted", Timestamp: time.Now(), Metadata: map[string]any{"path": p}})
		}
	}
	return events
}

// eventsToMessages converts events to Message envelopes.
func eventsToMessages(evts []contracts.Event, snapshot map[string]os.FileInfo, adapter *FilesystemAdapter) []contracts.Message {
	msgs := make([]contracts.Message, 0, len(evts))
	for _, e := range evts {
		msgs = append(msgs, &fsMessage{evt: e, adapter: adapter, snapshot: snapshot})
	}
	return msgs
}

// waitIntervalOrDone waits for the adapter interval or returns if context done.
func (f *FilesystemAdapter) waitIntervalOrDone(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(f.interval):
		return nil
	}
}
