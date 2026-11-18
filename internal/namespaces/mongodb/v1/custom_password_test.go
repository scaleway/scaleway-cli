package mongodb_test

import (
	"fmt"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	mongodb "github.com/scaleway/scaleway-cli/v2/internal/namespaces/mongodb/v1"
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
		Check:     core.TestCheckExitCode(0),
		AfterFunc: core.ExecAfterCmd("scw mongodb instance delete {{ .CmdResult.ID }}"),
	}))
}
