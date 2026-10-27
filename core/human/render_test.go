package human_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core/human"
)

func TestRenderGrid(t *testing.T) {
	red := "\x1b[31m"
	green := "\x1b[32m"
	reset := "\x1b[0m"

	cases := []struct {
		name     string
		minwidth int
		grid     [][]string
		want     string
	}{
		{
			name:     "simple",
			minwidth: 5,
			grid:     [][]string{{"ID", "NAME"}, {"123", "prod"}, {"456", "stg"}},
			want:     "ID   NAME\n123  prod\n456  stg",
		},
		{
			name:     "short-column-minwidth",
			minwidth: 5,
			grid:     [][]string{{"ID", "NAME"}, {"a", "production-server-name"}},
			want:     "ID   NAME\na    production-server-name",
		},
		{
			name:     "ansi-cells",
			minwidth: 5,
			grid: [][]string{
				{"NAME", "STATUS"},
				{red + "broken" + reset, green + "active" + reset},
			},
			want: "NAME    STATUS\n" + red + "broken" + reset + "  " + green + "active" + reset,
		},
		{
			name:     "multiline-cell",
			minwidth: 5,
			grid:     [][]string{{"NAME", "DESC"}, {"a", "line one\nline two"}, {"b", "x"}},
			want:     "NAME  DESC\na     line one\nline two\nb    x",
		},
		{
			name:     "empty-cells",
			minwidth: 5,
			grid:     [][]string{{"A", "B"}, {"", ""}, {"x", ""}},
			want:     "A    B\n     \nx    ",
		},
		{
			name:     "single-column",
			minwidth: 5,
			grid:     [][]string{{"A"}, {"value"}, {"b"}},
			want:     "A\nvalue\nb",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := human.RenderGrid(tc.grid, tc.minwidth); got != tc.want {
				t.Errorf("RenderGrid() = %q, want %q", got, tc.want)
			}
		})
	}
}
