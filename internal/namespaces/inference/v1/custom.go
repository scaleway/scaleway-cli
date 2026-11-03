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

	cmds.Add(listServerlessModels())

	cmds.MustFind("inference", "deployment", "list").Override(core.ListViewBuilder("ID", "Name", "ProjectID", "Status", "Tags", "NodeTypeName", "Endpoints", "Size", "MinSize", "MaxSize", "ErrorMessage", "ModelID", "ModelName", "CreatedAt", "UpdatedAt", "Region"))
	cmds.MustFind("inference", "model", "list").Override(core.ListViewBuilder("ID", "Name", "ProjectID", "Tags", "Status", "Description", "ErrorMessage", "HasEula", "CreatedAt", "UpdatedAt", "Region", "NodesSupport", "ParameterSizeBits", "SizeBytes"))
	cmds.MustFind("inference", "node-type", "list").Override(core.ListViewBuilder("Name", "StockStatus", "Description", "Vcpus", "Memory", "Vram", "Disabled", "Beta", "CreatedAt", "UpdatedAt", "Gpus", "Region"))

	return cmds
}
