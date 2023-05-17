package secret

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	secret "github.com/scaleway/scaleway-sdk-go/api/secret/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("secret", "version", "create").Override(dataCreateVersion)
	cmds.MustFind("secret", "version", "access").Override(accessVersion)
	return cmds
}

func dataCreateVersion(c *core.Command) *core.Command {
	*c.ArgSpecs.GetByName("data") = core.ArgSpec{
		Name:        "data",
		Short:       "Content of the secret version. Base64 is handled by the SDK",
		Required:    true,
		CanLoadFile: true,
	}
	return c
}

type accessSecretVersionRequest struct {
	secret.AccessSecretVersionRequest
	OnlyData bool
}

func accessVersion(command *core.Command) *core.Command {
	command.ArgsType = reflect.TypeOf(accessSecretVersionRequest{})
	command.ArgSpecs = core.ArgSpecs{
		{
			Name:       "secret-id",
			Short:      `ID of the secret`,
			Required:   true,
			Deprecated: false,
			Positional: false,
		},
		{
			Name:       "revision",
			Short:      `Version number`,
			Required:   true,
			Deprecated: false,
			Positional: false,
		},
		{
			Name:       "only-data",
			Short:      `Retrieve only data`,
			Required:   false,
			Deprecated: false,
			Positional: false,
			Default:    core.DefaultValueSetter("false"),
		},
		core.RegionArgSpec(scw.RegionFrPar),
	}

	command.Run = func(ctx context.Context, args interface{}) (i interface{}, e error) {
		argsAccessSecretVersion := args.(*accessSecretVersionRequest)

		request := &secret.AccessSecretVersionRequest{
			SecretID: argsAccessSecretVersion.SecretID,
			Revision: argsAccessSecretVersion.Revision,
		}
		client := core.ExtractClient(ctx)
		api := secret.NewAPI(client)
		resp, err := api.AccessSecretVersion(request)
		if err != nil {
			return nil, err
		}
		if argsAccessSecretVersion.OnlyData {
			return resp.Data, nil
		}

		return resp, err
	}
	return command
}
