// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/utils"
)

// LogsArgs are flags for the `RunLogs` function
type LogsArgs struct {
	Gateway string
	Server  string
}

// RunLogs is the handler for 'scw logs'
func RunLogs(ctx CommandContext, args LogsArgs) error {
	serverID := ctx.API.GetServerID(args.Server)
	server, err := ctx.API.GetServer(serverID)
	if err != nil {
		return fmt.Errorf("failed to get server information for %s: %v", serverID, err)
	}

	// FIXME: switch to serial history when API is ready

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

	command := []string{"dmesg"}
	err = utils.SSHExec(server.PublicAddress.IP, server.PrivateIP, command, true, gateway)
	if err != nil {
		return fmt.Errorf("command execution failed: %v", err)
	}
	return nil
}
