package container

import (
	"context"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/container/v1"
)

var (
	containerDeployTimeout = 12*time.Minute + 30*time.Second

	containerStatusMarshalSpecs = human.EnumMarshalSpecs{
		container.ContainerStatusCreating:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
		container.ContainerStatusDeleting:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
		container.ContainerStatusError:         &human.EnumMarshalSpec{Attribute: color.FgRed},
		container.ContainerStatusLocked:        &human.EnumMarshalSpec{Attribute: color.FgRed},
		container.ContainerStatusLocking:       &human.EnumMarshalSpec{Attribute: color.FgBlue},
		container.ContainerStatusReady:         &human.EnumMarshalSpec{Attribute: color.FgGreen},
		container.ContainerStatusUnknownStatus: &human.EnumMarshalSpec{Attribute: color.Faint},
		container.ContainerStatusUpdating:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
		container.ContainerStatusUpgrading:     &human.EnumMarshalSpec{Attribute: color.FgBlue},
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
