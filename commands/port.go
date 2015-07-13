// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	log "github.com/Sirupsen/logrus"

	types "github.com/scaleway/scaleway-cli/commands/types"
	"github.com/scaleway/scaleway-cli/utils"
)

var cmdPort = &types.Command{
	Exec:        runPort,
	UsageLine:   "port [OPTIONS] SERVER [PRIVATE_PORT[/PROTO]]",
	Description: "Lookup the public-facing port that is NAT-ed to PRIVATE_PORT",
	Help:        "List port mappings for the SERVER, or lookup the public-facing port that is NAT-ed to the PRIVATE_PORT",
}

func init() {
	cmdPort.Flag.BoolVar(&portHelp, []string{"h", "-help"}, false, "Print usage")
	cmdPort.Flag.StringVar(&portGateway, []string{"g", "-gateway"}, "", "Use a SSH gateway")
}

// FLags
var portHelp bool      // -h, --help flag
var portGateway string // -g, --gateway flag

func runPort(cmd *types.Command, args []string) {
	if portHelp {
		cmd.PrintUsage()
	}
	if len(args) < 1 {
		cmd.PrintShortUsage()
	}

	serverID := cmd.API.GetServerID(args[0])
	server, err := cmd.API.GetServer(serverID)
	if err != nil {
		log.Fatalf("Failed to get server information for %s: %v", serverID, err)
	}

	command := []string{"netstat -lutn 2>/dev/null | grep LISTEN"}
	err = utils.SSHExec(server.PublicAddress.IP, server.PrivateIP, command, true, portGateway)
	if err != nil {
		log.Fatalf("Command execution failed: %v", err)
	}
}
