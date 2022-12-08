package container

import (
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	container "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.Merge(core.NewCommands(
		containerContext(),
		containerContextCreate(),
		containerContextDelete(),
		containerContextStart(),
		containerContextStop(),
	))

	human.RegisterMarshalerFunc(container.NamespaceStatus(""), human.EnumMarshalFunc(namespaceStatusMarshalSpecs))
	human.RegisterMarshalerFunc(container.ContainerStatus(""), human.EnumMarshalFunc(containerStatusMarshalSpecs))
	human.RegisterMarshalerFunc(container.CronStatus(""), human.EnumMarshalFunc(cronStatusMarshalSpecs))

	return cmds
}
