package events_test

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	eventsPkg "github.com/brunojet/go-infra-backend/internal/events"
	adapters "github.com/brunojet/go-infra-backend/internal/events/adapters"
	contracts "github.com/brunojet/go-infra-backend/pkg/events/contracts"
)

func TestEventManager_FilesystemAdapter_EndToEnd(t *testing.T) {
	tmp := t.TempDir()

	adapter := adapters.NewFilesystemAdapter(tmp, 50*time.Millisecond)
	mgr := eventsPkg.NewManager(adapter)

	var mu sync.Mutex
	var got []contracts.Event

	handler := func(ctx context.Context, e contracts.Event) error {
		mu.Lock()
		got = append(got, e)
		mu.Unlock()
		return nil
	}

	// register handler specifically for this adapter
	mgr.RegisterHandlerForAdapter(adapter, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := mgr.Start(ctx); err != nil {
		t.Fatalf("start manager: %v", err)
	}

	// create a file that should trigger a file.created event
	fpath := filepath.Join(tmp, "foo.txt")
	if err := os.WriteFile(fpath, []byte("hello"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	// wait until we receive at least one event or timeout
	deadline := time.After(3 * time.Second)
	tick := time.Tick(30 * time.Millisecond)
	received := false
	for !received {
		select {
		case <-deadline:
			t.Fatal("timeout waiting for filesystem event")
		case <-tick:
			mu.Lock()
			if len(got) > 0 {
				received = true
			}
			mu.Unlock()
		}
	}

	// shutdown and assert
	if err := mgr.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(got) == 0 {
		t.Fatalf("expected events, got none")
	}
	// ensure at least one created event references our file
	found := false
	for _, e := range got {
		if e.Type != "file.created" {
			continue
		}
		if mp, ok := e.Metadata["path"].(string); ok {
			if mp == "foo.txt" || mp == filepath.Join(".", "foo.txt") || mp == "./foo.txt" || filepath.Base(mp) == "foo.txt" {
				found = true
				break
			}
		}
	}
	if !found {
		t.Fatalf("did not observe file.created for foo.txt; events: %+v", got)
	}
}
