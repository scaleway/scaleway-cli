package secret

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	secret "github.com/scaleway/scaleway-sdk-go/api/secret/v1beta1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("secret", "version", "create").Override(dataCreateVersion)
	cmds.MustFind("secret", "version", "access").Override(secretVersionAccessCommand)
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

func secretVersionAccessCommand(c *core.Command) *core.Command {
	c.Run = func(ctx context.Context, args interface{}) (i interface{}, e error) {
		request := args.(*secret.AccessSecretVersionRequest)

		client := core.ExtractClient(ctx)
		api := secret.NewAPI(client)
		res, err := api.AccessSecretVersion(request)

		if err != nil {
			return nil, err
		}

		return core.RawResult(res.Data), nil
	}

	return c
}
