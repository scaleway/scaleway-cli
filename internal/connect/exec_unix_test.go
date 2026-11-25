//go:build unix

package connect_test

import (
	"context"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/scaleway/scaleway-cli/v2/internal/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type subprocessStarted chan struct{}

func (ready subprocessStarted) Read([]byte) (int, error) {
	close(ready)

	return 0, io.EOF
}

func Test_ExecCmd_interrupted(t *testing.T) {
	// Keep the test process alive if the runner fails to register its handler.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signals)

	for _, sig := range []syscall.Signal{syscall.SIGINT, syscall.SIGTERM} {
		t.Run(sig.String(), func(t *testing.T) {
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()

			type result struct {
				code int
				err  error
			}
			const children = 2
			results := make(chan result, children)

			for range children {
				ready := make(subprocessStarted)
				cmd := exec.CommandContext(ctx, "sleep", "30")
				// os/exec starts copying stdin only after starting the child.
				cmd.Stdin = ready
				go func() {
					code, err := connect.ExecCmd(ctx, cmd)
					results <- result{code: code, err: err}
				}()

				select {
				case <-ready:
				case <-ctx.Done():
					t.Fatal("child did not start")
				}
			}

			require.NoError(t, syscall.Kill(os.Getpid(), sig))

			for range children {
				select {
				case result := <-results:
					require.NoError(t, result.err)
					assert.Equal(t, 128+int(sig), result.code)
				case <-ctx.Done():
					t.Fatal("child did not exit after signal")
				}
			}
		})
	}
}
