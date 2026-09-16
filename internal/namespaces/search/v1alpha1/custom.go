package search

import "github.com/scaleway/scaleway-cli/v2/core"

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("search").Groups = []string{"utility"}

	cmds.MustFind("search", "resource", "search").Override(searchBuilder)

	return cmds
}
