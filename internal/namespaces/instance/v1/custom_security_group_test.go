package instance_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
)

func Test_SecurityGroupCreate(t *testing.T) {
	// Keep the existing JSON response shape covered as an output contract.
	t.Run("Create", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd:      "scw instance security-group create name=cli-sg-charming-rubin -o json",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.ExecAfterCmd(
			"scw instance security-group delete {{ .CmdResult.SecurityGroup.ID }}",
		),
	}))

	// The API error must be surfaced when the create call fails.
	t.Run("Error", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd:      "scw instance security-group create name=cli-sg-charming-rubin -o json",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))
}

func Test_SecurityGroupGet(t *testing.T) {
	t.Run("Get", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createSecurityGroup("SecurityGroup"),
		),
		Cmd: "scw instance security-group get {{ .SecurityGroup.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteSecurityGroup("SecurityGroup"),
	}))
}
