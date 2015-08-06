// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/machine/log"
	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/utils"
)

// KillArgs are flags for the `RunKill` function
type KillArgs struct {
	Gateway string
	Server  string
}

// RunKill is the handler for 'scw kill'
func RunKill(ctx CommandContext, args KillArgs) error {
	serverID := ctx.API.GetServerID(args.Server)
	command := "halt"
	server, err := ctx.API.GetServer(serverID)
	if err != nil {
		return fmt.Errorf("failed to get server information for %s: %v", serverID, err)
	}

	// Resolve gateway
	if args.Gateway == "" {
		args.Gateway = ctx.Getenv("SCW_GATEWAY")
	}
	var gateway string
	if args.Gateway == serverID || args.Gateway == args.Server {
		gateway = ""
	} else {
		gateway, err = api.ResolveGateway(ctx.API, args.Gateway)
		if err != nil {
			return fmt.Errorf("cannot resolve Gateway '%s': %v", args.Gateway, err)
		}
	}

	execCmd := append(utils.NewSSHExecCmd(server.PublicAddress.IP, server.PrivateIP, true, nil, []string{command}, gateway))

	log.Debugf("Executing: ssh %s", strings.Join(execCmd, " "))

	spawn := exec.Command("ssh", execCmd...)
	spawn.Stdout = os.Stdout
	spawn.Stdin = os.Stdin
	spawn.Stderr = os.Stderr

	return spawn.Run()
}
