// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import "github.com/scaleway/scaleway-cli/pkg/commands"

var cmdRm = &Command{
	Exec:        runRm,
	UsageLine:   "rm [OPTIONS] SERVER [SERVER...]",
	Description: "Remove one or more servers",
	Help:        "Remove one or more servers.",
	Examples: `
    $ scw rm my-stopped-server my-second-stopped-server
    $ scw rm $(scw ps -q)
    $ scw rm $(scw ps | grep mysql | awk '{print $1}')
`,
}

func init() {
	cmdRm.Flag.BoolVar(&rmHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var rmHelp bool // -h, --help flag

func runRm(cmd *Command, rawArgs []string) error {
	if rmHelp {
		cmd.PrintUsage()
	}
	if len(rawArgs) < 1 {
		cmd.PrintShortUsage()
	}

	args := commands.RmArgs{
		Servers: rawArgs,
	}
	ctx := cmd.GetContext(rawArgs)
	return commands.RunRm(ctx, args)
}
