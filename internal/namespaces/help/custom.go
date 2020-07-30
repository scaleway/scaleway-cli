package help

import (
	"github.com/scaleway/scaleway-cli/internal/core"
)

func GetCommands() *core.Commands {
	return core.NewCommands(
		helpRoot(),
		helpOutput(),
	)
}

func helpRoot() *core.Command {
	return &core.Command{
		Short:                "Get help about how the CLI works",
		Namespace:            "help",
		AllowAnonymousClient: true,
	}
}
