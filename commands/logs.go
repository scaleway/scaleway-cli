// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"os"

	log "github.com/Sirupsen/logrus"

	"github.com/scaleway/scaleway-cli/api"
	types "github.com/scaleway/scaleway-cli/commands/types"
	"github.com/scaleway/scaleway-cli/utils"
)

var cmdLogs = &types.Command{
	Exec:        runLogs,
	UsageLine:   "logs [OPTIONS] SERVER",
	Description: "Fetch the logs of a server",
	Help:        "Fetch the logs of a server.",
}

func init() {
	cmdLogs.Flag.BoolVar(&logsHelp, []string{"h", "-help"}, false, "Print usage")
	cmdLogs.Flag.StringVar(&logsGateway, []string{"g", "-gateway"}, "", "Use a SSH gateway")
}

// FLags
var logsHelp bool      // -h, --help flag
var logsGateway string // -g, --gateway flag

func runLogs(cmd *types.Command, args []string) {
	if logsHelp {
		cmd.PrintUsage()
	}
	if len(args) != 1 {
		cmd.PrintShortUsage()
	}

	serverID := cmd.API.GetServerID(args[0])
	server, err := cmd.API.GetServer(serverID)
	if err != nil {
		log.Fatalf("Failed to get server information for %s: %v", serverID, err)
	}

	// FIXME: switch to serial history when API is ready

	// Resolve gateway
	if logsGateway == "" {
		logsGateway = os.Getenv("SCW_GATEWAY")
	}
	var gateway string
	if logsGateway == serverID || logsGateway == args[0] {
		gateway = ""
	} else {
		gateway, err = api.ResolveGateway(cmd.API, logsGateway)
		if err != nil {
			log.Fatalf("Cannot resolve Gateway '%s': %v", logsGateway, err)
		}
	}

	command := []string{"dmesg"}
	err = utils.SSHExec(server.PublicAddress.IP, server.PrivateIP, command, true, gateway)
	if err != nil {
		log.Fatalf("Command execution failed: %v", err)
	}
}
