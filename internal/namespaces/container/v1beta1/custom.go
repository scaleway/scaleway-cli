package container

import (
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	container "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	human.RegisterMarshalerFunc(container.NamespaceStatus(""), human.EnumMarshalFunc(namespaceStatusMarshalSpecs))
	human.RegisterMarshalerFunc(container.ContainerStatus(""), human.EnumMarshalFunc(containerStatusMarshalSpecs))
	human.RegisterMarshalerFunc(container.CronStatus(""), human.EnumMarshalFunc(cronStatusMarshalSpecs))

	cmds.MustFind("container", "container", "deploy").Override(containerContainerDeployBuilder)

	cmds.Add(containerDeployCommand())

	return cmds
}
