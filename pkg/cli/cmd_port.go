// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import "github.com/scaleway/scaleway-cli/pkg/commands"

var cmdPort = &Command{
	Exec:        runPort,
	UsageLine:   "port [OPTIONS] SERVER [PRIVATE_PORT[/PROTO]]",
	Description: "Lookup the public-facing port that is NAT-ed to PRIVATE_PORT",
	Help:        "List port mappings for the SERVER, or lookup the public-facing port that is NAT-ed to the PRIVATE_PORT",
}

func init() {
	cmdPort.Flag.BoolVar(&portHelp, []string{"h", "-help"}, false, "Print usage")
	cmdPort.Flag.StringVar(&portGateway, []string{"g", "-gateway"}, "", "Use a SSH gateway")
	cmdPort.Flag.StringVar(&portSSHUser, []string{"-user"}, "root", "Specify SSH user")
	cmdPort.Flag.IntVar(&portSSHPort, []string{"p", "-port"}, 22, "Specify SSH port")
}

// FLags
var portHelp bool      // -h, --help flag
var portGateway string // -g, --gateway flag
var portSSHUser string // --user flag
var portSSHPort int    // -p, --port flag

func runPort(cmd *Command, rawArgs []string) error {
	if portHelp {
		return cmd.PrintUsage()
	}
	if len(rawArgs) < 1 {
		return cmd.PrintShortUsage()
	}

	args := commands.PortArgs{
		Gateway: portGateway,
		Server:  rawArgs[0],
		SSHUser: portSSHUser,
		SSHPort: portSSHPort,
	}
	ctx := cmd.GetContext(rawArgs)
	return commands.RunPort(ctx, args)
}
