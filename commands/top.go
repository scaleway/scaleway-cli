// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/scaleway/scaleway-cli/api"
	types "github.com/scaleway/scaleway-cli/commands/types"
	"github.com/scaleway/scaleway-cli/utils"
)

var cmdTop = &types.Command{
	Exec:        runTop,
	UsageLine:   "top [OPTIONS] SERVER", // FIXME: add ps options
	Description: "Lookup the running processes of a server",
	Help:        "Lookup the running processes of a server.",
}

func init() {
	cmdTop.Flag.BoolVar(&topHelp, []string{"h", "-help"}, false, "Print usage")
	cmdTop.Flag.StringVar(&topGateway, []string{"g", "-gateway"}, "", "Use a SSH gateway")
}

// Flags
var topHelp bool      // -h, --help flag
var topGateway string // -g, --gateway flag

func runTop(cmd *types.Command, args []string) {
	if topHelp {
		cmd.PrintUsage()
	}
	if len(args) != 1 {
		cmd.PrintShortUsage()
	}

	serverID := cmd.API.GetServerID(args[0])
	command := "ps"
	server, err := cmd.API.GetServer(serverID)
	if err != nil {
		log.Fatalf("Failed to get server information for %s: %v", serverID, err)
	}

	// Resolve gateway
	var gateway string
	if topGateway == serverID || topGateway == args[0] {
		gateway = ""
	} else {
		gateway, err = api.ResolveGateway(cmd.API, topGateway)
		if err != nil {
			log.Fatalf("Cannot resolve Gateway '%s': %v", topGateway, err)
		}
	}

	execCmd := utils.NewSSHExecCmd(server.PublicAddress.IP, server.PrivateIP, true, nil, []string{command}, gateway)
	log.Debugf("Executing: ssh %s", strings.Join(execCmd, " "))
	out, err := exec.Command("ssh", execCmd...).CombinedOutput()
	fmt.Printf("%s", out)
	if err != nil {
		log.Fatal(err)
	}
}
