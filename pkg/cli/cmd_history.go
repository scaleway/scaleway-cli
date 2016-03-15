// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import "github.com/scaleway/scaleway-cli/pkg/commands"

var cmdHistory = &Command{
	Exec:        runHistory,
	UsageLine:   "history [OPTIONS] IMAGE",
	Description: "Show the history of an image",
	Help:        "Show the history of an image.",
}

func init() {
	cmdHistory.Flag.BoolVar(&historyNoTrunc, []string{"-no-trunc"}, false, "Don't truncate output")
	cmdHistory.Flag.BoolVar(&historyQuiet, []string{"q", "-quiet"}, false, "Only show numeric IDs")
	cmdHistory.Flag.BoolVar(&historyHelp, []string{"h", "-help"}, false, "Print usage")
	cmdHistory.Flag.StringVar(&historyArch, []string{"-arch"}, "*", "Specify architecture")
}

// Flags
var historyNoTrunc bool // --no-trunc flag
var historyQuiet bool   // -q, --quiet flag
var historyHelp bool    // -h, --help flag
var historyArch string  // --arch flag

func runHistory(cmd *Command, rawArgs []string) error {
	if historyHelp {
		return cmd.PrintUsage()
	}
	if len(rawArgs) != 1 {
		return cmd.PrintShortUsage()
	}

	args := commands.HistoryArgs{
		Quiet:   historyQuiet,
		NoTrunc: historyNoTrunc,
		Image:   rawArgs[0],
		Arch:    historyArch,
	}
	ctx := cmd.GetContext(rawArgs)
	return commands.RunHistory(ctx, args)
}
