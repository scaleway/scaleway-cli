// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"

	types "github.com/scaleway/scaleway-cli/commands/types"
	"github.com/scaleway/scaleway-cli/utils"
)

var cmdKill = &types.Command{
	Exec:        runKill,
	UsageLine:   "kill [OPTIONS] SERVER",
	Description: "Kill a running server",
	Help:        "Kill a running server.",
}

func init() {
	cmdKill.Flag.BoolVar(&killHelp, []string{"h", "-help"}, false, "Print usage")
	cmdKill.Flag.StringVar(&killGateway, []string{"g", "-gateway"}, "", "Use a SSH gateway")
	// FIXME: add --signal option
}

// Flags
var killHelp bool      // -h, --help flag
var killGateway string // -g, --gateway flag

func runKill(cmd *types.Command, args []string) {
	if killHelp {
		cmd.PrintUsage()
	}
	if len(args) < 1 {
		cmd.PrintShortUsage()
	}

	serverID := cmd.API.GetServerID(args[0])
	command := "halt"
	server, err := cmd.API.GetServer(serverID)
	if err != nil {
		log.Fatalf("Failed to get server information for %s: %v", serverID, err)
	}

	execCmd := append(utils.NewSSHExecCmd(server.PublicAddress.IP, server.PrivateIP, true, nil, []string{command}, killGateway))

	log.Debugf("Executing: ssh %s", strings.Join(execCmd, " "))

	spawn := exec.Command("ssh", execCmd...)
	spawn.Stdout = os.Stdout
	spawn.Stdin = os.Stdin
	spawn.Stderr = os.Stderr
	err = spawn.Run()
	if err != nil {
		log.Fatal(err)
	}
}
