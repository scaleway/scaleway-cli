package redis_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/redis/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// One cluster id per subtest so the temporary certificate files
// (redis-cert-<cluster-id>.crt in os.TempDir) never clash between
// parallel subtests.
const (
	connectSimpleClusterID       = "047ff1eb-5215-4bdf-bec3-26c11766bea9"
	connectTemplateClusterID     = "11111111-2222-3333-4444-555555555555"
	connectInterruptedClusterID  = "66666666-7777-8888-9999-000000000000"
	connectNotInstalledClusterID = "aabbccdd-eeff-0011-2233-445566778899"
	connectPasswordArgClusterID  = "99999999-8888-7777-6666-555555555555"
)

// envRedisAuth returns the value of REDISCLI_AUTH in the child process
// environment, or "" if it is not set.
func envRedisAuth(cmd *exec.Cmd) string {
	for _, e := range cmd.Env {
		if v, ok := strings.CutPrefix(e, "REDISCLI_AUTH="); ok {
			return v
		}
	}

	return ""
}

func Test_ClusterConnect(t *testing.T) {
	checkNoCertLeftover := func(clusterID string) func(t *testing.T, _ *core.CheckFuncCtx) {
		return func(t *testing.T, _ *core.CheckFuncCtx) {
			t.Helper()

			leftovers, err := filepath.Glob(
				filepath.Join(os.TempDir(), "redis-cert-"+clusterID+".crt"),
			)
			require.NoError(t, err)

			assert.Empty(t, leftovers, "temporary certificate must be removed")
		}
	}

	t.Run("not-installed", core.Test(&core.TestConfig{
		Commands: redis.GetCommands(),
		// No OverrideExec: the command really tries to exec a tool that does
		// not exist, and the child process must fail to start.
		Cmd: "scw redis cluster connect " + connectNotInstalledClusterID +
			` template="not-a-real-tool-xyz --host {{"{{"}}.Host{{"}}"}}"`,
		// The password prompt always fires, even with a custom template.
		PromptResponseMocks: []string{"P@sSw0Rd"},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckStderrContains(`Exec: "not-a-real-tool-xyz": executable file not found`),
		),
	}))

	t.Run("simple", core.Test(&core.TestConfig{
		Commands: redis.GetCommands(),
		Cmd:      "scw redis cluster connect " + connectSimpleClusterID,
		// The default template must build the redis-cli invocation: host,
		// port, --tls, --cacert and --user. The password is not in the
		// command: it must be exported as REDISCLI_AUTH.
		// No password=: the user is prompted.
		PromptResponseMocks: []string{"P@sSw0Rd"},
		OverrideExec: func(_ *core.ExecFuncCtx, cmd *exec.Cmd) (int, error) {
			if len(cmd.Args) != 10 {
				return 0, fmt.Errorf("unexpected args: %v", cmd.Args)
			}
			want := []string{
				"redis-cli", "-h", "78.232.49.55", "-p", "6379",
				"--tls", "--cacert",
			}
			if !slices.Equal(cmd.Args[:7], want) {
				return 0, fmt.Errorf("unexpected args: %v", cmd.Args)
			}
			if _, err := os.Stat(cmd.Args[7]); err != nil {
				return 0, fmt.Errorf("certificate file not found: %w", err)
			}
			if cmd.Args[8] != "--user" || cmd.Args[9] != "admin" {
				return 0, fmt.Errorf("unexpected args: %v", cmd.Args)
			}
			if got := envRedisAuth(cmd); got != "P@sSw0Rd" {
				return 0, fmt.Errorf("unexpected REDISCLI_AUTH: %q", got)
			}

			return 0, nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			checkNoCertLeftover(connectSimpleClusterID),
		),
	}))

	t.Run("password-arg", core.Test(&core.TestConfig{
		Commands: redis.GetCommands(),
		Cmd:      "scw redis cluster connect " + connectPasswordArgClusterID + " password=P@sSw0Rd",
		// The user must not be prompted when password= is provided: if the
		// prompt fires it returns this trap value and the args assertion
		// below fails.
		PromptResponseMocks: []string{"should-not-be-prompted"},
		OverrideExec: func(_ *core.ExecFuncCtx, cmd *exec.Cmd) (int, error) {
			if len(cmd.Args) != 10 {
				return 0, fmt.Errorf("unexpected args: %v", cmd.Args)
			}
			want := []string{
				"redis-cli", "-h", "78.232.49.55", "-p", "6379",
				"--tls", "--cacert",
			}
			if !slices.Equal(cmd.Args[:7], want) {
				return 0, fmt.Errorf("unexpected args: %v", cmd.Args)
			}
			if cmd.Args[8] != "--user" || cmd.Args[9] != "admin" {
				return 0, fmt.Errorf("unexpected args: %v", cmd.Args)
			}
			// The prompted trap value must not leak into REDISCLI_AUTH.
			if got := envRedisAuth(cmd); got != "P@sSw0Rd" {
				return 0, fmt.Errorf("unexpected REDISCLI_AUTH: %q", got)
			}

			return 0, nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			checkNoCertLeftover(connectPasswordArgClusterID),
		),
	}))

	t.Run("template", core.Test(&core.TestConfig{
		Commands: redis.GetCommands(),
		// The template braces are escaped so the test framework's own templating leaves them for the CLI
		Cmd: "scw redis cluster connect " + connectTemplateClusterID +
			` template="sleep {{"{{"}}.Port{{"}}"}}"`,
		// The password prompt always fires, even with a custom template.
		PromptResponseMocks: []string{"P@sSw0Rd"},
		OverrideExec: func(_ *core.ExecFuncCtx, cmd *exec.Cmd) (int, error) {
			want := []string{"sleep", "6379"}
			if !slices.Equal(cmd.Args, want) {
				return 0, fmt.Errorf("unexpected rendered command: %v", cmd.Args)
			}

			return 0, nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
		),
	}))

	t.Run("interrupted", core.Test(&core.TestConfig{
		Commands: redis.GetCommands(),
		// No OverrideExec: the child process is really executed. It sends SIGINT
		// to scw (its parent) one second in, simulating a Ctrl-C during a
		// long-running session, then keeps running until terminated.
		Cmd: "scw redis cluster connect " + connectInterruptedClusterID +
			` template="sh -c 'sleep 1; kill -INT $PPID; exec sleep 300'"`,
		// The password prompt always fires, even with a custom template.
		PromptResponseMocks: []string{"P@sSw0Rd"},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(130),
			checkNoCertLeftover(connectInterruptedClusterID),
		),
	}))
}
