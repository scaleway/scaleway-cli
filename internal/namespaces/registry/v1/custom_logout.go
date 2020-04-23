package registry

import (
	"context"
	"fmt"
	"os/exec"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type registryLogoutArgs struct {
	Region  scw.Region
	Program string
}

func registryLogoutCommand() *core.Command {
	return &core.Command{
		Short: `Logout of a registry`,
		Long: `This command will run the correct command in order to log you out of the registry with the chosen program.
You will need to have the chosen binary installed on your system and in your PATH.`,
		Namespace: "registry",
		Resource:  "logout",
		ArgsType:  reflect.TypeOf(registryLogoutArgs{}),
		ArgSpecs: []*core.ArgSpec{
			{
				Name:       "program",
				Short:      "Program used to log in to the namespace",
				Default:    core.DefaultValueSetter(string(docker)),
				EnumValues: availablePrograms.StringArray(),
			},
			core.RegionArgSpec(),
		},
		Run: registryLogoutRun,
	}
}

func registryLogoutRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	args := argsI.(*registryLogoutArgs)
	client := core.ExtractClient(ctx)

	region := args.Region.String()
	if region == "" {
		scwRegion, ok := client.GetDefaultRegion()
		if !ok {
			return nil, fmt.Errorf("no default region configured")
		}
		region = scwRegion.String()
	}
	endpoint := endpointPrefix + region + endpointSuffix

	var commandArgs []string

	switch program(args.Program) {
	case docker, podman:
		commandArgs = []string{"logout", endpoint}
	default:
		return nil, fmt.Errorf("unknown program")
	}

	logoutCommand := core.ExtractAndExecCommand(ctx, args.Program, commandArgs...)

	if err := logoutCommand.Run(); err != nil {
		if execErr, ok := err.(*exec.ExitError); ok {
			return nil, &core.CliError{Empty: true, Code: execErr.ExitCode()}
		}
		return nil, fmt.Errorf("could not logout of namespace: %w", err)
	}

	return &core.SuccessResult{
		Empty: true, // the program will output the sucess message
	}, nil
}
