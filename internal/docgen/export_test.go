package docgen

// Exported for black-box tests of the table cell escaping.

func EscapeTableCell(s string) string {
	return escapeTableCell(s)
}
