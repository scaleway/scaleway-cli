package human

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// RenderGrid renders a grid of cells the way the previous tabwriter-based
// code did:
//   - cells are left-aligned, column widths are driven by the widest cell
//     (ANSI escape sequences are not counted);
//   - each column except the last has a colPadding-wide gutter after it and
//     a minimum width of minwidth (including the gutter);
//   - the last cell of a row is never padded;
//   - a multi-line cell flushes its continuation lines as raw lines at
//     column 0, and column widths restart from scratch after them.
//
// ponytail: assumes uniform rows (every row has the same number of cells),
// which is what all call sites produce. Ragged rows render with slightly
// different padding than the old tabwriter did.
func RenderGrid(grid [][]string, minwidth int) string {
	var sb strings.Builder
	start := 0
	for i := range grid {
		if !hasMultilineCell(grid[i]) && i < len(grid)-1 {
			continue
		}
		if sb.Len() > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(renderBlock(grid[start:i+1], minwidth))
		for _, c := range continuations(grid[i]) {
			sb.WriteString("\n")
			sb.WriteString(c)
		}
		start = i + 1
	}

	return sb.String()
}

func hasMultilineCell(row []string) bool {
	return strings.Contains(strings.Join(row, ""), "\n")
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

// renderBlock renders a block of rows containing no multi-line cell, except
// possibly the last row, whose continuation lines are flushed by the caller.
// Each column except the last is padded with a colPadding-wide gutter to at
// least minwidth; the last cell of a row is never padded.
func renderBlock(block [][]string, minwidth int) string {
	ncols := 0
	for _, row := range block {
		ncols = max(ncols, len(row))
	}

	// Column slots (display width), for every column except the last one:
	// the widest cell plus the gutter, at least minwidth.
	slots := make([]int, max(ncols-1, 0))
	for j := range slots {
		w := minwidth - colPadding
		for _, row := range block {
			if j < len(row) {
				w = max(w, lipgloss.Width(firstLine(row[j])))
			}
		}
		slots[j] = w + colPadding
	}

	var sb strings.Builder
	for i, row := range block {
		if i > 0 {
			sb.WriteString("\n")
		}
		for j, c := range row {
			c = firstLine(c)
			sb.WriteString(c)
			if j < len(row)-1 {
				sb.WriteString(strings.Repeat(" ", slots[j]-lipgloss.Width(c)))
			}
		}
	}

	return sb.String()
}

func firstLine(s string) string {
	before, _, _ := strings.Cut(s, "\n")

	return before
}
