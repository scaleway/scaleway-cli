package k8s_test

import (
	"errors"
	"os/exec"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1"
)

func Test_ClusterKubectl(t *testing.T) {
	t.Run("missing_cluster_id", core.Test(&core.TestConfig{
		Commands: k8s.GetCommands(),
		Cmd:      "scw k8s cluster kubectl",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckStderrContains("A cluster ID is required for this command"),
		),
	}))

	t.Run("missing_kubectl_args", core.Test(&core.TestConfig{
		Commands: k8s.GetCommands(),
		Cmd:      "scw k8s cluster kubectl 3313184d-3767-4535-902a-2bba7b7a4414",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckStderrContains(
				"At least one kubectl argument is required after the cluster ID",
			),
		),
	}))

	t.Run("kubectl_not_installed", core.Test(&core.TestConfig{
		Commands: k8s.GetCommands(),
		Cmd:      "scw k8s cluster kubectl 3313184d-3767-4535-902a-2bba7b7a4414 get nodes",
		OverrideExec: func(_ *core.ExecFuncCtx, cmd *exec.Cmd) (int, error) {
			return 0, &exec.Error{Name: cmd.Args[0], Err: exec.ErrNotFound}
		},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckStderrContains(`Exec: "kubectl": executable file not found`),
		),
	}))

	t.Run("simple", core.Test(&core.TestConfig{
		Commands: k8s.GetCommands(),
		Cmd:      "scw k8s cluster kubectl 3313184d-3767-4535-902a-2bba7b7a4414 get nodes",
		OverrideExec: func(_ *core.ExecFuncCtx, cmd *exec.Cmd) (int, error) {
			if strings.Join(cmd.Args, " ") != "kubectl get nodes" {
				return 0, errors.New("unexpected kubectl command: " + strings.Join(cmd.Args, " "))
			}

			// Ensure the temporary kubeconfig is passed to kubectl
			for _, e := range cmd.Env {
				if strings.HasPrefix(e, "KUBECONFIG=") && strings.HasSuffix(e, ".yaml") {
					return 0, nil
				}
			}

			return 0, errors.New("KUBECONFIG env not set for kubectl")
		},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
		),
	}))

	t.Run("flags_with_separator", core.Test(&core.TestConfig{
		Commands: k8s.GetCommands(),
		Cmd:      "scw k8s cluster kubectl 3313184d-3767-4535-902a-2bba7b7a4414 -- get pods -o wide -l app=web",
		OverrideExec: func(_ *core.ExecFuncCtx, cmd *exec.Cmd) (int, error) {
			if strings.Join(cmd.Args, " ") != "kubectl get pods -o wide -l app=web" {
				return 0, errors.New("unexpected kubectl command: " + strings.Join(cmd.Args, " "))
			}

			return 0, nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
		),
	}))

	t.Run("exit_code", core.Test(&core.TestConfig{
		Commands:     k8s.GetCommands(),
		Cmd:          "scw k8s cluster kubectl 3313184d-3767-4535-902a-2bba7b7a4414 get pods",
		OverrideExec: core.OverrideExecSimple("kubectl get pods", 3),
		Check:        core.TestCheckExitCode(3),
	}))
}
