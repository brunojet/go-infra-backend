package nethttp

import "testing"

func TestTranslatePath_Empty(t *testing.T) {
	if got := translatePath(""); got != "" {
		t.Fatalf("translatePath(\"\") = %q, want empty", got)
	}
}

func TestTranslatePath_SingleParam(t *testing.T) {
	got := translatePath("/apps/:id")
	want := "/apps/{id}"
	if got != want {
		t.Fatalf("translatePath() = %q, want %q", got, want)
	}
}

func TestTranslatePath_MultipleParams(t *testing.T) {
	got := translatePath("/a/:x/b/:y")
	want := "/a/{x}/b/{y}"
	if got != want {
		t.Fatalf("translatePath() = %q, want %q", got, want)
	}
}

func TestTranslatePath_NoParams_Unchanged(t *testing.T) {
	got := translatePath("/health")
	want := "/health"
	if got != want {
		t.Fatalf("translatePath() = %q, want %q", got, want)
	}
}

func TestTranslatePath_ParamAtRoot(t *testing.T) {
	got := translatePath("/:id")
	want := "/{id}"
	if got != want {
		t.Fatalf("translatePath() = %q, want %q", got, want)
	}
}
