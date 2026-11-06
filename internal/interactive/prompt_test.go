package interactive_test

import (
	"bytes"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPromptStringWithConfig(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		buffer := &bytes.Buffer{}

		interactive.IsInteractive = false

		interactive.SetOutputWriter(buffer)

		ctx := t.Context()
		ctx = interactive.InjectMockResponseToContext(ctx, []string{"mock1", "mock2"})

		s, err := interactive.PromptString(ctx, "prompt1", "default1", "", nil)
		require.NoError(t, err)
		assert.Equal(t, "mock1", s)
		s, err = interactive.PromptString(ctx, "prompt2", "default2", "", nil)
		assert.Equal(t, "mock2", s)
		require.NoError(t, err)
	})
}
