// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import "github.com/scaleway/scaleway-cli/pkg/commands"

var cmdKill = &Command{
	Exec:        runKill,
	UsageLine:   "kill [OPTIONS] SERVER",
	Description: "Kill a running server",
	Help:        "Kill a running server.",
}

func init() {
	cmdKill.Flag.BoolVar(&killHelp, []string{"h", "-help"}, false, "Print usage")
	cmdKill.Flag.StringVar(&killGateway, []string{"g", "-gateway"}, "", "Use a SSH gateway")
	cmdKill.Flag.StringVar(&killSSHUser, []string{"u", "-user"}, "root", "Specify SSH user")
	cmdKill.Flag.IntVar(&killSSHPort, []string{"p", "-port"}, 22, "Specify SSH port")
	// FIXME: add --signal option
}

// Flags
var killHelp bool      // -h, --help flag
var killGateway string // -g, --gateway flag
var killSSHUser string // -u, --user flag
var killSSHPort int    // -p, --port flag

func runKill(cmd *Command, rawArgs []string) error {
	if killHelp {
		return cmd.PrintUsage()
	}
	if len(rawArgs) < 1 {
		return cmd.PrintShortUsage()
	}

	args := commands.KillArgs{
		Gateway: killGateway,
		Server:  rawArgs[0],
		SSHUser: killSSHUser,
		SSHPort: killSSHPort,
	}
	ctx := cmd.GetContext(rawArgs)
	return commands.RunKill(ctx, args)
}
