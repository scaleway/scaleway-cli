package mongodb_test

import (
	"fmt"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	mongodb "github.com/scaleway/scaleway-cli/v2/internal/namespaces/mongodb/v1"
	"github.com/stretchr/testify/require"
)

const (
	mongodbNodeType = "MGDB-PLAY2-NANO"
	mongodbVersion  = "7.0"
	mongodbUser     = "cliuser"
)

func Test_PasswordGenerator(t *testing.T) {
	t.Run("Create instance", core.Test(&core.TestConfig{
		Commands: mongodb.GetCommands(),
		Cmd: fmt.Sprintf(
			"scw mongodb instance create --wait name=cli-pwd-gen-create node-type=%s version=%s user-name=%s generate-password=true endpoints.0.public-network={}",
			mongodbNodeType,
			mongodbVersion,
			mongodbUser,
		),
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				result, ok := ctx.Result.(mongodb.CreateInstanceResult)
				require.True(t, ok)
				require.NotEmpty(t, result.Password)
			},
		),
		AfterFunc: core.ExecAfterCmd("scw mongodb instance delete {{ .CmdResult.ID }}"),
	}))
}
