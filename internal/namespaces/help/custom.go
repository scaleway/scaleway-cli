package help

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func GetCommands() *core.Commands {
	return core.NewCommands(
		helpRoot(),
		newHelpCommand("output", shortOutput, longOutput),
		newHelpCommand("date", shortDate, longDate),
	)
}

func helpRoot() *core.Command {
	return &core.Command{
		Short:                "Get help about how the CLI works",
		Namespace:            "help",
		AllowAnonymousClient: true,
	}
}

func newHelpCommand(resource string, short string, long string) *core.Command {
	return &core.Command{
		Short:                short,
		Long:                 long,
		Namespace:            "help",
		Resource:             resource,
		AllowAnonymousClient: true,
		ArgsType:             reflect.TypeOf(struct{}{}),
		ArgSpecs:             core.ArgSpecs{},
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			return long, nil
		},
	}
}
