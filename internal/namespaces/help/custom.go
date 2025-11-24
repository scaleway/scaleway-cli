package help

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
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
		Short:                "Get help about how specific topics inside the CLI work",
		Namespace:            "help",
		AllowAnonymousClient: true,
		Groups:               []string{"utility"},
		SeeAlsos: []*core.SeeAlso{
			{
				Short:   "Print general help",
				Command: "scw --help",
			},
		},
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
		Run: func(_ context.Context, _ any) (any, error) {
			return long, nil
		},
	}
}
