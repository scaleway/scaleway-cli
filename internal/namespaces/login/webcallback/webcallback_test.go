package webcallback_test

import (
	"context"
	"testing"
	"time"

	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/login/webcallback"
	"github.com/stretchr/testify/assert"
)

func TestWebCallback(t *testing.T) {
	wb := webcallback.New()

	t.Cleanup(func() {
		wb.Close()
	})
	assert.NoError(t, wb.Start())
	assert.NoError(t, wb.Trigger("test-token", time.Second))

	ctx, cancelFunc := context.WithTimeout(t.Context(), time.Second)
	t.Cleanup(cancelFunc)

	resp, err := wb.Wait(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "test-token", resp)
}
