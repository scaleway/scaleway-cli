package container

import (
	"context"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/container/v1"
)

var containerDeployTimeout = 12*time.Minute + 30*time.Second

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
