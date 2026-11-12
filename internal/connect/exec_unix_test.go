//go:build unix

package connect_test

import (
	"context"
	"io"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/scaleway/scaleway-cli/v2/internal/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ExecCmd_interrupted(t *testing.T) {
	cmd := exec.Command("sleep", "300")
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard

	res := make(chan struct {
		code int
		err  error
	}, 1)
	go func() {
		code, err := connect.ExecCmd(context.Background(), cmd)
		res <- struct {
			code int
			err  error
		}{code, err}
	}()

	// Give the child time to start. The signal handler is registered before
	// the child starts, so there is no window where the test process would
	// use the default (terminate) handler.
	time.Sleep(200 * time.Millisecond)
	require.NoError(t, syscall.Kill(os.Getpid(), syscall.SIGINT))

	select {
	case r := <-res:
		require.NoError(t, r.err)
		assert.Equal(t, 130, r.code, "child killed by SIGINT must exit with 130")
		assert.False(t, cmd.ProcessState.Success())
	case <-time.After(10 * time.Second):
		t.Fatal("ExecCmd did not return after SIGINT")
	}
}
