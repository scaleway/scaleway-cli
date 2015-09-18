// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import "github.com/scaleway/scaleway-cli/pkg/commands"

var cmdRestart = &Command{
	Exec:        runRestart,
	UsageLine:   "restart [OPTIONS] SERVER [SERVER...]",
	Description: "Restart a running server",
	Help:        "Restart a running server.",
}

func init() {
	cmdRestart.Flag.BoolVar(&restartW, []string{"w", "-wait"}, false, "Synchronous restart. Wait for SSH to be ready")
	cmdRestart.Flag.Float64Var(&restartTimeout, []string{"T", "-timeout"}, 0, "Set timeout values to seconds")
	cmdRestart.Flag.BoolVar(&restartHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var restartW bool          // -w flag
var restartTimeout float64 // -T flag
var restartHelp bool       // -h, --help flag

func runRestart(cmd *Command, rawArgs []string) error {
	if restartHelp {
		return cmd.PrintUsage()
	}
	if len(rawArgs) < 1 {
		return cmd.PrintShortUsage()
	}

	args := commands.RestartArgs{
		Timeout: restartTimeout,
		Wait:    restartW,
		Servers: rawArgs,
	}
	ctx := cmd.GetContext(rawArgs)
	return commands.RunRestart(ctx, args)
}
