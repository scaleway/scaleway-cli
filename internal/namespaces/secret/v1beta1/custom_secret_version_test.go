package secret_test

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/v2/core"
	secret "github.com/scaleway/scaleway-cli/v2/internal/namespaces/secret/v1beta1"
)

func Test_AccessSecret(t *testing.T) {
	cmds := secret.GetCommands()

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createSecret("cli-test-access-secret"),
			createSecretVersion("{\"key\":\"value\"}"),
		),
		Cmd: "scw secret version access {{ .Secret.ID }} revision={{ .SecretVersion.Revision }} field=key raw=true",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				rawResult, isBytes := ctx.Result.(core.RawResult)
				if !isBytes {
					t.Fatalf("Expecting result to be bytes")
				}
				assert.Equal(t, core.RawResult("value"), rawResult)
			},
		),
		AfterFunc: core.AfterFuncCombine(
			deleteSecret(),
		),
	}))
}
