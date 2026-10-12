package tui

import (
	"regexp"
	"strings"
	"testing"
)

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*[A-Za-z]`)

// TestViewHeader verifies the header layout: project/organization/
// zone/region on the left, the current screen's shortcuts in the
// middle, and the Scaleway logo on the right.
func TestViewHeader(t *testing.T) {
	a := &App{
		loc:         NewLocality("fr-par-1"),
		project:     "11111111-aaaa-bbbb-cccc-000000000000",
		projectName: "my-project",
		orgName:     "my-org",
		orgAlias:    "myorg",
		profile:     "prod",
		resources:   map[string]*Resource{},
		stack:       []Screen{NewBrowser(&Resource{Name: "instance server"}, nil)},
		width:       120,
	}
	out := a.viewHeader()
	plain := ansiRe.ReplaceAllString(out, "")

	for _, want := range []string{
		"profile: prod",
		"project:", "my-project (11111111)",
		"organization:", "my-org/myorg",
		"zone: fr-par-1", "region: fr-par",
		// center: the current screen's shortcuts
		"[ enter ] open",
		// right: the logo
		"#############",
	} {
		if !strings.Contains(plain, want) {
			t.Errorf("header missing %q:\n%s", want, plain)
		}
	}
	// The info block starts on the first line, top-left.
	if !strings.HasPrefix(plain, "profile:") {
		t.Errorf("header info should start on the first line:\n%s", plain)
	}
	// The header has one line per logo row.
	if got := strings.Count(plain, "\n") + 1; got != len(scalewayLogo) {
		t.Errorf("header lines = %d, want %d", got, len(scalewayLogo))
	}
}

// TestViewHeaderNoNames verifies the header degrades to short IDs when
// the project/org names have not loaded (or the API failed).
func TestViewHeaderNoNames(t *testing.T) {
	a := &App{
		loc:       NewLocality("nl-ams-1"),
		project:   "22222222-aaaa-bbbb-cccc-000000000000",
		resources: map[string]*Resource{},
		stack:     []Screen{NewBrowser(&Resource{Name: "instance server"}, nil)},
		width:     120,
	}
	plain := ansiRe.ReplaceAllString(a.viewHeader(), "")
	if !strings.Contains(plain, "22222222") {
		t.Errorf("expected the short project id, got:\n%s", plain)
	}
	if !strings.Contains(plain, "region: nl-ams") {
		t.Errorf("expected the region, got:\n%s", plain)
	}
}

// TestLayoutHints verifies hints wrap in at most 3 columns per row,
// with columns padded to the widest hint of the column.
func TestLayoutHints(t *testing.T) {
	got := layoutHints([]string{"a", "bb", "ccc", "dddd", "e"}, 3)
	want := []string{"a     bb  ccc", "dddd  e"}
	if len(got) != len(want) {
		t.Fatalf("rows = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("row %d = %q, want %q", i, got[i], want[i])
		}
	}
}
