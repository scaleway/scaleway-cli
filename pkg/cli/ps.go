// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"github.com/Sirupsen/logrus"
	"github.com/scaleway/scaleway-cli/pkg/commands"
)

var cmdPs = &Command{
	Exec:        cmdExecPs,
	UsageLine:   "ps [OPTIONS]",
	Description: "List servers",
	Help:        "List servers. By default, only running servers are displayed.",
}

func init() {
	cmdPs.Flag.BoolVar(&psA, []string{"a", "-all"}, false, "Show all servers. Only running servers are shown by default")
	cmdPs.Flag.BoolVar(&psL, []string{"l", "-latest"}, false, "Show only the latest created server, include non-running ones")
	cmdPs.Flag.IntVar(&psN, []string{"n"}, 0, "Show n last created servers, include non-running ones")
	cmdPs.Flag.BoolVar(&psNoTrunc, []string{"-no-trunc"}, false, "Don't truncate output")
	cmdPs.Flag.BoolVar(&psQ, []string{"q", "-quiet"}, false, "Only display numeric IDs")
	cmdPs.Flag.BoolVar(&psHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var psA bool       // -a flag
var psL bool       // -l flag
var psQ bool       // -q flag
var psNoTrunc bool // -no-trunc flag
var psN int        // -n flag
var psHelp bool    // -h, --help flag

func cmdExecPs(cmd *Command, rawArgs []string) {
	if psHelp {
		cmd.PrintUsage()
	}
	if len(rawArgs) != 0 {
		cmd.PrintShortUsage()
	}

	args := commands.PsArgs{
		All:     psA,
		Latest:  psL,
		Quiet:   psQ,
		NoTrunc: psNoTrunc,
		NLast:   psN,
	}
	ctx := cmd.GetContext(rawArgs)
	err := commands.RunPs(ctx, args)
	if err != nil {
		logrus.Fatalf("Cannot execute 'ps': %v", err)
	}
}
