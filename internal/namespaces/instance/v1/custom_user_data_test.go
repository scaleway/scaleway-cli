package instance_test

import (
	"os"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
)

func Test_UserDataGet(t *testing.T) {
	t.Run("Get an existing key", core.Test(&core.TestConfig{
		BeforeFunc: core.BeforeFuncCombine(
			createServer("Server"),
			core.ExecBeforeCmd(
				"scw instance user-data set server-id={{.Server.ID}} key=happy content=true",
			),
		),
		Commands:  instance.GetCommands(),
		Cmd:       "scw instance user-data get server-id={{.Server.ID}} key=happy",
		AfterFunc: deleteServer("Server"),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))

	t.Run("Get an nonexistent key", core.Test(&core.TestConfig{
		BeforeFunc: createServer("Server"),
		Commands:   instance.GetCommands(),
		Cmd:        "scw instance user-data get server-id={{.Server.ID}} key=happy",
		AfterFunc:  deleteServer("Server"),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))
}

func Test_UserDataList(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		BeforeFunc: core.BeforeFuncCombine(
			createServer("Server"),
			core.ExecBeforeCmd(
				"scw instance user-data set server-id={{ .Server.ID }} key=foo content=bar",
			),
			core.ExecBeforeCmd(
				"scw instance user-data set server-id={{ .Server.ID }} key=bar content=foo",
			),
		),
		Commands:  instance.GetCommands(),
		Cmd:       "scw instance user-data list server-id={{ .Server.ID }}",
		AfterFunc: deleteServer("Server"),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}

func Test_UserDataFileUploadOn(t *testing.T) {
	content := "cloud-init file content"
	file, err := os.CreateTemp(t.TempDir(), "test")
	if err != nil {
		t.Fatalf("%s", err)
	}
	_, err = file.WriteString(content)
	if err != nil {
		t.Fatalf("%s", err)
	}

	t.Run("cloud-init", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd("Server", testServerCommand("stopped=true")),
			func(ctx *core.BeforeFuncCtx) error {
				ctx.Meta["filePath"] = file.Name()

				return nil
			},
		),
		Cmd: `scw instance user-data set key=cloud-init server-id={{ .Server.ID }} content=@{{ .filePath }}`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			func(_ *core.AfterFuncCtx) error {
				// We need to close this file explicitly because it is not closed by the os.Remove call on windows
				// https://github.com/golang/go/issues/50510
				err = file.Close()
				if err != nil {
					return err
				}

				err = os.RemoveAll(file.Name())
				if err != nil {
					return err
				}

				return nil
			},
			deleteServer("Server"),
		),
	}))
}

func Test_UserDataFileUploadOnRandom(t *testing.T) {
	content := "cloud-init file content"
	file, err := os.CreateTemp(t.TempDir(), "test")
	if err != nil {
		t.Fatalf("%s", err)
	}
	_, err = file.WriteString(content)
	if err != nil {
		t.Fatalf("%s", err)
	}

	t.Run("key", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd("Server", testServerCommand("stopped=true")),
			func(ctx *core.BeforeFuncCtx) error {
				ctx.Meta["filePath"] = file.Name()

				return nil
			},
		),
		Cmd: `scw instance user-data set key=foobar server-id={{ .Server.ID }} content=@{{ .filePath }}`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			func(_ *core.AfterFuncCtx) error {
				// We need to close this file explicitly because it is not closed by the os.Remove call on windows
				// https://github.com/golang/go/issues/50510
				err = file.Close()
				if err != nil {
					return err
				}

				err = os.RemoveAll(file.Name())
				if err != nil {
					return err
				}

				return nil
			},
			deleteServer("Server"),
		),
	}))
}
