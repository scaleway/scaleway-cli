package inference

import (
	"reflect"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/inference/v1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("inference").Groups = []string{"ai"}

	cmds.MustFind("inference", "deployment", "create").Override(deploymentCreateBuilder)
	cmds.MustFind("inference", "deployment", "delete").Override(deploymentDeleteBuilder)

	cmds.Add(listServerlessModels())

	humanv2.RegisterEnumFormat(humanv2.EnumFormats[inference.DeploymentStatus]{
		inference.DeploymentStatusCreating:  humanv2.EnumFormat{Color: color.FgBlue},
		inference.DeploymentStatusDeploying: humanv2.EnumFormat{Color: color.FgBlue},
		inference.DeploymentStatusDeleting:  humanv2.EnumFormat{Color: color.FgBlue},
		inference.DeploymentStatusError:     humanv2.EnumFormat{Color: color.FgRed},
		inference.DeploymentStatusReady:     humanv2.EnumFormat{Color: color.FgGreen},
		inference.DeploymentStatusLocked:    humanv2.EnumFormat{Color: color.FgRed},
	})

	humanv2.RegisterLayout(reflect.TypeFor[*inference.Deployment](), humanv2.Layout{
		Sections: []humanv2.Section{
			{Path: "Endpoints", Fields: []humanv2.Field{
				{Path: "ID"},
				{Path: "URL"},
				{Path: "DisableAuth"},
			}},
		},
	})

	return cmds
}
