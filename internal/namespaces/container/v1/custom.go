package container

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/container/v1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("container").Groups = []string{"serverless"}

	humanv2.RegisterEnumFormat(humanv2.EnumFormats[container.NamespaceStatus]{
		container.NamespaceStatusCreating:      humanv2.EnumFormat{Color: color.FgBlue},
		container.NamespaceStatusDeleting:      humanv2.EnumFormat{Color: color.FgBlue},
		container.NamespaceStatusError:         humanv2.EnumFormat{Color: color.FgRed},
		container.NamespaceStatusLocked:        humanv2.EnumFormat{Color: color.FgRed},
		container.NamespaceStatusLocking:       humanv2.EnumFormat{Color: color.FgBlue},
		container.NamespaceStatusReady:         humanv2.EnumFormat{Color: color.FgGreen},
		container.NamespaceStatusUpdating:      humanv2.EnumFormat{Color: color.FgBlue},
		container.NamespaceStatusUpgrading:     humanv2.EnumFormat{Color: color.FgBlue},
		container.NamespaceStatusUnknownStatus: humanv2.EnumFormat{Color: color.Faint},
	})
	humanv2.RegisterEnumFormat(humanv2.EnumFormats[container.ContainerStatus]{
		container.ContainerStatusCreating:      humanv2.EnumFormat{Color: color.FgBlue},
		container.ContainerStatusDeleting:      humanv2.EnumFormat{Color: color.FgBlue},
		container.ContainerStatusError:         humanv2.EnumFormat{Color: color.FgRed},
		container.ContainerStatusLocked:        humanv2.EnumFormat{Color: color.FgRed},
		container.ContainerStatusLocking:       humanv2.EnumFormat{Color: color.FgBlue},
		container.ContainerStatusReady:         humanv2.EnumFormat{Color: color.FgGreen},
		container.ContainerStatusUnknownStatus: humanv2.EnumFormat{Color: color.Faint},
		container.ContainerStatusUpdating:      humanv2.EnumFormat{Color: color.FgBlue},
		container.ContainerStatusUpgrading:     humanv2.EnumFormat{Color: color.FgBlue},
	})
	humanv2.RegisterEnumFormat(humanv2.EnumFormats[container.DomainStatus]{
		container.DomainStatusCreating:      humanv2.EnumFormat{Color: color.FgBlue},
		container.DomainStatusDeleting:      humanv2.EnumFormat{Color: color.FgBlue},
		container.DomainStatusError:         humanv2.EnumFormat{Color: color.FgRed},
		container.DomainStatusLocked:        humanv2.EnumFormat{Color: color.FgRed},
		container.DomainStatusLocking:       humanv2.EnumFormat{Color: color.FgBlue},
		container.DomainStatusReady:         humanv2.EnumFormat{Color: color.FgGreen},
		container.DomainStatusUnknownStatus: humanv2.EnumFormat{Color: color.Faint},
		container.DomainStatusUpdating:      humanv2.EnumFormat{Color: color.FgBlue},
		container.DomainStatusUpgrading:     humanv2.EnumFormat{Color: color.FgBlue},
	})
	humanv2.RegisterEnumFormat(humanv2.EnumFormats[container.TriggerStatus]{
		container.TriggerStatusCreating:      humanv2.EnumFormat{Color: color.FgBlue},
		container.TriggerStatusDeleting:      humanv2.EnumFormat{Color: color.FgBlue},
		container.TriggerStatusError:         humanv2.EnumFormat{Color: color.FgRed},
		container.TriggerStatusLocked:        humanv2.EnumFormat{Color: color.FgRed},
		container.TriggerStatusLocking:       humanv2.EnumFormat{Color: color.FgBlue},
		container.TriggerStatusReady:         humanv2.EnumFormat{Color: color.FgGreen},
		container.TriggerStatusUnknownStatus: humanv2.EnumFormat{Color: color.Faint},
		container.TriggerStatusUpdating:      humanv2.EnumFormat{Color: color.FgBlue},
		container.TriggerStatusUpgrading:     humanv2.EnumFormat{Color: color.FgBlue},
	})

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

	cmds.Add(containerCockpitCommand())

	return cmds
}
