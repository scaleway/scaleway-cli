// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import "github.com/scaleway/scaleway-cli/pkg/commands"

var cmdWait = &Command{
	Exec:        runWait,
	UsageLine:   "wait [OPTIONS] SERVER [SERVER...]",
	Description: "Block until a server stops",
	Help:        "Block until a server stops.",
}

func init() {
	cmdWait.Flag.BoolVar(&waitHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var waitHelp bool // -h, --help flag

func runWait(cmd *Command, rawArgs []string) error {
	if waitHelp {
		return cmd.PrintUsage()
	}
	if len(rawArgs) < 1 {
		return cmd.PrintShortUsage()
	}

	args := commands.WaitArgs{
		Servers: rawArgs,
	}
	ctx := cmd.GetContext(rawArgs)
	return commands.RunWait(ctx, args)
}
