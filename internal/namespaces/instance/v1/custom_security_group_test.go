package instance

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func Test_SecurityGroupGet(t *testing.T) {
	t.Run("Get", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createSecurityGroup("SecurityGroup"),
		),
		Cmd:       "scw instance security-group get {{ .SecurityGroup.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteSecurityGroup("SecurityGroup"),
	}))
}
