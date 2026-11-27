package mongodb_test

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"slices"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/mongodb/v1"
	mongodbSDK "github.com/scaleway/scaleway-sdk-go/api/mongodb/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	mongodbTestVersion = "7.0"
	mongodbTestNode    = "MGDB-PLAY2-NANO"
	mongodbTestUser    = "admin"
	mongodbTestPass    = "P@ssw0rd123"

	// connectNoAPIInstanceID is a syntactically valid uuid used by subtests that
	// never resolve API-backed template fields.
	connectNoAPIInstanceID = "66666666-7777-8888-9999-000000000000"

	baseConnectCreate = "scw mongodb instance create --wait name=%s version=" + mongodbTestVersion +
		" node-type=" + mongodbTestNode + " node-amount=1 user-name=" + mongodbTestUser +
		" password=" + mongodbTestPass + " endpoints.0.public-network={}"
)

func createConnectInstance(name string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd("Instance", fmt.Sprintf(baseConnectCreate, name))
}

func deleteConnectInstance() core.AfterFunc {
	return core.ExecAfterCmd("scw mongodb instance delete {{ .Instance.ID }}")
}

func assertCertRemoved(t *testing.T, certPath string) {
	t.Helper()
	require.NotEmpty(t, certPath, "certificate path was not captured")
	_, err := os.Stat(certPath)
	assert.True(t, os.IsNotExist(err), "temporary certificate must be removed: %s", certPath)
}

func parseDefaultMongoshArgs(cmd *exec.Cmd) (connURL, username, certPath string, err error) {
	// mongosh <url> --username <user> --tlsCAFile <path>
	if len(cmd.Args) != 6 {
		return "", "", "", fmt.Errorf("unexpected args: %v", cmd.Args)
	}
	if cmd.Args[0] != "mongosh" || cmd.Args[2] != "--username" || cmd.Args[4] != "--tlsCAFile" {
		return "", "", "", fmt.Errorf("unexpected args: %v", cmd.Args)
	}

	return cmd.Args[1], cmd.Args[3], cmd.Args[5], nil
}

func Test_InstanceConnect(t *testing.T) {
	t.Run("not-installed", core.Test(&core.TestConfig{
		Commands:        mongodb.GetCommands(),
		DisableParallel: true, // MgdbNodesTypePlay2Nano quota is 3
		BeforeFunc:      createConnectInstance("connect-test-not-installed"),
		Cmd: "scw mongodb instance connect {{ .Instance.ID }}" +
			` template="not-a-real-tool-xyz --uri {{"{{"}}.ConnURL{{"}}"}}"`,
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckStderrContains(`Exec: "not-a-real-tool-xyz": executable file not found`),
		),
		AfterFunc: deleteConnectInstance(),
	}))

	t.Run("simple", func(t *testing.T) {
		var certPath string

		core.Test(&core.TestConfig{
			Commands:        mongodb.GetCommands(),
			DisableParallel: true,
			BeforeFunc:      createConnectInstance("connect-test-simple"),
			Cmd:             "scw mongodb instance connect {{ .Instance.ID }}",
			OverrideExec: func(ctx *core.ExecFuncCtx, cmd *exec.Cmd) (int, error) {
				instance := ctx.Meta["Instance"].(*mongodbSDK.Instance)
				connURL, username, path, err := parseDefaultMongoshArgs(cmd)
				if err != nil {
					return 0, err
				}
				certPath = path
				if _, err := os.Stat(certPath); err != nil {
					return 0, fmt.Errorf("certificate file not found: %w", err)
				}
				if username != mongodbTestUser {
					return 0, fmt.Errorf("unexpected --username: %s", username)
				}

				parsed, err := url.Parse(connURL)
				if err != nil {
					return 0, fmt.Errorf("invalid ConnURL: %w", err)
				}
				if parsed.Scheme != "mongodb+srv" {
					return 0, fmt.Errorf("unexpected scheme: %s", parsed.Scheme)
				}
				if parsed.User != nil {
					return 0, fmt.Errorf(
						"interactive path must not put credentials in ConnURL: %s",
						connURL,
					)
				}
				if len(instance.Endpoints) == 0 || parsed.Host != instance.Endpoints[0].DNSRecord {
					return 0, fmt.Errorf("unexpected host: %s", parsed.Host)
				}

				return 0, nil
			},
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				func(t *testing.T, _ *core.CheckFuncCtx) {
					t.Helper()
					assertCertRemoved(t, certPath)
				},
			),
			AfterFunc: deleteConnectInstance(),
		})(t)
	})

	t.Run("password-arg", func(t *testing.T) {
		var certPath string

		core.Test(&core.TestConfig{
			Commands:        mongodb.GetCommands(),
			DisableParallel: true,
			BeforeFunc:      createConnectInstance("connect-test-password-arg"),
			Cmd: "scw mongodb instance connect {{ .Instance.ID }} password=" +
				mongodbTestPass,
			OverrideExec: func(_ *core.ExecFuncCtx, cmd *exec.Cmd) (int, error) {
				connURL, _, path, err := parseDefaultMongoshArgs(cmd)
				if err != nil {
					return 0, err
				}
				certPath = path

				parsed, err := url.Parse(connURL)
				if err != nil {
					return 0, fmt.Errorf("invalid ConnURL: %w", err)
				}
				pass, ok := parsed.User.Password()
				if !ok || pass != mongodbTestPass {
					return 0, fmt.Errorf("unexpected password in ConnURL: %q", pass)
				}

				return 0, nil
			},
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				func(t *testing.T, _ *core.CheckFuncCtx) {
					t.Helper()
					assertCertRemoved(t, certPath)
				},
			),
			AfterFunc: deleteConnectInstance(),
		})(t)
	})

	t.Run("username-override", func(t *testing.T) {
		var certPath string

		core.Test(&core.TestConfig{
			Commands:        mongodb.GetCommands(),
			DisableParallel: true,
			BeforeFunc:      createConnectInstance("connect-test-username-override"),
			Cmd: "scw mongodb instance connect {{ .Instance.ID }}" +
				" username=not-the-provisioned-admin password=" + mongodbTestPass,
			OverrideExec: func(_ *core.ExecFuncCtx, cmd *exec.Cmd) (int, error) {
				connURL, username, path, err := parseDefaultMongoshArgs(cmd)
				if err != nil {
					return 0, err
				}
				certPath = path
				if username != "not-the-provisioned-admin" {
					return 0, fmt.Errorf("username override ignored in --username: %s", username)
				}

				parsed, err := url.Parse(connURL)
				if err != nil {
					return 0, fmt.Errorf("invalid ConnURL: %w", err)
				}
				if parsed.User.Username() != "not-the-provisioned-admin" {
					return 0, fmt.Errorf(
						"username override ignored in ConnURL: %s",
						parsed.User.Username(),
					)
				}

				return 0, nil
			},
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				func(t *testing.T, _ *core.CheckFuncCtx) {
					t.Helper()
					assertCertRemoved(t, certPath)
				},
			),
			AfterFunc: deleteConnectInstance(),
		})(t)
	})

	t.Run("template", core.Test(&core.TestConfig{
		Commands: mongodb.GetCommands(),
		// No API: username/password set and template only uses {{.Username}}.
		Cmd: "scw mongodb instance connect " + connectNoAPIInstanceID +
			` username=` + mongodbTestUser + ` password=` + mongodbTestPass +
			` template="echo {{"{{"}}.Username{{"}}"}}"`,
		OverrideExec: func(_ *core.ExecFuncCtx, cmd *exec.Cmd) (int, error) {
			want := []string{"echo", mongodbTestUser}
			if !slices.Equal(cmd.Args, want) {
				return 0, fmt.Errorf("unexpected rendered command: %v", cmd.Args)
			}

			return 0, nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
		),
	}))
}
