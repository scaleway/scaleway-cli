package k8s_test

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ClusterConnect(t *testing.T) {
	t.Run("k9s_not_installed", core.Test(&core.TestConfig{
		Commands: k8s.GetCommands(),
		Cmd:      "scw k8s cluster connect 11111111-1111-1111-1111-111111111111",
		OverrideExec: func(_ *core.ExecFuncCtx, cmd *exec.Cmd) (int, error) {
			return 0, &exec.Error{Name: cmd.Args[0], Err: exec.ErrNotFound}
		},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckStderrContains(`Exec: "k9s": executable file not found`),
		),
	}))

	t.Run("simple", core.Test(&core.TestConfig{
		Commands: k8s.GetCommands(),
		Cmd:      "scw k8s cluster connect 3313184d-3767-4535-902a-2bba7b7a4414",
		OverrideExec: func(ctx *core.ExecFuncCtx, cmd *exec.Cmd) (int, error) {
			switch {
			case strings.HasPrefix(cmd.Args[0], "k9s") && len(cmd.Args) > 1 && cmd.Args[1] == "version":
				return 0, nil
			case strings.HasPrefix(cmd.Args[0], "k9s"):
				// Ensure the temporary kubeconfig is passed to k9s
				for _, e := range cmd.Env {
					if strings.HasPrefix(e, "KUBECONFIG=") && strings.HasSuffix(e, ".yaml") {
						return 0, nil
					}
				}

				return 0, errors.New("KUBECONFIG env not set for k9s")
			default:
				return 0, errors.New("unexpected exec: " + cmd.Args[0])
			}
		},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
		),
	}))

	t.Run("template", core.Test(&core.TestConfig{
		Commands: k8s.GetCommands(),
		// The template braces are escaped so the test framework's own templating leaves them for the CLI
		Cmd: `scw k8s cluster connect 3313184d-3767-4535-902a-2bba7b7a4414 ` +
			`template="kubectl --kubeconfig {{"{{"}}.Kubeconfig{{"}}"}} get pods"`,
		OverrideExec: func(ctx *core.ExecFuncCtx, cmd *exec.Cmd) (int, error) {
			// The rendered command must be executed directly, no k9s involved
			if len(cmd.Args) != 5 || cmd.Args[0] != "kubectl" || cmd.Args[1] != "--kubeconfig" ||
				cmd.Args[3] != "get" || cmd.Args[4] != "pods" || !strings.HasSuffix(cmd.Args[2], ".yaml") {
				return 0, errors.New("unexpected rendered command: " + strings.Join(cmd.Args, " "))
			}

			if !slices.Contains(cmd.Env, "KUBECONFIG="+cmd.Args[2]) {
				return 0, errors.New("KUBECONFIG env not set for rendered command")
			}

			return 0, nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
		),
	}))

	t.Run("interrupted", core.Test(&core.TestConfig{
		Commands: k8s.GetCommands(),
		// No OverrideExec: the child process is really executed. It sends SIGINT
		// to scw (its parent) one second in, simulating a Ctrl-C during a
		// long-running session, then keeps running until terminated.
		Cmd: `scw k8s cluster connect 3313184d-3767-4535-902a-2bba7b7a4414 ` +
			`template="sh -c 'sleep 1; kill -INT $PPID; exec sleep 300'"`,
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(130),
			func(t *testing.T, _ *core.CheckFuncCtx) {
				t.Helper()

				pattern := filepath.Join(os.TempDir(), "scw-k8s-kubeconfig-*.yaml")
				leftovers, err := filepath.Glob(pattern)
				require.NoError(t, err)

				assert.Empty(t, leftovers, "temporary kubeconfig must be removed")
			},
		),
	}))
}
