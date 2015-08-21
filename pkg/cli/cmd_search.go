// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import "github.com/scaleway/scaleway-cli/pkg/commands"

var cmdSearch = &Command{
	Exec:        runSearch,
	UsageLine:   "search [OPTIONS] TERM",
	Description: "Search the Scaleway Hub for images",
	Help:        "Search the Scaleway Hub for images.",
}

func init() {
	cmdSearch.Flag.BoolVar(&searchNoTrunc, []string{"-no-trunc"}, false, "Don't truncate output")
	cmdSearch.Flag.BoolVar(&searchHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var searchNoTrunc bool // --no-trunc flag
var searchHelp bool    // -h, --help flag

func runSearch(cmd *Command, rawArgs []string) error {
	if searchHelp {
		return cmd.PrintUsage()
	}
	if len(rawArgs) != 1 {
		return cmd.PrintShortUsage()
	}

	args := commands.SearchArgs{
		Term:    rawArgs[0],
		NoTrunc: searchNoTrunc,
	}
	ctx := cmd.GetContext(rawArgs)
	return commands.RunSearch(ctx, args)
}
