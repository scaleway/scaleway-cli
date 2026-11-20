package mongodb_test

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/mongodb/v1"
	mongodbSDK "github.com/scaleway/scaleway-sdk-go/api/mongodb/v1"
	"github.com/stretchr/testify/assert"
)

const (
	mongodbTestVersion = "7.0"
	mongodbTestNode    = "MGDB-PLAY2-NANO"
	mongodbTestUser    = "admin"
	mongodbTestPass    = "P@ssw0rd123"

	// connectInterruptedInstanceID is a syntactically valid uuid. The interrupted
	// subtest references no template field resolved from the API, so it never
	// calls the API and does not need a real instance.
	connectInterruptedInstanceID = "66666666-7777-8888-9999-000000000000"

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

func checkNoCertLeftover(t *testing.T, instanceID string) {
	t.Helper()

	_, err := os.Stat(filepath.Join(os.TempDir(), "mongodb-cert-"+instanceID+".crt"))
	assert.True(t, os.IsNotExist(err), "temporary certificate must be removed")
}

func Test_InstanceConnect(t *testing.T) {
	t.Run("not-installed", core.Test(&core.TestConfig{
		Commands:        mongodb.GetCommands(),
		DisableParallel: true, // MgdbNodesTypePlay2Nano quota is 3
		BeforeFunc:      createConnectInstance("connect-test-not-installed"),
		Cmd: "scw mongodb instance connect {{ .Instance.ID }}" +
			` template="not-a-real-tool-xyz --uri {{"{{"}}.ConnURL{{"}}"}}"`,
		PromptResponseMocks: []string{mongodbTestPass},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckStderrContains(`Exec: "not-a-real-tool-xyz": executable file not found`),
		),
		AfterFunc: deleteConnectInstance(),
	}))

	t.Run("simple", core.Test(&core.TestConfig{
		Commands:            mongodb.GetCommands(),
		DisableParallel:     true,
		BeforeFunc:          createConnectInstance("connect-test-simple"),
		Cmd:                 "scw mongodb instance connect {{ .Instance.ID }}",
		PromptResponseMocks: []string{mongodbTestPass},
		OverrideExec: func(ctx *core.ExecFuncCtx, cmd *exec.Cmd) (int, error) {
			instance := ctx.Meta["Instance"].(*mongodbSDK.Instance)
			if len(cmd.Args) != 4 {
				return 0, fmt.Errorf("unexpected args: %v", cmd.Args)
			}
			if cmd.Args[0] != "mongosh" || cmd.Args[2] != "--tlsCAFile" {
				return 0, fmt.Errorf("unexpected args: %v", cmd.Args)
			}
			if _, err := os.Stat(cmd.Args[3]); err != nil {
				return 0, fmt.Errorf("certificate file not found: %w", err)
			}

			parsed, err := url.Parse(cmd.Args[1])
			if err != nil {
				return 0, fmt.Errorf("invalid ConnURL: %w", err)
			}
			if parsed.Scheme != "mongodb+srv" {
				return 0, fmt.Errorf("unexpected scheme: %s", parsed.Scheme)
			}
			if parsed.User.Username() != mongodbTestUser {
				return 0, fmt.Errorf("unexpected username: %s", parsed.User.Username())
			}
			pass, ok := parsed.User.Password()
			if !ok || pass != mongodbTestPass {
				return 0, errors.New("unexpected password in ConnURL")
			}
			if len(instance.Endpoints) == 0 || parsed.Host != instance.Endpoints[0].DNSRecord {
				return 0, fmt.Errorf("unexpected host: %s", parsed.Host)
			}

			return 0, nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				checkNoCertLeftover(t, ctx.Meta["Instance"].(*mongodbSDK.Instance).ID)
			},
		),
		AfterFunc: deleteConnectInstance(),
	}))

	t.Run("password-arg", core.Test(&core.TestConfig{
		Commands:        mongodb.GetCommands(),
		DisableParallel: true,
		BeforeFunc:      createConnectInstance("connect-test-password-arg"),
		Cmd: "scw mongodb instance connect {{ .Instance.ID }} password=" +
			mongodbTestPass,
		PromptResponseMocks: []string{"should-not-be-prompted"},
		OverrideExec: func(_ *core.ExecFuncCtx, cmd *exec.Cmd) (int, error) {
			if len(cmd.Args) < 2 {
				return 0, fmt.Errorf("unexpected args: %v", cmd.Args)
			}
			parsed, err := url.Parse(cmd.Args[1])
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
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				checkNoCertLeftover(t, ctx.Meta["Instance"].(*mongodbSDK.Instance).ID)
			},
		),
		AfterFunc: deleteConnectInstance(),
	}))

	t.Run("username-override", core.Test(&core.TestConfig{
		Commands:        mongodb.GetCommands(),
		DisableParallel: true,
		BeforeFunc:      createConnectInstance("connect-test-username-override"),
		Cmd: "scw mongodb instance connect {{ .Instance.ID }}" +
			" username=not-the-provisioned-admin password=" + mongodbTestPass,
		OverrideExec: func(_ *core.ExecFuncCtx, cmd *exec.Cmd) (int, error) {
			if len(cmd.Args) < 2 {
				return 0, fmt.Errorf("unexpected args: %v", cmd.Args)
			}
			parsed, err := url.Parse(cmd.Args[1])
			if err != nil {
				return 0, fmt.Errorf("invalid ConnURL: %w", err)
			}
			if parsed.User.Username() != "not-the-provisioned-admin" {
				return 0, fmt.Errorf("username override ignored: %s", parsed.User.Username())
			}

			return 0, nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				checkNoCertLeftover(t, ctx.Meta["Instance"].(*mongodbSDK.Instance).ID)
			},
		),
		AfterFunc: deleteConnectInstance(),
	}))

	t.Run("template", core.Test(&core.TestConfig{
		Commands: mongodb.GetCommands(),
		// No API: username/password set and template only uses {{.Username}}.
		Cmd: "scw mongodb instance connect " + connectInterruptedInstanceID +
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

	t.Run("interrupted", core.Test(&core.TestConfig{
		Commands: mongodb.GetCommands(),
		Cmd: "scw mongodb instance connect " + connectInterruptedInstanceID +
			` username=admin password=secret template="sh -c 'sleep 1; kill -INT $PPID; exec sleep 300'"`,
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(130),
			func(t *testing.T, _ *core.CheckFuncCtx) {
				t.Helper()
				checkNoCertLeftover(t, connectInterruptedInstanceID)
			},
		),
	}))
}
