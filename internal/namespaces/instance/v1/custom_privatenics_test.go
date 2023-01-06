package instance

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/vpc/v1"
)

func Test_ListNICs(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		// Temporary in waiting for the private network support in the CLI
		Cmd: "scw instance private-nic list server-id=4fe24c2a-3c65-4530-b274-574b22ba3d14",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
	}))
}

func Test_GetPrivateNIC(t *testing.T) {
	cmds := GetCommands()
	cmds.Merge(vpc.GetCommands())

	t.Run("Get from ID", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createPN(),
			createServer("Server"),
			createNIC(),
		),
		Cmd: "scw instance private-nic get server-id={{ .Server.ID }} private-nic-id={{ .NIC.PrivateNic.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
	}))

	t.Run("Get from MAC address", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createPN(),
			createServer("Server"),
			createNIC(),
		),
		Cmd: "scw instance private-nic get server-id={{ .Server.ID }} private-nic-id={{ .NIC.PrivateNic.MacAddress }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
	}))
}
