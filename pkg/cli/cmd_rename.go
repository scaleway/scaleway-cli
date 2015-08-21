// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import "github.com/scaleway/scaleway-cli/pkg/commands"

var cmdRename = &Command{
	Exec:        runRename,
	UsageLine:   "rename [OPTIONS] SERVER NEW_NAME",
	Description: "Rename a server",
	Help:        "Rename a server.",
}

func init() {
	cmdRename.Flag.BoolVar(&renameHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var renameHelp bool // -h, --help flag

func runRename(cmd *Command, rawArgs []string) error {
	if renameHelp {
		return cmd.PrintUsage()
	}
	if len(rawArgs) != 2 {
		return cmd.PrintShortUsage()
	}

	args := commands.RenameArgs{
		Server:  rawArgs[0],
		NewName: rawArgs[1],
	}
	ctx := cmd.GetContext(rawArgs)
	return commands.RunRename(ctx, args)
}
