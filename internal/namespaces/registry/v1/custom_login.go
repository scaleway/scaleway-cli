package registry

import (
	"context"
	"fmt"
	"os/exec"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type registryLoginArgs struct {
	Region  scw.Region
	Program string
}

func registryLoginCommand() *core.Command {
	return &core.Command{
		Short: `Login to a registry`,
		Long: `This command will run the correct command in order to log you in on the registry with the chosen program.
You will need to have the chosen binary installed on your system and in your PATH.`,
		Namespace: "registry",
		Resource:  "login",
		ArgsType:  reflect.TypeOf(registryLoginArgs{}),
		ArgSpecs: []*core.ArgSpec{
			{
				Name:       "program",
				Short:      "Program used to log in to the namespace",
				Default:    core.DefaultValueSetter(string(docker)),
				EnumValues: availablePrograms.StringArray(),
			},
			core.RegionArgSpec(scw.AllRegions...),
		},
		Run: registryLoginRun,
	}
}

func registryLoginRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	args := argsI.(*registryLoginArgs)
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

	secretKey, ok := client.GetSecretKey()
	if !ok {
		return nil, fmt.Errorf("could not get secret key")
	}

	var commandArgs []string

	switch program(args.Program) {
	case docker, podman:
		commandArgs = []string{"login", "-u", "scaleway", "--password-stdin", endpoint}
	default:
		return nil, fmt.Errorf("unknown program")
	}

	loginCommand := core.ExtractAndExecCommand(ctx, args.Program, commandArgs...)

	loginCommand.SetStdin(strings.NewReader(secretKey))

	if err := loginCommand.Run(); err != nil {
		if execErr, ok := err.(*exec.ExitError); ok {
			return nil, &core.CliError{Empty: true, Code: execErr.ExitCode()}
		}
		return nil, fmt.Errorf("could not login to namespace: %w", err)
	}

	return &core.SuccessResult{
		Empty: true, // the program will output the sucess message
	}, nil
}
