package ipfs

import (
	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func newIpfsRoot() *core.Command {
	return &core.Command{
		Short:     `IPFS Pinning service API`,
		Long:      `IPFS Pinning service API.`,
		Namespace: "ipfs",
		Groups:    []string{"labs"},
	}
}

func newIpnsRoot() *core.Command {
	return &core.Command{
		Short:     `IPFS Naming service API`,
		Long:      ``,
		Namespace: "ipns",
		Groups:    []string{"labs"},
	}
}

// GetCommands returns the list of commands for the 'ipfs' namespace.
func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	commands := cmds.GetAll()
	commands = append(commands[:0], commands[2:]...)
	commands = append([]*core.Command{newIpfsRoot(), newIpnsRoot()}, commands...)

	return core.NewCommands(commands...)
}
