package bench

import (
	"database/sql"
	"encoding/json"
	"testing"
)

type ValDTO struct {
	Name string `json:"name,omitempty"`
}

type PtrDTO struct {
	Name *string `json:"name,omitempty"`
}

var test sql.NullString

func Test_OmitEmpty_StringValue(t *testing.T) {
	var v ValDTO
	if err := json.Unmarshal([]byte(`{"name":""}`), &v); err != nil {
		t.Fatal(err)
	}
	if v.Name != "" {
		t.Fatalf("expected empty string, got %q", v.Name)
	}
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "{}" {
		t.Fatalf("expected {}, got %s", b)
	}
}

func Test_OmitEmpty_Pointer(t *testing.T) {
	var p PtrDTO

	// absent field -> nil pointer -> omitted on marshal
	if err := json.Unmarshal([]byte(`{}`), &p); err != nil {
		t.Fatal(err)
	}
	if p.Name != nil {
		t.Fatalf("expected nil, got %v", p.Name)
	}
	b, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "{}" {
		t.Fatalf("expected {}, got %s", b)
	}

	// explicit empty string -> pointer to empty string -> marshals as "name":""
	if err := json.Unmarshal([]byte(`{"name":""}`), &p); err != nil {
		t.Fatal(err)
	}
	if p.Name == nil || *p.Name != "" {
		t.Fatalf("expected ptr to empty string, got %v", p.Name)
	}
	b, err = json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"name":""}` {
		t.Fatalf("expected {\"name\":\"\"}, got %s", b)
	}
}
