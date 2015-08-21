// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import "github.com/scaleway/scaleway-cli/pkg/commands"

var cmdAttach = &Command{
	Exec:        runAttach,
	UsageLine:   "attach [OPTIONS] SERVER",
	Description: "Attach to a server serial console",
	Help:        "Attach to a running server serial console.",
	Examples: `
    $ scw attach my-running-server
    $ scw attach $(scw start my-stopped-server)
    $ scw attach $(scw start $(scw create ubuntu-vivid))
`,
}

func init() {
	cmdAttach.Flag.BoolVar(&attachHelp, []string{"h", "-help"}, false, "Print usage")
	cmdAttach.Flag.BoolVar(&attachNoStdin, []string{"-no-stdin"}, false, "Do not attach stdin")
}

// Flags
var attachHelp bool    // -h, --help flag
var attachNoStdin bool // --no-stdin flag

func runAttach(cmd *Command, rawArgs []string) error {
	if attachHelp {
		return cmd.PrintUsage()
	}
	if len(rawArgs) != 1 {
		return cmd.PrintShortUsage()
	}

	args := commands.AttachArgs{
		NoStdin: attachNoStdin,
		Server:  rawArgs[0],
	}
	ctx := cmd.GetContext(rawArgs)
	return commands.RunAttach(ctx, args)
}
