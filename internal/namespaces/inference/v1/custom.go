package inference

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/inference/v1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("inference").Groups = []string{"ai"}

	human.RegisterMarshalerFunc(
		inference.DeploymentStatus(""),
		human.EnumMarshalFunc(deploymentStateMarshalSpecs),
	)

	human.RegisterMarshalerFunc(inference.Deployment{}, DeploymentMarshalerFunc)

	cmds.MustFind("inference", "deployment", "create").Override(deploymentCreateBuilder)
	cmds.MustFind("inference", "deployment", "delete").Override(deploymentDeleteBuilder)

	return cmds
}
