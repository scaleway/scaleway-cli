// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import "github.com/scaleway/scaleway-cli/pkg/commands"

var cmdStop = &Command{
	Exec:        runStop,
	UsageLine:   "stop [OPTIONS] SERVER [SERVER...]",
	Description: "Stop a running server",
	Help:        "Stop a running server.",
	Examples: `
    $ scw stop my-running-server my-second-running-server
    $ scw stop -t my-running-server my-second-running-server
    $ scw stop $(scw ps -q)
    $ scw stop $(scw ps | grep mysql | awk '{print $1}')
    $ scw stop server && stop wait server
    $ scw stop -w server
`,
}

func init() {
	cmdStop.Flag.BoolVar(&stopT, []string{"t", "-terminate"}, false, "Stop and trash a server with its volumes")
	cmdStop.Flag.BoolVar(&stopHelp, []string{"h", "-help"}, false, "Print usage")
	cmdStop.Flag.BoolVar(&stopW, []string{"w", "-wait"}, false, "Synchronous stop. Wait for SSH to be ready")
}

// Flags
var stopT bool    // -t flag
var stopHelp bool // -h, --help flag
var stopW bool    // -w, --wait flat

func runStop(cmd *Command, rawArgs []string) error {
	if stopHelp {
		return cmd.PrintUsage()
	}
	if len(rawArgs) < 1 {
		return cmd.PrintShortUsage()
	}

	args := commands.StopArgs{
		Terminate: stopT,
		Wait:      stopW,
		Servers:   rawArgs,
	}
	ctx := cmd.GetContext(rawArgs)
	return commands.RunStop(ctx, args)
}
