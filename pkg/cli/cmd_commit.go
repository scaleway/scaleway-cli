// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import "github.com/scaleway/scaleway-cli/pkg/commands"

var cmdCommit = &Command{
	Exec:        runCommit,
	UsageLine:   "commit [OPTIONS] SERVER [NAME]",
	Description: "Create a new snapshot from a server's volume",
	Help:        "Create a new snapshot from a server's volume.",
	Examples: `
    $ scw commit my-stopped-server
    $ scw commit -v 1 my-stopped-server
`,
}

func init() {
	cmdCommit.Flag.IntVar(&commitVolume, []string{"v", "-volume"}, 0, "Volume slot")
	cmdCommit.Flag.BoolVar(&commitHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var commitVolume int // -v, --volume flag
var commitHelp bool  // -h, --help flag

func runCommit(cmd *Command, rawArgs []string) error {
	if commitHelp {
		return cmd.PrintUsage()
	}
	if len(rawArgs) < 1 {
		return cmd.PrintShortUsage()
	}

	args := commands.CommitArgs{
		Volume: commitVolume,
		Server: rawArgs[0],
		Name:   "",
	}
	if len(rawArgs) > 1 {
		args.Name = rawArgs[1]
	}

	ctx := cmd.GetContext(rawArgs)
	return commands.RunCommit(ctx, args)
}
