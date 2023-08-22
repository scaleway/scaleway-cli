package document_db

import "github.com/scaleway/scaleway-cli/v2/internal/core"

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("document-db", "engine", "list").Override(engineListBuilder)

	return cmds
}
