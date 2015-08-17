// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import "github.com/scaleway/scaleway-cli/pkg/commands"

var cmdStart = &Command{
	Exec:        runStart,
	UsageLine:   "start [OPTIONS] SERVER [SERVER...]",
	Description: "Start a stopped server",
	Help:        "Start a stopped server.",
}

func init() {
	cmdStart.Flag.BoolVar(&startW, []string{"w", "-wait"}, false, "Synchronous start. Wait for SSH to be ready")
	cmdStart.Flag.Float64Var(&startTimeout, []string{"T", "-timeout"}, 0, "Set timeout values to seconds")
	cmdStart.Flag.BoolVar(&startHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var startW bool          // -w flag
var startTimeout float64 // -T flag
var startHelp bool       // -h, --help flag

func runStart(cmd *Command, rawArgs []string) error {
	if startHelp {
		cmd.PrintUsage()
	}
	if len(rawArgs) < 1 {
		cmd.PrintShortUsage()
	}

	args := commands.StartArgs{
		Servers: rawArgs,
		Timeout: startTimeout,
		Wait:    startW,
	}
	ctx := cmd.GetContext(rawArgs)
	return commands.RunStart(ctx, args)
}
