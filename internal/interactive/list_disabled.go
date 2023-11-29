//go:build wasm

package interactive

import (
	"context"
	"fmt"
)

type ListPrompt struct {
	Prompt       string
	Choices      []string
	DefaultIndex int
}

func (m *ListPrompt) Execute(ctx context.Context) (int, error) {
	return -1, fmt.Errorf("not implemented for current platform")
}
