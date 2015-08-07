// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/utils"
	"github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"
)

// TopArgs are flags for the `RunTop` function
type TopArgs struct {
	Server  string
	Gateway string
}

// RunTop is the handler for 'scw top'
func RunTop(ctx CommandContext, args TopArgs) error {
	serverID := ctx.API.GetServerID(args.Server)
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

	execCmd := utils.NewSSHExecCmd(server.PublicAddress.IP, server.PrivateIP, true, nil, []string{command}, gateway)
	logrus.Debugf("Executing: ssh %s", strings.Join(execCmd, " "))
	out, err := exec.Command("ssh", execCmd...).CombinedOutput()
	fmt.Printf("%s", out)
	return err
}
