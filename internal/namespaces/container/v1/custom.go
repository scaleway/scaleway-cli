package container

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/container/v1"
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

	cmds.MustFind("container", "container", "redeploy").Override(containerContainerRedeployBuilder)
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
