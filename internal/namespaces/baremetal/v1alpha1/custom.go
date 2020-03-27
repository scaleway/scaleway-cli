package baremetal

import "github.com/scaleway/scaleway-cli/internal/core"

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.Merge(core.NewCommands(
		serverWaitCommand(),
	))

	cmds.MustFind("baremetal", "server", "create").Override(serverCreateBuilder)

	// Action commands
	cmds.MustFind("baremetal", "server", "start").Override(serverStartBuilder)
	cmds.MustFind("baremetal", "server", "stop").Override(serverStopBuilder)
	cmds.MustFind("baremetal", "server", "reboot").Override(serverRebootBuilder)

	return cmds
}
