package function_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	function "github.com/scaleway/scaleway-cli/v2/internal/namespaces/function/v1beta1"
	functionSDK "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func Test_Deploy(t *testing.T) {
	functionName := "cli-test-function-deploy"
	testZip := "testfixture/gofunction.zip"

	_, err := os.Stat(testZip)
	if err != nil {
		t.Fatal("test zip not found", err)
	}

	commands := function.GetCommands()

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: commands,
		Cmd: fmt.Sprintf(
			"scw function deploy name=%s runtime=go124 zip-file=%s",
			functionName,
			testZip,
		),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			testDeleteFunctionNamespaceAfter(functionName),
		),
	}))
}

func testDeleteFunctionNamespaceAfter(functionName string) func(*core.AfterFuncCtx) error {
	return func(ctx *core.AfterFuncCtx) error {
		api := functionSDK.NewAPI(ctx.Client)

		namespaces, err := api.ListNamespaces(&functionSDK.ListNamespacesRequest{
			Name: &functionName,
		}, scw.WithAllPages())
		if err != nil {
			return err
		}

		var namespaceID string
		for _, namespace := range namespaces.Namespaces {
			if namespace.Name == functionName {
				namespaceID = namespace.ID

				break
			}
		}

		if namespaceID == "" {
			return errors.New("namespace not found")
		}

		return core.ExecAfterCmd("scw function namespace delete " + namespaceID)(ctx)
	}
}
