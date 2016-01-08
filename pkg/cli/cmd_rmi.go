// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import "github.com/scaleway/scaleway-cli/pkg/commands"

var cmdRmi = &Command{
	Exec:        runRmi,
	UsageLine:   "rmi [OPTIONS] IDENTIFIER [IDENTIFIER...]",
	Description: "Remove one or more image(s)/volume(s)/snapshot(s)",
	Help:        "Remove one or more image(s)/volume(s)/snapshot(s)",
	Examples: `
    $ scw rmi myimage
    $ scw rmi mysnapshot
    $ scw rmi myvolume
    $ scw rmi $(scw images -q)
`,
}

func init() {
	cmdRmi.Flag.BoolVar(&rmiHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var rmiHelp bool // -h, --help flag

func runRmi(cmd *Command, rawArgs []string) error {
	if rmiHelp {
		return cmd.PrintUsage()
	}
	if len(rawArgs) < 1 {
		return cmd.PrintShortUsage()
	}

	args := commands.RmiArgs{
		Identifier: rawArgs,
	}
	ctx := cmd.GetContext(rawArgs)
	return commands.RunRmi(ctx, args)
}
