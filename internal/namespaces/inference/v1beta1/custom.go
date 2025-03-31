package inference

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	inference "github.com/scaleway/scaleway-sdk-go/api/inference/v1beta1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("inference").Groups = []string{"ai"}

	human.RegisterMarshalerFunc(
		inference.DeploymentStatus(""),
		human.EnumMarshalFunc(deployementStateMarshalSpecs),
	)

	human.RegisterMarshalerFunc(inference.Deployment{}, DeploymentMarshalerFunc)
	human.RegisterMarshalerFunc([]*inference.Model{}, ListModelMarshalerFunc)

	cmds.MustFind("inference", "deployment", "create").Override(deploymentCreateBuilder)
	cmds.MustFind("inference", "deployment", "delete").Override(deploymentDeleteBuilder)
	cmds.MustFind("inference", "endpoint", "create").Override(endpointCreateBuilder)

	return cmds
}
