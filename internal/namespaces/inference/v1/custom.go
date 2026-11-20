package inference

import (
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/inference/v1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("inference").Groups = []string{"ai"}

	humanv2.RegisterEnumFormat(deploymentStateFormat)
	humanv2.RegisterLayout(reflect.TypeFor[inference.Deployment](), deploymentLayout)

	cmds.MustFind("inference", "deployment", "create").Override(deploymentCreateBuilder)
	cmds.MustFind("inference", "deployment", "delete").Override(deploymentDeleteBuilder)

	cmds.Add(listServerlessModels())

	return cmds
}
