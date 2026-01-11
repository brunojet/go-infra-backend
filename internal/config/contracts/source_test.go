package contracts_test

import (
	"testing"
)

type mapSource map[string]string

func (m mapSource) Lookup(key string) (string, bool) {
	v, ok := m[key]
	return v, ok
}

func TestMapSource_Lookup(t *testing.T) {
	src := mapSource{"FOO": "bar"}
	val, ok := src.Lookup("FOO")
	if !ok || val != "bar" {
		t.Errorf("expected 'bar', got '%v' (ok=%v)", val, ok)
	}
	_, ok = src.Lookup("NOPE")
	if ok {
		t.Error("expected false for missing key")
	}
}
