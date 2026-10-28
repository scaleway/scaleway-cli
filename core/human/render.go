package human

import (
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
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
	if len(grid) == 0 {
		return ""
	}

	// A multi-line cell ends the current block: tabwriter flushes its
	// continuation lines as raw lines at column 0 and restarts the column
	// widths from scratch after them.
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

func renderBlock(block [][]string, minwidth int) string {
	ncols := 0
	for _, row := range block {
		ncols = max(ncols, len(row))
	}

	if ncols <= 1 {
		// A single-column line is a tabwriter "last cell": printed raw,
		// without padding or minwidth.
		lines := make([]string, len(block))
		for i, row := range block {
			if len(row) > 0 {
				lines[i] = firstLine(row[0])
			}
		}

		return strings.Join(lines, "\n")
	}

	// Collapse multi-line cells (only the last row of a block can have one):
	// the first line participates in the block, the rest are flushed by the
	// caller.
	// raw keeps the unpadded cells: the last cell of a row is printed raw,
	// without the minwidth padding.
	cells := make([][]string, len(block))
	raw := make([][]string, len(block))
	for i, row := range block {
		cells[i] = make([]string, ncols)
		raw[i] = make([]string, len(row))
		for j, c := range row {
			cells[i][j] = firstLine(c)
			raw[i][j] = cells[i][j]
		}
	}

	// Column slots (display width), for every column except the last one:
	// the widest cell plus the gutter, at least minwidth.
	slots := make([]int, ncols)
	for j := range ncols - 1 {
		maxW := 0
		for _, row := range cells {
			if j < len(row) {
				maxW = max(maxW, lipgloss.Width(row[j]))
			}
		}
		if minwidth > 0 && maxW+colPadding < minwidth {
			// lipgloss has no minimum column width, so pad the cells
			// ourselves to reach it.
			// ponytail: pads with plain (unstyled) spaces; a column whose
			// cells are fully colored and shorter than minwidth gets
			// uncolored trailing spaces. Nobody notices; remove if it ever does.
			for r := range cells {
				if j < len(cells[r]) {
					cells[r][j] = padWidth(cells[r][j], minwidth-colPadding)
				}
			}
			maxW = minwidth - colPadding
		}
		slots[j] = maxW + colPadding
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
		// Wrap(true) is required: the table only computes multi-line row
		// heights when wrapping is enabled. Column widths equal the widest
		// cell, so no wrapping actually happens.
		Wrap(true).
		// The gutter is the right padding of the non-last columns, like in
		// tabwriter. The last column is rendered but discarded (see below).
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

	rendered := strings.Split(t.String(), "\n")

	// Rebuild each line: the non-last columns exactly as rendered, then the
	// last cell raw (tabwriter never pads the last cell of a line).
	lines := make([]string, len(block))
	for i, row := range block {
		if len(row) <= 1 {
			lines[i] = raw[i][0]

			continue
		}
		offset := 0
		for j := range len(row) - 1 {
			offset += len(cells[i][j]) + slots[j] - lipgloss.Width(cells[i][j])
		}
		lines[i] = rendered[i][:offset] + raw[i][len(row)-1]
	}

	return strings.Join(lines, "\n")
}

func firstLine(s string) string {
	before, _, _ := strings.Cut(s, "\n")

	return before
}

// padWidth appends spaces so the display width of s reaches at least w.
func padWidth(s string, w int) string {
	d := w - lipgloss.Width(s)
	if d <= 0 {
		return s
	}

	return s + strings.Repeat(" ", d)
}
