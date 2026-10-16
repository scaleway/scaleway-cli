package k8s_test

import (
	"errors"
	"os/exec"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1"
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
			core.TestCheckStderrContains("is not installed. Please install k9s"),
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
}
