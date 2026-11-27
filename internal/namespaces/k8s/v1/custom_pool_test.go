package k8s_test

import (
	"bytes"
	"io"
	"os/exec"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/platform/terminal"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/require"
)

func TestPoolAddExternalNode(t *testing.T) {
	const (
		nodeIP    = "192.0.2.9"
		poolID    = "22222222-2222-4222-8222-222222222222"
		secretKey = "11111111-1111-1111-1111-111111111111"
		script    = "/home/alice/init_kosmos_node.sh"
	)
	client, err := scw.NewClient(scw.WithAuth("SCWXXXXXXXXXXXXXXXXX", secretKey))
	require.NoError(t, err)
	// A command that bypasses the hook must never reach a real SSH client.
	t.Setenv("PATH", t.TempDir())

	for _, tc := range []struct {
		name   string
		failAt int
	}{
		{name: "success"},
		{name: "copy fails", failAt: 1},
		{name: "chmod fails", failAt: 2},
		{name: "execution fails", failAt: 3},
	} {
		t.Run(tc.name, func(t *testing.T) {
			commands := k8s.GetCommands()
			commands.MustFind("k8s", "pool", "add-external-node").DisableAfterChecks = true
			var stderr bytes.Buffer
			var calls [][]string
			code, _, err := core.Bootstrap(t.Context(), &core.BootstrapConfig{
				Args: []string{
					"scw", "k8s", "pool", "add-external-node",
					"node-ip=" + nodeIP, "pool-id=" + poolID, "username=alice", "region=fr-par",
				},
				Commands:         commands,
				Client:           client,
				Stdout:           io.Discard,
				Stderr:           &stderr,
				Platform:         terminal.NewPlatform("test"),
				BuildInfo:        &core.BuildInfo{},
				DisableTelemetry: true,
				OverrideEnv:      map[string]string{"HOME": t.TempDir()},
				OverrideExec: func(cmd *exec.Cmd) (int, error) {
					require.Greater(t, len(cmd.Args), 5)
					require.Equal(t, []string{"ssh", nodeIP, "-t", "-l", "alice"}, cmd.Args[:5])
					calls = append(calls, cmd.Args[5:])
					if len(calls) == tc.failAt {
						return 17, nil
					}

					return 0, nil
				},
			})
			if tc.failAt != 0 {
				require.ErrorContains(t, err, "ssh command failed with exit code 17")
				require.Equal(t, 1, code)
				require.Len(t, calls, tc.failAt)

				return
			}

			require.NoError(t, err)
			require.Zero(t, code)
			require.Len(t, calls, 3)
			require.GreaterOrEqual(t, len(calls[0]), 6)
			require.Equal(t, []string{"cat", "<<", "EOF", ">", script}, calls[0][:5])
			require.Contains(
				t,
				strings.Join(calls[0][5:], " "),
				"POOL_ID="+poolID+" POOL_REGION=fr-par",
			)
			require.Equal(t, []string{"chmod", "+x", script}, calls[1])
			require.Equal(
				t,
				[]string{"SCW_SECRET_KEY=" + secretKey, "./init_kosmos_node.sh"},
				calls[2],
			)
			require.NotContains(t, stderr.String(), secretKey)
		})
	}
}
