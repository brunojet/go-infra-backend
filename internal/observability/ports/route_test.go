package ports

import "testing"

func TestSelectRoute_PrefersTemplate(t *testing.T) {
	template := "/apps/{id}"
	actual := "/apps/123"

	got := SelectRoute(template, actual)
	if got != template {
		t.Fatalf("SelectRoute() = %q, want %q", got, template)
	}
}

func TestSelectRoute_UnmatchedWhenTemplateMissing(t *testing.T) {
	cases := []struct {
		name     string
		template string
	}{
		{name: "empty", template: ""},
		{name: "whitespace", template: "  \t\n"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := SelectRoute(tc.template, "/apps/123")
			if got != RouteUnmatched {
				t.Fatalf("SelectRoute() = %q, want %q", got, RouteUnmatched)
			}
		})
	}
}
