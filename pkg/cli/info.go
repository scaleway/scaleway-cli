// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"github.com/scaleway/scaleway-cli/pkg/commands"
	"github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"
)

var cmdInfo = &Command{
	Exec:        runInfo,
	UsageLine:   "info [OPTIONS]",
	Description: "Display system-wide information",
	Help:        "Display system-wide information.",
}

func init() {
	cmdInfo.Flag.BoolVar(&infoHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var infoHelp bool // -h, --help flag

func runInfo(cmd *Command, rawArgs []string) {
	if infoHelp {
		cmd.PrintUsage()
	}
	if len(rawArgs) != 0 {
		cmd.PrintShortUsage()
	}

	args := commands.InfoArgs{}
	ctx := cmd.GetContext(rawArgs)
	err := commands.RunInfo(ctx, args)
	if err != nil {
		logrus.Fatalf("Cannot execute 'info': %v", err)
	}
}
