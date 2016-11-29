// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"os/exec"

	"github.com/Sirupsen/logrus"
	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/utils"
)

// TopArgs are flags for the `RunTop` function
type TopArgs struct {
	Server  string
	Gateway string
	SSHUser string
	SSHPort int
}

// RunTop is the handler for 'scw top'
func RunTop(ctx CommandContext, args TopArgs) error {
	serverID, err := ctx.API.GetServerID(args.Server)
	if err != nil {
		return err
	}
	command := "ps"
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

	sshCommand := utils.NewSSHExecCmd(server.PublicAddress.IP, server.PrivateIP, args.SSHUser, args.SSHPort, true, []string{command}, gateway, false)
	logrus.Debugf("Executing: %s", sshCommand)
	out, err := exec.Command("ssh", sshCommand.Slice()[1:]...).CombinedOutput()
	if err == nil {
		fmt.Fprintf(ctx.Stdout, "%s", out)
	}
	return err
}
