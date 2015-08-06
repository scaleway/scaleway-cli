// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"os"
	"os/exec"
	"strings"

	log "github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"

	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/utils"
)

var cmdKill = &Command{
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

func runKill(cmd *Command, args []string) {
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

	// Resolve gateway
	if killGateway == "" {
		killGateway = os.Getenv("SCW_GATEWAY")
	}
	var gateway string
	if killGateway == serverID || killGateway == args[0] {
		gateway = ""
	} else {
		gateway, err = api.ResolveGateway(cmd.API, killGateway)
		if err != nil {
			log.Fatalf("Cannot resolve Gateway '%s': %v", killGateway, err)
		}
	}

	execCmd := append(utils.NewSSHExecCmd(server.PublicAddress.IP, server.PrivateIP, true, nil, []string{command}, gateway))

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
