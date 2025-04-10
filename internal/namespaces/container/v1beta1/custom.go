package container

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	container "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("container").Groups = []string{"serverless"}

	human.RegisterMarshalerFunc(
		container.NamespaceStatus(""),
		human.EnumMarshalFunc(namespaceStatusMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		container.ContainerStatus(""),
		human.EnumMarshalFunc(containerStatusMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		container.CronStatus(""),
		human.EnumMarshalFunc(cronStatusMarshalSpecs),
	)

	cmds.MustFind("container", "container", "deploy").Override(containerContainerDeployBuilder)
	cmds.MustFind("container", "container", "create").Override(containerContainerCreateBuilder)
	cmds.MustFind("container", "container", "update").Override(containerContainerUpdateBuilder)
	cmds.MustFind("container", "namespace", "create").Override(containerNamespaceCreateBuilder)
	cmds.MustFind("container", "namespace", "update").Override(containerNamespaceUpdateBuilder)
	cmds.MustFind("container", "namespace", "delete").Override(containerNamespaceDeleteBuilder)

	if cmdDeploy := containerDeployCommand(); cmdDeploy != nil {
		cmds.Add(cmdDeploy)
	}

	return cmds
}
