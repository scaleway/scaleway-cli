package baremetal

import "github.com/scaleway/scaleway-cli/internal/core"

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.Merge(core.NewCommands(
		serverWaitCommand(),
		serverStartCommand(),
		serverStopCommand(),
		serverRebootCommand(),
	))

	cmds.MustFind("baremetal", "server", "create").Override(serverCreateBuilder)

	return cmds
}
