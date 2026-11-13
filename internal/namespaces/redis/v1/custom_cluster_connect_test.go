package redis_test

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/redis/v1"
	redisSDK "github.com/scaleway/scaleway-sdk-go/api/redis/v1"
	"github.com/stretchr/testify/assert"
)

// connectInterruptedClusterID is a syntactically valid uuid. The interrupted
// subtest references no template field resolved from the API, so it never
// calls the API and does not need a real cluster (it also passes in cassette
// update mode).
const connectInterruptedClusterID = "66666666-7777-8888-9999-000000000000"

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

// createConnectCluster creates a real cluster for a connect subtest and
// stores it in the test meta as "Cluster". TLS is enabled explicitly so the
// default template exercises the certificate download path. In cassette
// replay mode it replays the recorded create; in cassette update mode
// (nightly) it creates a real cluster, so the subtest stays re-recordable.
func createConnectCluster(name string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd("Cluster", fmt.Sprintf(baseCommand+"tls-enabled=true", name))
}

func deleteConnectCluster() core.AfterFunc {
	return core.ExecAfterCmd("scw redis cluster delete {{ .Cluster.ID }} --wait")
}

// checkNoCertLeftover asserts the temporary TLS certificate of the given
// cluster has been removed from os.TempDir.
func checkNoCertLeftover(t *testing.T, clusterID string) {
	t.Helper()

	_, err := os.Stat(filepath.Join(os.TempDir(), "redis-cert-"+clusterID+".crt"))
	assert.True(t, os.IsNotExist(err), "temporary certificate must be removed")
}

func Test_ClusterConnect(t *testing.T) {
	t.Run("not-installed", core.Test(&core.TestConfig{
		Commands:   redis.GetCommands(),
		BeforeFunc: createConnectCluster("connect-test-not-installed"),
		// No OverrideExec: the command really tries to exec a tool that does
		// not exist, and the child process must fail to start.
		Cmd: "scw redis cluster connect {{ .Cluster.ID }}" +
			` template="not-a-real-tool-xyz --host {{"{{"}}.Endpoint.Host{{"}}"}}"`,
		// The password prompt always fires, even with a custom template.
		PromptResponseMocks: []string{"P@sSw0Rd"},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckStderrContains(`Exec: "not-a-real-tool-xyz": executable file not found`),
		),
		AfterFunc: deleteConnectCluster(),
	}))

	t.Run("simple", core.Test(&core.TestConfig{
		Commands:   redis.GetCommands(),
		BeforeFunc: createConnectCluster("connect-test-simple"),
		// The default template must build the redis-cli invocation: host,
		// port, --tls, --cacert and --user. The password is not in the
		// command: it must be exported as REDISCLI_AUTH.
		// No password=: the user is prompted.
		Cmd:                 "scw redis cluster connect {{ .Cluster.ID }}",
		PromptResponseMocks: []string{"P@sSw0Rd"},
		OverrideExec: func(_ *core.ExecFuncCtx, cmd *exec.Cmd) (int, error) {
			if len(cmd.Args) != 10 {
				return 0, fmt.Errorf("unexpected args: %v", cmd.Args)
			}
			if cmd.Args[0] != "redis-cli" ||
				cmd.Args[1] != "-h" ||
				net.ParseIP(cmd.Args[2]) == nil ||
				cmd.Args[3] != "-p" ||
				cmd.Args[4] != "6379" ||
				cmd.Args[5] != "--tls" ||
				cmd.Args[6] != "--cacert" {
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
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				checkNoCertLeftover(t, ctx.Meta["Cluster"].(*redisSDK.Cluster).ID)
			},
		),
		AfterFunc: deleteConnectCluster(),
	}))

	t.Run("password-arg", core.Test(&core.TestConfig{
		Commands:   redis.GetCommands(),
		BeforeFunc: createConnectCluster("connect-test-password-arg"),
		Cmd:        "scw redis cluster connect {{ .Cluster.ID }} password=P@sSw0Rd",
		// The user must not be prompted when password= is provided: if the
		// prompt fires it returns this trap value and the args assertion
		// below fails.
		PromptResponseMocks: []string{"should-not-be-prompted"},
		OverrideExec: func(_ *core.ExecFuncCtx, cmd *exec.Cmd) (int, error) {
			if len(cmd.Args) != 10 {
				return 0, fmt.Errorf("unexpected args: %v", cmd.Args)
			}
			if cmd.Args[0] != "redis-cli" || cmd.Args[1] != "-h" ||
				net.ParseIP(cmd.Args[2]) == nil ||
				cmd.Args[3] != "-p" ||
				cmd.Args[4] != "6379" ||
				cmd.Args[5] != "--tls" ||
				cmd.Args[6] != "--cacert" {
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
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				checkNoCertLeftover(t, ctx.Meta["Cluster"].(*redisSDK.Cluster).ID)
			},
		),
		AfterFunc: deleteConnectCluster(),
	}))

	t.Run("template", core.Test(&core.TestConfig{
		Commands:   redis.GetCommands(),
		BeforeFunc: createConnectCluster("connect-test-template"),
		// The template braces are escaped so the test framework's own templating leaves them for the CLI
		Cmd: "scw redis cluster connect {{ .Cluster.ID }}" +
			` template="sleep {{"{{"}}.Endpoint.Port{{"}}"}}"`,
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
		AfterFunc: deleteConnectCluster(),
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
			func(t *testing.T, _ *core.CheckFuncCtx) {
				t.Helper()

				checkNoCertLeftover(t, connectInterruptedClusterID)
			},
		),
	}))
}
