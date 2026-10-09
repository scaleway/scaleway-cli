package docgen_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/docgen"
	"github.com/stretchr/testify/assert"
)

func TestEscapeTableCell(t *testing.T) {
	assert.Equal(
		t,
		`(new \| ipv4 \| ipv6 \| both)`,
		docgen.EscapeTableCell("(new | ipv4 | ipv6 | both)"),
	)
}
