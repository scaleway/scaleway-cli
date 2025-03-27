package container

import (
	"context"
	"reflect"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
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

func waitForContainer(ctx context.Context, _, respI interface{}) (interface{}, error) {
	c := respI.(*container.Container)

	client := core.ExtractClient(ctx)
	api := container.NewAPI(client)

	return api.WaitForContainer(&container.WaitForContainerRequest{
		ContainerID:   c.ID,
		Region:        c.Region,
		Timeout:       scw.TimeDurationPtr(containerDeployTimeout),
		RetryInterval: core.DefaultRetryInterval,
	})
}

func containerContainerDeployBuilder(command *core.Command) *core.Command {
	command.WaitFunc = waitForContainer

	return command
}

func containerContainerCreateBuilder(command *core.Command) *core.Command {
	// Add an interceptor that will deploy container after it was created
	type CustomCreateContainerRequest struct {
		*container.CreateContainerRequest
		Deploy bool `json:"deploy"`
	}

	command.ArgSpecs.AddBefore("region", &core.ArgSpec{
		Name:     "deploy",
		Short:    "Deploy container after creation",
		Required: false,
		Default:  core.DefaultValueSetter("true"),
	})

	command.ArgsType = reflect.TypeOf(CustomCreateContainerRequest{})

	command.AddInterceptors(
		func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
			args := argsI.(*CustomCreateContainerRequest)
			resI, err := runner(ctx, args.CreateContainerRequest)
			if err != nil {
				return resI, err
			}

			c := resI.(*container.Container)
			containerAPI := container.NewAPI(core.ExtractClient(ctx))

			return containerAPI.DeployContainer(&container.DeployContainerRequest{
				Region:      c.Region,
				ContainerID: c.ID,
			}, scw.WithContext(ctx))
		},
	)

	command.WaitFunc = waitForContainer

	return command
}

func containerContainerUpdateBuilder(command *core.Command) *core.Command {
	command.WaitFunc = waitForContainer

	return command
}
