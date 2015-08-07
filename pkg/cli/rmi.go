// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"github.com/scaleway/scaleway-cli/pkg/commands"
	"github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"
)

var cmdRmi = &Command{
	Exec:        runRmi,
	UsageLine:   "rmi [OPTIONS] IMAGE [IMAGE...]",
	Description: "Remove one or more images",
	Help:        "Remove one or more images.",
	Examples: `
    $ scw rmi myimage
    $ scw rmi $(scw images -q)
`,
}

func init() {
	cmdRmi.Flag.BoolVar(&rmiHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var rmiHelp bool // -h, --help flag

func runRmi(cmd *Command, rawArgs []string) {
	if rmiHelp {
		cmd.PrintUsage()
	}
	if len(rawArgs) < 1 {
		cmd.PrintShortUsage()
	}

	args := commands.RmiArgs{
		Images: rawArgs,
	}
	ctx := cmd.GetContext(rawArgs)
	err := commands.RunRmi(ctx, args)
	if err != nil {
		logrus.Fatalf("Cannot execute 'rmi': %v", err)
	}
}
