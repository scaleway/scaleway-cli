package ipfs

import (
	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

// GetCommands returns the list of commands for the 'ipfs' namespace.
func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("ipns").Groups = []string{"labs"}
	cmds.MustFind("ipfs").Groups = []string{"labs"}
	return cmds
}
