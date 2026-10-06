package core_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/stretchr/testify/require"
)

type testFileContentPool struct {
	UserData map[string][]byte
}

type testFileContentRequest struct {
	CloudInit []byte
	Content   string
	Reader    io.Reader
	UserData  map[string][]byte
	Pools     []*testFileContentPool
}

const testFileContent = "hello from file"

func testFileContentCommands() *core.Commands {
	return core.NewCommands(
		&core.Command{
			Namespace: "test",
			Resource:  "file-content",
			ArgSpecs: core.ArgSpecs{
				{Name: "cloud-init", CanLoadFile: true},
				{Name: "content", CanLoadFile: true},
				{Name: "reader", CanLoadFile: true},
				{Name: "user-data.{key}", CanLoadFile: true},
				{Name: "pools.{index}.user-data.{key}", CanLoadFile: true},
			},
			ArgsType:             reflect.TypeFor[testFileContentRequest](),
			AllowAnonymousClient: true,
			Run: func(_ context.Context, argsI any) (i any, e error) {
				req := argsI.(*testFileContentRequest)

				readerContent, err := io.ReadAll(req.Reader)
				if err != nil {
					return nil, err
				}

				return fmt.Sprintf(
					"cloud-init=%q\ncontent=%q\nreader=%q\nuser-data-cloud-init=%q\nuser-data-other=%q\npool-0-cloud-init=%q\npool-1-other=%q",
					string(req.CloudInit),
					req.Content,
					string(readerContent),
					string(req.UserData["cloud-init"]),
					string(req.UserData["other"]),
					string(req.Pools[0].UserData["cloud-init"]),
					string(req.Pools[1].UserData["other"]),
				), nil
			},
		},
		&core.Command{
			Namespace: "test",
			Resource:  "file-content-missing-field",
			ArgSpecs: core.ArgSpecs{
				{Name: "nonexistent.{key}", CanLoadFile: true},
				{Name: "user-data.{key}", CanLoadFile: true},
			},
			ArgsType:             reflect.TypeFor[testFileContentRequest](),
			AllowAnonymousClient: true,
			Run: func(_ context.Context, argsI any) (i any, e error) {
				req := argsI.(*testFileContentRequest)

				return fmt.Sprintf(
					"user-data-cloud-init=%q",
					string(req.UserData["cloud-init"]),
				), nil
			},
		},
	)
}

func Test_FileContentArgs(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "cloud-init.yml")
	require.NoError(t, os.WriteFile(filePath, []byte(testFileContent), 0o644))

	t.Run("loads files", core.Test(&core.TestConfig{
		Commands:   testFileContentCommands(),
		BeforeFunc: core.BeforeFuncStoreInMeta("filePath", filePath),
		Cmd: "scw test file-content " +
			"cloud-init=@{{ .filePath }} " +
			"content=@{{ .filePath }} " +
			"reader=@{{ .filePath }} " +
			"user-data.cloud-init=@{{ .filePath }} " +
			"user-data.other=plain " +
			"pools.0.user-data.cloud-init=@{{ .filePath }} " +
			"pools.1.user-data.other=plain",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckStdout(
				fmt.Sprintf(
					"cloud-init=%q\ncontent=%q\nreader=%q\nuser-data-cloud-init=%q\nuser-data-other=%q\npool-0-cloud-init=%q\npool-1-other=%q\n",
					testFileContent,
					testFileContent,
					testFileContent,
					testFileContent,
					"plain",
					testFileContent,
					"plain",
				),
			),
		),
	}))

	t.Run("missing file", core.Test(&core.TestConfig{
		Commands: testFileContentCommands(),
		Cmd:      "scw test file-content user-data.cloud-init=@does-not-exist.yml",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckStderrContains("Could not open requested file"),
		),
	}))

	t.Run("missing field", core.Test(&core.TestConfig{
		Commands: testFileContentCommands(),
		Cmd:      "scw test file-content-missing-field user-data.cloud-init=plain",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckStdout("user-data-cloud-init=\"plain\"\n"),
		),
	}))
}
