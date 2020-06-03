package interactive

import (
	"bytes"
	"context"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/stretchr/testify/require"
)

func TestPromptStringWithConfig(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		buffer := &bytes.Buffer{}

		IsInteractive = false

		SetOutputWriter(buffer)

		ctx := context.Background()
		ctx = InjectMockResponseToContext(ctx, []string{"mock1", "mock2"})

		s, err := PromptStringWithConfig(&PromptStringConfig{
			Ctx:          ctx,
			DefaultValue: "default1",
		})
		require.NoError(t, err)
		assert.Equal(t, "mock1", s)
		s, err = PromptStringWithConfig(&PromptStringConfig{
			Ctx:          ctx,
			DefaultValue: "default2",
		})
		assert.Equal(t, "mock2", s)
		require.NoError(t, err)
	})
}
