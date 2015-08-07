// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"github.com/scaleway/scaleway-cli/pkg/commands"
	"github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"
)

var cmdVersion = &Command{
	Exec:        runVersion,
	UsageLine:   "version [OPTIONS]",
	Description: "Show the version information",
	Help:        "Show the version information.",
}

func init() {
	cmdVersion.Flag.BoolVar(&versionHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var versionHelp bool // -h, --help flag

func runVersion(cmd *Command, rawArgs []string) {
	if versionHelp {
		cmd.PrintUsage()
	}
	if len(rawArgs) != 0 {
		cmd.PrintShortUsage()
	}

	args := commands.VersionArgs{}
	ctx := cmd.GetContext(rawArgs)
	err := commands.Version(ctx, args)
	if err != nil {
		logrus.Fatalf("Cannot execute 'version': %v", err)
	}
}
