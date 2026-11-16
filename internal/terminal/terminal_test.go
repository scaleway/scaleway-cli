//go:build !wasm

package terminal_test

import (
	"os"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetWidthNoTTY(t *testing.T) {
	old := os.Stderr
	t.Cleanup(func() { os.Stderr = old })

	// Point stderr at a regular file so term.GetSize fails as it does in CI.
	f, err := os.CreateTemp(t.TempDir(), "stderr")
	require.NoError(t, err)
	os.Stderr = f

	// GetWidth must never be negative: -1 fed to strings.Repeat panics.
	assert.GreaterOrEqual(t, terminal.GetWidth(), 0)
}
