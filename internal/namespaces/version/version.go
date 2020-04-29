package version

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func GetCommands() *core.Commands {
	return core.NewCommands(versionCommand())
}

func versionCommand() *core.Command {
	return &core.Command{
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
