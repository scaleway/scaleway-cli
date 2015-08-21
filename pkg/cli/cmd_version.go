// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import "github.com/scaleway/scaleway-cli/pkg/commands"

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

func runVersion(cmd *Command, rawArgs []string) error {
	if versionHelp {
		return cmd.PrintUsage()
	}
	if len(rawArgs) != 0 {
		return cmd.PrintShortUsage()
	}

	args := commands.VersionArgs{}
	ctx := cmd.GetContext(rawArgs)
	return commands.Version(ctx, args)
}
