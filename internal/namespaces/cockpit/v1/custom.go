package cockpit

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("cockpit").Groups = []string{"monitoring"}

	cmds.MustFind("cockpit", "token", "get").Override(cockpitTokenGetBuilder)
	cmds.Add(cockpitConfigRoot())
	cmds.Add(cockpitConfigGetCommand())
	addCockpitConfigGetSeeAlso(cmds)

	return cmds
}

func addCockpitConfigGetSeeAlso(cmds *core.Commands) {
	for _, cmd := range cmds.GetAll() {
		if cmd.Namespace != "cockpit" || cmd.Run == nil || cmd.Resource == "config" {
			continue
		}

		cmd.SeeAlsos = append(cmd.SeeAlsos, &core.SeeAlso{
			Command: "scw cockpit config get",
			Short:   "Generate a data source configuration snippet",
		})
	}
}
