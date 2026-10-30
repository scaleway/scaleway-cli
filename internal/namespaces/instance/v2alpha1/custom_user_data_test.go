package instance_test

import (
	"os"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	instanceCLIV1 "github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/testhelpers"
	instanceSDKV2 "github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
)

func Test_UserData(t *testing.T) {
	cmds := instanceV1AndV2Commands()

	t.Run("Get an existing key", core.Test(&core.TestConfig{
		BeforeFunc: core.BeforeFuncCombine(
			createServerV1("Server"),
			core.ExecBeforeCmd(
				"scw instance user-data set server-id={{.Server.ID}} key=happy content=true",
			),
		),
		Commands:  cmds,
		Cmd:       "scw instance user-data get server-id={{.Server.ID}} key=happy",
		AfterFunc: deleteServerV1("Server"),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))

	t.Run("Get a nonexistent key", core.Test(&core.TestConfig{
		BeforeFunc: createServerV1("Server"),
		Commands:   cmds,
		Cmd:        "scw instance user-data get server-id={{.Server.ID}} key=happy",
		AfterFunc:  deleteServerV1("Server"),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))

	t.Run("List", core.Test(&core.TestConfig{
		BeforeFunc: core.BeforeFuncCombine(
			createServerV1("Server"),
			core.ExecBeforeCmd(
				"scw instance user-data set {{ .Server.ID }} key=foo content=bar",
			),
			core.ExecBeforeCmd(
				"scw instance user-data set {{ .Server.ID }} key=bar content=foo",
			),
		),
		Commands:  cmds,
		Cmd:       "scw instance user-data list server-id={{ .Server.ID }}",
		AfterFunc: deleteServerV1("Server"),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}

func Test_UserData_SetUploadFile(t *testing.T) {
	content := "cloud-init or other file content"
	file, err := os.CreateTemp(t.TempDir(), "test")
	if err != nil {
		t.Fatalf("%s", err)
	}
	_, err = file.WriteString(content)
	if err != nil {
		t.Fatalf("%s", err)
	}

	t.Run("on random key", core.Test(&core.TestConfig{
		Commands: instanceV1AndV2Commands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd("Server", testServerCommand("stopped=true")),
			func(ctx *core.BeforeFuncCtx) error {
				ctx.Meta["filePath"] = file.Name()

				return nil
			},
		),
		Cmd: `scw instance user-data set key=foobar {{ .Server.ID }} content=@{{ .filePath }}`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				err = assertUserDataContentAtKey(t, ctx, "foobar", content)
				if err != nil {
					t.Error(err)
				}
			},
		),
		AfterFunc: core.AfterFuncCombine(
			func(_ *core.AfterFuncCtx) error {
				return deleteTmpFile(file)
			},
			deleteServerV1("Server"),
		),
	}))
}

func Test_UserData_SetCloudInit(t *testing.T) {
	content := `#cloud-config
groups:
  - ubuntu: [root,sys]
  - cloud-users
users:
  - default
  - name: foobar
  - name: nosshlogins
    ssh_redirect_user: true
`
	file, err := os.CreateTemp(t.TempDir(), "test")
	if err != nil {
		t.Fatalf("%s", err)
	}
	_, err = file.WriteString(content)
	if err != nil {
		t.Fatalf("%s", err)
	}

	t.Run("from file", core.Test(&core.TestConfig{
		Commands: instanceV1AndV2Commands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd("Server", testServerCommand("stopped=true")),
			func(ctx *core.BeforeFuncCtx) error {
				ctx.Meta["filePath"] = file.Name()

				return nil
			},
		),
		Cmd: `scw instance user-data set-cloud-init {{ .Server.ID }} content=@{{ .filePath }}`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				err = assertUserDataContentAtKey(t, ctx, "cloud-init", content)
				if err != nil {
					t.Error(err)
				}
			},
		),
		AfterFunc: core.AfterFuncCombine(
			func(_ *core.AfterFuncCtx) error {
				return deleteTmpFile(file)
			},
			deleteServerV1("Server"),
		),
	}))
}

func assertUserDataContentAtKey(t *testing.T,
	ctx *core.CheckFuncCtx,
	key, expected string,
) error {
	t.Helper()
	server := testhelpers.MapValue[*instanceCLIV1.ServerWithWarningsResponse](
		t,
		ctx.Meta,
		"Server",
	)

	ud, err := instanceSDKV2.NewAPI(ctx.Client).GetUserData(&instanceSDKV2.GetUserDataRequest{
		Zone:     "fr-par-1",
		ServerID: server.ID,
		Key:      key,
	}, scw.WithContext(t.Context()))
	if err != nil {
		return err
	}

	assert.Equal(t, expected, string(ud.Content))

	return nil
}

func deleteTmpFile(file *os.File) error {
	// We need to close this file explicitly because it is not closed by the os.Remove call on Windows
	// https://github.com/golang/go/issues/50510
	err := file.Close()
	if err != nil {
		return err
	}

	err = os.RemoveAll(file.Name())
	if err != nil {
		return err
	}

	return nil
}
