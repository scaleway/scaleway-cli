package docgen

import (
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
)

// TestRenderNamespaceHeadingSurroundedByBlankLines ensures the generated top
// heading is surrounded by blank lines (markdownlint MD022). See issue #5666.
func TestRenderNamespaceHeadingSurroundedByBlankLines(t *testing.T) {
	ns := &Namespace{
		Cmd: &core.Command{
			Namespace: "testns",
			Short:     "Namespace description.",
		},
		Resources: map[string]*Resource{},
	}

	out, err := renderNamespace(ns)
	if err != nil {
		t.Fatalf("renderNamespace failed: %v", err)
	}

	lines := strings.Split(out, "\n")
	headingIdx := -1
	for i, l := range lines {
		if strings.HasPrefix(l, "# Documentation for ") {
			headingIdx = i

			break
		}
	}
	if headingIdx == -1 {
		t.Fatalf("no top-level heading found in output:\n%s", out)
	}

	// Print the region around the heading so both RED and GREEN runs carry
	// evidence of what was actually generated.
	start := headingIdx - 1
	if start < 0 {
		start = 0
	}
	end := headingIdx + 2
	if end > len(lines) {
		end = len(lines)
	}
	t.Logf("heading region (lines %d-%d):\n%q", start, end-1, lines[start:end])

	if headingIdx == 0 || lines[headingIdx-1] != "" {
		t.Errorf("MD022: heading %q must be preceded by a blank line, got %q", lines[headingIdx], lines[headingIdx-1])
	}
	if headingIdx+1 >= len(lines) || lines[headingIdx+1] != "" {
		t.Errorf("MD022: heading %q must be followed by a blank line, got %q", lines[headingIdx], lines[headingIdx+1])
	}
}
