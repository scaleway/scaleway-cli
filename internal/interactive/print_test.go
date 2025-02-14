//go:build !wasm

package interactive_test

import (
	"bytes"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrint(t *testing.T) {
	buffer := &bytes.Buffer{}
	interactive.SetOutputWriter(buffer)
	interactive.IsInteractive = true

	n, err := interactive.Print("Test", 42)
	require.NoError(t, err)
	assert.Equal(t, 6, n)

	interactive.IsInteractive = false
	n, err = interactive.Print("Test", 42)
	require.NoError(t, err)
	assert.Equal(t, 0, n)

	assert.Equal(t, "Test42", buffer.String())
}

func TestPrintln(t *testing.T) {
	buffer := &bytes.Buffer{}
	interactive.SetOutputWriter(buffer)
	interactive.IsInteractive = true

	n, err := interactive.Println("Test", 42)
	require.NoError(t, err)
	assert.Equal(t, 8, n)

	interactive.IsInteractive = false
	n, err = interactive.Println("Test", 42)
	require.NoError(t, err)
	assert.Equal(t, 0, n)

	assert.Equal(t, "Test 42\n", buffer.String())
}

func TestPrintf(t *testing.T) {
	buffer := &bytes.Buffer{}
	interactive.SetOutputWriter(buffer)
	interactive.IsInteractive = true

	n, err := interactive.Printf("%s %d", "Test", 42)
	require.NoError(t, err)
	assert.Equal(t, 7, n)

	interactive.IsInteractive = false
	n, err = interactive.Printf("%s %d", "Test", 42)
	require.NoError(t, err)
	assert.Equal(t, 0, n)

	assert.Equal(t, "Test 42", buffer.String())
}
