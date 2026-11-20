package container

import (
	"context"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/container/v1"
)

var (
	containerDeployTimeout = 12*time.Minute + 30*time.Second

	containerStatusFormat = humanv2.EnumFormats[container.ContainerStatus]{
		container.ContainerStatusCreating:      humanv2.EnumFormat{Color: color.FgBlue},
		container.ContainerStatusDeleting:      humanv2.EnumFormat{Color: color.FgBlue},
		container.ContainerStatusError:         humanv2.EnumFormat{Color: color.FgRed},
		container.ContainerStatusLocked:        humanv2.EnumFormat{Color: color.FgRed},
		container.ContainerStatusLocking:       humanv2.EnumFormat{Color: color.FgBlue},
		container.ContainerStatusReady:         humanv2.EnumFormat{Color: color.FgGreen},
		container.ContainerStatusUnknownStatus: humanv2.EnumFormat{Color: color.Faint},
		container.ContainerStatusUpdating:      humanv2.EnumFormat{Color: color.FgBlue},
		container.ContainerStatusUpgrading:     humanv2.EnumFormat{Color: color.FgBlue},
	}
)

func waitForContainer(ctx context.Context, _, respI any) (any, error) {
	c := respI.(*container.Container)

	client := core.ExtractClient(ctx)
	api := container.NewAPI(client)

	return api.WaitForContainer(&container.WaitForContainerRequest{
		ContainerID:   c.ID,
		Region:        c.Region,
		Timeout:       new(containerDeployTimeout),
		RetryInterval: core.DefaultRetryInterval,
	})
}

func containerContainerRedeployBuilder(command *core.Command) *core.Command {
	command.WaitFunc = waitForContainer

	return command
}

func containerContainerCreateBuilder(command *core.Command) *core.Command {
	command.WaitFunc = waitForContainer

	return command
}

func containerContainerUpdateBuilder(command *core.Command) *core.Command {
	command.WaitFunc = waitForContainer

	return command
}
