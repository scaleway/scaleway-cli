package human

import (
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

// RenderGrid renders a grid of cells as borderless tables with a
// colPadding-wide gutter between columns.
//
// A multi-line cell ends the current block: its continuation lines are
// flushed as raw lines at column 0 and column widths restart from scratch
// after them (matching the old tabwriter behaviour, kept so wide multi-line
// tables like `k8s version list` stay compact).
//
// ponytail: short columns render at their natural width (no minimum width);
// ragged rows render with slightly different padding than the old tabwriter.
func RenderGrid(grid [][]string) string {
	if len(grid) == 0 {
		return ""
	}

	var sb strings.Builder
	start := 0
	for i := range grid {
		if !hasMultilineCell(grid[i]) && i < len(grid)-1 {
			continue
		}
		if sb.Len() > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(renderBlock(grid[start : i+1]))
		for _, c := range continuations(grid[i]) {
			sb.WriteString("\n")
			sb.WriteString(c)
		}
		start = i + 1
	}

	return sb.String()
}

func hasMultilineCell(row []string) bool {
	for _, c := range row {
		if strings.Contains(c, "\n") {
			return true
		}
	}

	return false
}

// continuations returns the lines following the first line of a multi-line
// cell in the row.
func continuations(row []string) []string {
	var out []string
	for _, c := range row {
		if _, after, ok := strings.Cut(c, "\n"); ok {
			out = append(out, strings.Split(after, "\n")...)
		}
	}

	return out
}

func renderBlock(block [][]string) string {
	ncols := 0
	for _, row := range block {
		ncols = max(ncols, len(row))
	}

	if ncols <= 1 {
		// A single-column line is a tabwriter "last cell": printed raw.
		lines := make([]string, len(block))
		for i, row := range block {
			if len(row) > 0 {
				lines[i] = firstLine(row[0])
			}
		}

		return strings.Join(lines, "\n")
	}

	// Collapse multi-line cells to their first line; the caller flushes the
	// rest.
	cells := make([][]string, len(block))
	for i, row := range block {
		cells[i] = make([]string, ncols)
		for j, c := range row {
			cells[i][j] = firstLine(c)
		}
	}

	t := table.New().
		Border(lipgloss.Border{}).
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false).
		BorderHeader(false).
		BorderColumn(false).
		BorderRow(false).
		Wrap(true).
		// The gutter is the right padding of the non-last columns; the last
		// column is unpadded and its trailing spaces are trimmed below.
		StyleFunc(func(_, col int) lipgloss.Style {
			if col == ncols-1 {
				return lipgloss.NewStyle()
			}

			return lipgloss.NewStyle().PaddingRight(colPadding)
		}).
		Headers(cells[0]...)

	if len(cells) > 1 {
		t = t.Rows(cells[1:]...)
	}

	lines := strings.Split(t.String(), "\n")
	for i, l := range lines {
		lines[i] = strings.TrimRight(l, " ")
	}

	return strings.Join(lines, "\n")
}

func firstLine(s string) string {
	before, _, _ := strings.Cut(s, "\n")

	return before
}
