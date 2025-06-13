package registry

import (
	"context"
	"os/exec"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
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

func registryLogoutRun(ctx context.Context, argsI any) (i any, e error) {
	args := argsI.(*registryLogoutArgs)

	region := args.Region.String()
	endpoint := endpointPrefix + region + endpointSuffix

	cmdArgs := []string{"logout", endpoint}
	cmd := exec.Command(args.Program, cmdArgs...) //nolint:gosec
	exitCode, err := core.ExecCmd(ctx, cmd)
	if err != nil {
		return nil, err
	}
	if exitCode != 0 {
		return nil, &core.CliError{Empty: true, Code: exitCode}
	}

	return &core.SuccessResult{
		Empty: true, // the program will output the success message
	}, nil
}
