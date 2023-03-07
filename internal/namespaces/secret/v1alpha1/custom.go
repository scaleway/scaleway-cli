package secret

import (
	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("secret", "version", "create").Override(dataCreateVersion)
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
