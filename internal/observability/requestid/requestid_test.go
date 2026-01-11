package requestid

import (
	"context"
	"encoding/hex"
	"testing"
)

func TestNew_GeneratesUniqueID(t *testing.T) {
	id1 := New()
	id2 := New()
	if id1 == "" || id2 == "" {
		t.Error("expected non-empty id")
	}
	if id1 == id2 {
		t.Error("expected unique ids")
	}
	_, err := hex.DecodeString(id1)
	if err != nil {
		t.Errorf("id1 is not hex: %v", err)
	}
}

func TestWithAndFrom(t *testing.T) {
	ctx := context.Background()
	id := "test-id"
	ctx2 := With(ctx, id)
	val, ok := From(ctx2)
	if !ok {
		t.Error("expected to extract id from context")
	}
	if val != id {
		t.Errorf("expected id %q, got %q", id, val)
	}
	// Test From returns false for missing or empty
	ctx3 := context.Background()
	_, ok = From(ctx3)
	if ok {
		t.Error("expected false for missing id")
	}
	ctx4 := With(ctx, "")
	_, ok = From(ctx4)
	if ok {
		t.Error("expected false for empty id")
	}
}
