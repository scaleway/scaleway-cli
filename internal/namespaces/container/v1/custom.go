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
	human.RegisterMarshalerFunc(
		container.DomainStatus(""),
		human.EnumMarshalFunc(domainStatusMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		container.TriggerStatus(""),
		human.EnumMarshalFunc(triggerStatusMarshalSpecs),
	)

	cmds.MustFind("container", "container", "redeploy").Override(containerContainerRedeployBuilder)
	cmds.MustFind("container", "container", "create").Override(containerContainerCreateBuilder)
	cmds.MustFind("container", "container", "update").Override(containerContainerUpdateBuilder)
	cmds.MustFind("container", "namespace", "create").Override(containerNamespaceCreateBuilder)
	cmds.MustFind("container", "namespace", "update").Override(containerNamespaceUpdateBuilder)
	cmds.MustFind("container", "namespace", "delete").Override(containerNamespaceDeleteBuilder)

	// Logs and Metrics
	cmds.Merge(core.NewCommands(
		containerLogs(),
	))

	// Deploy
	if cmdDeploy := containerDeployCommand(); cmdDeploy != nil {
		cmds.Add(cmdDeploy)
	}

	cmds.MustFind("container", "container", "list").
		Override(core.ListViewBuilder("ID", "Name", "NamespaceID", "Description", "Status", "ErrorMessage", "CreatedAt", "UpdatedAt", "MinScale", "MaxScale", "MemoryLimitBytes", "MvcpuLimit", "LocalStorageLimitBytes", "Timeout", "Privacy", "Image", "Protocol", "Port", "HTTPSConnectionsOnly", "Sandbox", "Tags", "PrivateNetworkID", "Command", "Args", "PublicEndpoint", "PrivateEndpoint", "DefaultPublicEndpointEnabled", "Region"))
	cmds.MustFind("container", "domain", "list").
		Override(core.ListViewBuilder("ID", "ContainerID", "Hostname", "Status", "ErrorMessage", "CreatedAt", "UpdatedAt", "Tags", "Region"))
	cmds.MustFind("container", "namespace", "list").
		Override(core.ListViewBuilder("ID", "Name", "OrganizationID", "ProjectID", "Description", "Status", "ErrorMessage", "Tags", "CreatedAt", "UpdatedAt", "Region"))
	cmds.MustFind("container", "trigger", "list").
		Override(core.ListViewBuilder("ID", "Name", "Description", "Tags", "Status", "ErrorMessage", "ContainerID", "SourceType", "CreatedAt", "UpdatedAt", "Region"))

	return cmds
}
