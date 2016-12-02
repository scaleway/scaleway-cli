// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/utils"
)

// PortArgs are flags for the `RunPort` function
type PortArgs struct {
	Gateway string
	Server  string
	SSHUser string
	SSHPort int
}

// RunPort is the handler for 'scw port'
func RunPort(ctx CommandContext, args PortArgs) error {
	serverID, err := ctx.API.GetServerID(args.Server)
	if err != nil {
		return err
	}
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

	command := []string{"netstat -lutn 2>/dev/null | grep LISTEN"}
	err = utils.SSHExec(server.PublicAddress.IP, server.PrivateIP, args.SSHUser, args.SSHPort, command, true, gateway, false)
	if err != nil {
		return fmt.Errorf("command execution failed: %v", err)
	}

	return nil
}
