package ports_test

import (
	"testing"
	"time"

	"github.com/brunojet/go-infra-backend/internal/config/ports"
)

type mapSource map[string]string

func (m mapSource) Lookup(key string) (string, bool) {
	v, ok := m[key]
	return v, ok
}

func TestTrimmed(t *testing.T) {
	src := mapSource{"FOO": "  bar  "}
	val := ports.Trimmed(src, "FOO")
	if val != "bar" {
		t.Errorf("expected 'bar', got '%s'", val)
	}
}

func TestLowerTrimmed(t *testing.T) {
	src := mapSource{"FOO": "  BAR  "}
	val := ports.LowerTrimmed(src, "FOO")
	if val != "bar" {
		t.Errorf("expected 'bar', got '%s'", val)
	}
}

func TestBool(t *testing.T) {
	src := mapSource{"YES": "true", "NO": "false"}
	if !ports.Bool(src, "YES", false) {
		t.Error("expected true for YES")
	}
	if ports.Bool(src, "NO", true) {
		t.Error("expected false for NO")
	}
	// Test default when key not found
	if !ports.Bool(src, "MISSING", true) {
		t.Error("expected default true for missing key")
	}
	// Test default when value is empty
	src2 := mapSource{"EMPTY": "   "}
	if !ports.Bool(src2, "EMPTY", true) {
		t.Error("expected default true for empty value")
	}
	// Test default for unknown value
	src3 := mapSource{"UNKNOWN": "maybe"}
	if !ports.Bool(src3, "UNKNOWN", true) {
		t.Error("expected default true for unknown value")
	}
	if ports.Bool(src3, "UNKNOWN", false) {
		t.Error("expected default false for unknown value")
	}
}

func TestInt(t *testing.T) {
	src := mapSource{"NUM": "42"}
	if ports.Int(src, "NUM", 0) != 42 {
		t.Error("expected 42 for NUM")
	}
	// Test default when key not found
	if ports.Int(src, "MISSING", 7) != 7 {
		t.Error("expected default 7 for missing key")
	}
	// Test default when value is empty
	src2 := mapSource{"EMPTY": "   "}
	if ports.Int(src2, "EMPTY", 5) != 5 {
		t.Error("expected default 5 for empty value")
	}
	// Test default for invalid int
	src3 := mapSource{"BAD": "abc"}
	if ports.Int(src3, "BAD", 9) != 9 {
		t.Error("expected default 9 for invalid int")
	}
}

func TestDuration(t *testing.T) {
	src := mapSource{"TIME": "2s"}
	val := ports.Duration(src, "TIME", 1*time.Second)
	if val != 2*time.Second {
		t.Errorf("expected 2s, got %s", val)
	}
	// Test default when key not found
	if ports.Duration(src, "MISSING", 3*time.Second) != 3*time.Second {
		t.Error("expected default 3s for missing key")
	}
	// Test default when value is empty
	src2 := mapSource{"EMPTY": "   "}
	if ports.Duration(src2, "EMPTY", 4*time.Second) != 4*time.Second {
		t.Error("expected default 4s for empty value")
	}
	// Test parsing seconds as int
	src3 := mapSource{"SEC": "5"}
	if ports.Duration(src3, "SEC", 1*time.Second) != 5*time.Second {
		t.Error("expected 5s for SEC")
	}
	// Test default for invalid duration
	src4 := mapSource{"BAD": "notaduration"}
	if ports.Duration(src4, "BAD", 6*time.Second) != 6*time.Second {
		t.Error("expected default 6s for invalid duration")
	}
}

func TestSplitCSV(t *testing.T) {
	src := mapSource{"CSV": "a, b, c"}
	vals := ports.SplitCSV(ports.Trimmed(src, "CSV"), nil)
	if len(vals) != 3 || vals[0] != "a" || vals[1] != "b" || vals[2] != "c" {
		t.Errorf("unexpected SplitCSV result: %#v", vals)
	}
	// Test empty input returns default
	def := []string{"x"}
	vals2 := ports.SplitCSV("   ", def)
	if len(vals2) != 1 || vals2[0] != "x" {
		t.Errorf("expected default slice for empty input, got %#v", vals2)
	}
	// Test all empty items returns default
	vals3 := ports.SplitCSV(", , ", def)
	if len(vals3) != 1 || vals3[0] != "x" {
		t.Errorf("expected default slice for all empty items, got %#v", vals3)
	}
}

func TestSplitCSVUpper(t *testing.T) {
	src := mapSource{"CSV": "a, b, c"}
	vals := ports.SplitCSVUpper(ports.Trimmed(src, "CSV"), nil)
	if len(vals) != 3 || vals[0] != "A" || vals[1] != "B" || vals[2] != "C" {
		t.Errorf("unexpected SplitCSVUpper result: %#v", vals)
	}
	// Test default for empty input
	def := []string{"Y"}
	vals2 := ports.SplitCSVUpper("   ", def)
	if len(vals2) != 1 || vals2[0] != "Y" {
		t.Errorf("expected default upper slice for empty input, got %#v", vals2)
	}
}

func TestParseKeyValueCSV(t *testing.T) {
	csv := "k1=v1,k2=v2"
	m := ports.ParseKeyValueCSV(csv)
	if m["k1"] != "v1" || m["k2"] != "v2" {
		t.Errorf("unexpected ParseKeyValueCSV result: %#v", m)
	}
	// Test empty input returns nil
	if ports.ParseKeyValueCSV("   ") != nil {
		t.Error("expected nil for empty input")
	}
	// Test invalid parts are ignored
	m2 := ports.ParseKeyValueCSV("k1=v1,foo,=,k2=,=v3")
	if m2["k1"] != "v1" || m2["k2"] != "" || len(m2) != 2 {
		t.Errorf("unexpected ParseKeyValueCSV with invalid parts: %#v", m2)
	}
	// Test all invalid parts returns nil
	m3 := ports.ParseKeyValueCSV("foo,bar,=,=, ,")
	if m3 != nil {
		t.Errorf("expected nil for all invalid input, got %#v", m3)
	}
}
