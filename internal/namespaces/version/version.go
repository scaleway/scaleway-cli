package version

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func GetCommands() *core.Commands {
	return core.NewCommands(versionCommand())
}

func versionCommand() *core.Command {
	return &core.Command{
		Groups:               []string{"utility"},
		Short:                `Display cli version`,
		Namespace:            "version",
		AllowAnonymousClient: true,
		ArgsType:             reflect.TypeOf(struct{}{}),
		ArgSpecs:             core.ArgSpecs{},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			return core.ExtractBuildInfo(ctx), nil
		},
	}
}
