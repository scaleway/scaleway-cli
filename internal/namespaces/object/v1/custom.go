package object

import (
	"github.com/scaleway/scaleway-cli/internal/core"
)

func GetCommands() *core.Commands {
	return core.NewCommands(
		objectRoot(),
		objectConfig(),
		configGetCommand(),
		configInstallCommand(),
	)
}

func objectRoot() *core.Command {
	return &core.Command{
		Short:     `Object-storage utils`,
		Namespace: "object",
	}
}

func objectConfig() *core.Command {
	return &core.Command{
		Short:     `Manage configuration files for popular S3 tools`,
		Long:      `Configuration generation for S3 tools.`,
		Namespace: "object",
		Resource:  `config`,
	}
}
