package container

import (
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	container "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	human.RegisterMarshalerFunc(container.NamespaceStatus(""), human.EnumMarshalFunc(namespaceStatusMarshalSpecs))
	human.RegisterMarshalerFunc(container.ContainerStatus(""), human.EnumMarshalFunc(containerStatusMarshalSpecs))
	human.RegisterMarshalerFunc(container.CronStatus(""), human.EnumMarshalFunc(cronStatusMarshalSpecs))

	return cmds
}
