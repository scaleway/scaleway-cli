package instance_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/vpc/v2"
	"github.com/scaleway/scaleway-cli/v2/internal/testhelpers"
)

func Test_ListNICs(t *testing.T) {
	cmds := instance.GetCommands()
	cmds.Merge(vpc.GetCommands())

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			testhelpers.CreatePN(),
			createServer("Server"),
			createNIC(),
		),
		Cmd: "scw instance private-nic list server-id={{ .Server.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServer("Server"),
			testhelpers.DeletePN(),
		),
	}))
}

func Test_GetPrivateNIC(t *testing.T) {
	cmds := instance.GetCommands()
	cmds.Merge(vpc.GetCommands())

	t.Run("Get from ID", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			testhelpers.CreatePN(),
			createServer("Server"),
			createNIC(),
		),
		Cmd: "scw instance private-nic get server-id={{ .Server.ID }} private-nic-id={{ .NIC.PrivateNic.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServer("Server"),
			testhelpers.DeletePN(),
		),
	}))

	t.Run("Get from MAC address", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			testhelpers.CreatePN(),
			createServer("Server"),
			createNIC(),
		),
		Cmd: "scw instance private-nic get server-id={{ .Server.ID }} private-nic-id={{ .NIC.PrivateNic.MacAddress }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServer("Server"),
			testhelpers.DeletePN(),
		),
	}))
}
