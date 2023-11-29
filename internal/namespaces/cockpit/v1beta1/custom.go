package cockpit

import (
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	cockpit "github.com/scaleway/scaleway-sdk-go/api/cockpit/v1beta1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.Merge(core.NewCommands(
		cockpitWaitCommand(),
	))

	human.RegisterMarshalerFunc(cockpit.CockpitStatus(""), human.EnumMarshalFunc(cockpitStatusMarshalSpecs))

	cmds.MustFind("cockpit", "cockpit", "activate").Override(cockpitCockpitActivateBuilder)
	cmds.MustFind("cockpit", "cockpit", "deactivate").Override(cockpitCockpitDeactivateBuilder)
	cmds.MustFind("cockpit", "cockpit", "get").Override(cockpitCockpitGetBuilder)
	cmds.MustFind("cockpit", "token", "get").Override(cockpitTokenGetBuilder)

	return cmds
}
