// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

// Docker-style commands to manage BareMetal servers

// Package cli contains CLI commands
package cli

// Commands is the list of enabled CLI commands
var Commands = []*Command{
	CmdHelp,

	cmdAttach,
	cmdCommit,
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
	cmdPort,
	cmdProducts,
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
	cmdUserdata,
	cmdVersion,
	cmdWait,

	cmdBilling,
	cmdCompletion,
	cmdFlushCache,
	cmdMarketplace,
	cmdPatch,
	cmdSecurityGroups,
	cmdIPS,
	cmdCS,
}
