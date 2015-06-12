package commands

import types "github.com/scaleway/scaleway-cli/commands/types"

var Commands = []*types.Command{
	CmdHelp,

	cmdAttach,
	cmdCommit,
	cmdCompletion,
	cmdCp,
	cmdCreate,
	cmdEvents,
	cmdExec,
	cmdHistory,
	cmdImages,
	cmdInfo,
	cmdInspect,
	cmdKill,
	cmdLogin,
	cmdLogout,
	cmdLogs,
	cmdPatch,
	cmdPort,
	cmdPs,
	cmdRename,
	cmdRestart,
	cmdRm,
	cmdRmi,
	cmdRun,
	cmdSearch,
	cmdStart,
	cmdStop,
	cmdTag,
	cmdTop,
	cmdVersion,
	cmdWait,
}
