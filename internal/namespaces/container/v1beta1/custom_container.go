package container

import (
	"context"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	container "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var (
	containerDeployTimeout = 12*time.Minute + 30*time.Second

	containerStatusMarshalSpecs = human.EnumMarshalSpecs{
		container.ContainerStatusCreated:  &human.EnumMarshalSpec{Attribute: color.FgGreen},
		container.ContainerStatusCreating: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		container.ContainerStatusDeleting: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		container.ContainerStatusError:    &human.EnumMarshalSpec{Attribute: color.FgRed},
		container.ContainerStatusLocked:   &human.EnumMarshalSpec{Attribute: color.FgRed},
		container.ContainerStatusPending:  &human.EnumMarshalSpec{Attribute: color.FgBlue},
		container.ContainerStatusReady:    &human.EnumMarshalSpec{Attribute: color.FgGreen},
		container.ContainerStatusUnknown:  &human.EnumMarshalSpec{Attribute: color.Faint},
	}
)

func containerContainerDeployBuilder(command *core.Command) *core.Command {
	command.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		req := argsI.(*container.DeployContainerRequest)

		client := core.ExtractClient(ctx)
		api := container.NewAPI(client)
		return api.WaitForContainer(&container.WaitForContainerRequest{
			ContainerID:   req.ContainerID,
			Region:        req.Region,
			Timeout:       scw.TimeDurationPtr(containerDeployTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}
	return command
}
