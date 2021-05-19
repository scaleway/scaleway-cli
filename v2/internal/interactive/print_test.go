package interactive

import (
	"bytes"
	"testing"

	"github.com/alecthomas/assert"
)

func TestPrint(t *testing.T) {
	buffer := &bytes.Buffer{}
	outputWriter = buffer
	IsInteractive = true

	n, err := Print("Test", 42)
	assert.NoError(t, err)
	assert.Equal(t, 6, n)

	IsInteractive = false
	n, err = Print("Test", 42)
	assert.NoError(t, err)
	assert.Equal(t, 0, n)

	assert.Equal(t, "Test42", buffer.String())
}

func TestPrintln(t *testing.T) {
	buffer := &bytes.Buffer{}
	outputWriter = buffer
	IsInteractive = true

	n, err := Println("Test", 42)
	assert.NoError(t, err)
	assert.Equal(t, 8, n)

	IsInteractive = false
	n, err = Println("Test", 42)
	assert.NoError(t, err)
	assert.Equal(t, 0, n)

	assert.Equal(t, "Test 42\n", buffer.String())
}

func TestPrintf(t *testing.T) {
	buffer := &bytes.Buffer{}
	outputWriter = buffer
	IsInteractive = true

	n, err := Printf("%s %d", "Test", 42)
	assert.NoError(t, err)
	assert.Equal(t, 7, n)

	IsInteractive = false
	n, err = Printf("%s %d", "Test", 42)
	assert.NoError(t, err)
	assert.Equal(t, 0, n)

	assert.Equal(t, "Test 42", buffer.String())
}
