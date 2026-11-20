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
	"time"

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

// The gateway subtests replay pre-recorded cassettes of a connect session;
// the IDs below are the resources recorded in them, so nothing is created
// or deleted in replay.
const (
	connectGatewayClusterID   = "955763a4-77ed-4ce9-9380-7932e3d3592c"
	connectGatewayID          = "a86e6447-2a9f-475f-9ef6-98f34e32b3f3"
	connectGatewayPrivateIP   = "172.16.12.2"
	connectGatewayBastionPort = "61000"
	connectGatewayPublicIP    = "151.115.59.64"
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

	// The gateway subtests put a fake ssh binary (a Go program compiled at
	// test time) on the PATH: it records the command line it is invoked
	// with and listens on the local port it is asked to forward, which is
	// enough for the tunnel to become ready.
	t.Run("gateway", func(t *testing.T) {
		argsFile := setupFakeSSH(t)

		core.Test(&core.TestConfig{
			Commands: redis.GetCommands(),
			Cmd: "scw redis cluster connect " + connectGatewayClusterID +
				" zone=pl-waw-1 vpc-gateway-id=" + connectGatewayID,
			PromptResponseMocks: []string{"P@sSw0Rd"},
			OverrideExec:        gatewayConnectOverride(argsFile),
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				func(t *testing.T, _ *core.CheckFuncCtx) {
					t.Helper()

					checkNoCertLeftover(t, connectGatewayClusterID)
				},
				checkTunnelStopped(argsFile),
			),
		})
	})

	// No vpc-gateway-id: the gateway must be auto-detected from the
	// cluster's private network. The cluster has no public endpoint, so the
	// plain connect must fall back to the private one.
	t.Run("private-network-autodetect", func(t *testing.T) {
		argsFile := setupFakeSSH(t)

		core.Test(&core.TestConfig{
			Commands:            redis.GetCommands(),
			Cmd:                 "scw redis cluster connect " + connectGatewayClusterID + " zone=pl-waw-1",
			PromptResponseMocks: []string{"P@sSw0Rd"},
			OverrideExec:        gatewayConnectOverride(argsFile),
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				func(t *testing.T, _ *core.CheckFuncCtx) {
					t.Helper()

					checkNoCertLeftover(t, connectGatewayClusterID)
				},
				checkTunnelStopped(argsFile),
			),
		})
	})

	t.Run("no-gateway", core.Test(&core.TestConfig{
		Commands: redis.GetCommands(),
		Cmd:      "scw redis cluster connect " + connectGatewayClusterID + " zone=pl-waw-1",
		// The password prompt always fires, even before the gateway error.
		PromptResponseMocks: []string{"P@sSw0Rd"},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckStderrContains(
				"No VPC gateway with the SSH bastion enabled is attached to private network",
			),
		),
	}))

	t.Run("multiple-gateways", core.Test(&core.TestConfig{
		Commands:            redis.GetCommands(),
		Cmd:                 "scw redis cluster connect " + connectGatewayClusterID + " zone=pl-waw-1",
		PromptResponseMocks: []string{"P@sSw0Rd"},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckStderrContains(
				"Multiple VPC gateways with the SSH bastion enabled are attached",
			),
			core.TestCheckStderrContains("Specify one with vpc-gateway-id"),
		),
	}))

	t.Run("gateway-bastion-disabled", core.Test(&core.TestConfig{
		Commands: redis.GetCommands(),
		Cmd: "scw redis cluster connect " + connectGatewayClusterID +
			" zone=pl-waw-1 vpc-gateway-id=" + connectGatewayID,
		PromptResponseMocks: []string{"P@sSw0Rd"},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckStderrContains("the SSH bastion is not enabled on gateway"),
		),
	}))
}

// setupFakeSSH puts the fake ssh binary on the PATH for the duration of the
// test and returns the path of the file where it records its args.
func setupFakeSSH(t *testing.T) string {
	t.Helper()

	fakeDir := t.TempDir()
	writeFakeSSH(t, fakeDir)
	t.Setenv("PATH", fakeDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	argsFile := filepath.Join(fakeDir, "ssh-args")
	t.Setenv("FAKE_SSH_ARGS_FILE", argsFile)

	return argsFile
}

// gatewayConnectOverride asserts the connect command targets the tunnel's
// local port on loopback, that the ssh command line is the expected one,
// and that the certificate and the cluster user are handled as usual.
func gatewayConnectOverride(argsFile string) func(_ *core.ExecFuncCtx, cmd *exec.Cmd) (int, error) {
	return func(_ *core.ExecFuncCtx, cmd *exec.Cmd) (int, error) {
		// The connect command must target the tunnel's local port
		// on loopback, not the private IP.
		if !slices.Equal(cmd.Args[:4], []string{"redis-cli", "-h", "127.0.0.1", "-p"}) {
			return 0, fmt.Errorf("unexpected args: %v", cmd.Args)
		}

		sshArgs, err := readFakeSSHArgs(argsFile)
		if err != nil {
			return 0, err
		}
		want := []string{
			"-N",
			"-o", "ExitOnForwardFailure=yes",
			"-L", "127.0.0.1:" + cmd.Args[4] + ":" + connectGatewayPrivateIP + ":6379",
			"-p", connectGatewayBastionPort,
			"bastion@" + connectGatewayPublicIP,
		}
		if !slices.Equal(sshArgs, want) {
			return 0, fmt.Errorf("unexpected ssh command line: %v", sshArgs)
		}

		// The TLS certificate and the cluster user are handled as usual.
		if !slices.Equal(cmd.Args[5:7], []string{"--tls", "--cacert"}) {
			return 0, fmt.Errorf("unexpected args: %v", cmd.Args)
		}
		if _, err := os.Stat(cmd.Args[7]); err != nil {
			return 0, fmt.Errorf("certificate file not found: %w", err)
		}
		if !slices.Equal(cmd.Args[8:], []string{"--user", "default"}) {
			return 0, fmt.Errorf("unexpected args: %v", cmd.Args)
		}
		if got := envRedisAuth(cmd); got != "P@sSw0Rd" {
			return 0, fmt.Errorf("unexpected REDISCLI_AUTH: %q", got)
		}

		return 0, nil
	}
}

// writeFakeSSH compiles a Go program that stands in for the ssh binary of
// the tunnel: it writes its args to $FAKE_SSH_ARGS_FILE and listens on the
// local port of its -L spec, until killed.
func writeFakeSSH(t *testing.T, dir string) {
	t.Helper()

	const main = `package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	if err := os.WriteFile(os.Getenv("FAKE_SSH_ARGS_FILE"), []byte(strings.Join(args, "\n")), 0o600); err != nil {
		panic(err)
	}

	var spec string
	for i, a := range args {
		if a == "-L" && i+1 < len(args) {
			spec = args[i+1]
		}
	}
	parts := strings.Split(spec, ":")
	if len(parts) != 4 {
		fmt.Fprintln(os.Stderr, "fake ssh: unexpected -L spec:", spec)
		os.Exit(1)
	}

	l, err := net.Listen("tcp", "127.0.0.1:"+parts[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "fake ssh: listen:", err)
		os.Exit(1)
	}
	defer l.Close()

	// Accept and close the connections so the tunnel readiness dial succeeds.
	for {
		c, err := l.Accept()
		if err != nil {
			return
		}
		_ = c.Close()
	}
}
`

	src := filepath.Join(dir, "fake_ssh.go")
	if err := os.WriteFile(src, []byte(main), 0o600); err != nil {
		t.Fatal(err)
	}
	if out, err := exec.Command("go", "build", "-o", filepath.Join(dir, "ssh"), src).
		CombinedOutput(); err != nil {
		t.Fatalf("failed to build the fake ssh binary: %v\n%s", err, out)
	}
}

// readFakeSSHArgs returns the command line the fake ssh binary was invoked
// with.
func readFakeSSHArgs(argsFile string) ([]string, error) {
	argsOut, err := os.ReadFile(argsFile)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(argsOut), "\n"), nil
}

// checkTunnelStopped asserts the tunnel local port, read from the recorded
// ssh args, is not listening anymore once the command has exited.
func checkTunnelStopped(argsFile string) func(t *testing.T, _ *core.CheckFuncCtx) {
	return func(t *testing.T, _ *core.CheckFuncCtx) {
		t.Helper()

		sshArgs, err := readFakeSSHArgs(argsFile)
		if err != nil {
			t.Error(err)

			return
		}
		var local string
		for i, a := range sshArgs {
			if a == "-L" && i+1 < len(sshArgs) {
				local = strings.Split(sshArgs[i+1], ":")[1]
			}
		}

		conn, err := net.DialTimeout("tcp", "127.0.0.1:"+local, time.Second)
		if err == nil {
			_ = conn.Close()
			t.Error("the ssh tunnel is still listening on " + local + " after the command exited")
		}
	}
}
