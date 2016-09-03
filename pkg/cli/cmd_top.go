// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import "github.com/scaleway/scaleway-cli/pkg/commands"

var cmdTop = &Command{
	Exec:        runTop,
	UsageLine:   "top [OPTIONS] SERVER", // FIXME: add ps options
	Description: "Lookup the running processes of a server",
	Help:        "Lookup the running processes of a server.",
}

func init() {
	cmdTop.Flag.BoolVar(&topHelp, []string{"h", "-help"}, false, "Print usage")
	cmdTop.Flag.StringVar(&topGateway, []string{"g", "-gateway"}, "", "Use a SSH gateway")
	cmdTop.Flag.StringVar(&topSSHUser, []string{"u", "-user"}, "root", "Specify SSH user")
}

// Flags
var topHelp bool      // -h, --help flag
var topGateway string // -g, --gateway flag
var topSSHUser string // -u, --user flag

func runTop(cmd *Command, rawArgs []string) error {
	if topHelp {
		return cmd.PrintUsage()
	}
	if len(rawArgs) != 1 {
		return cmd.PrintShortUsage()
	}

	args := commands.TopArgs{
		Gateway: topGateway,
		Server:  rawArgs[0],
		SSHUser: topSSHUser,
	}
	ctx := cmd.GetContext(rawArgs)
	return commands.RunTop(ctx, args)
}
