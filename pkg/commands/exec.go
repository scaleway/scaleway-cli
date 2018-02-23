// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/utils"
)

// ExecArgs are flags for the `RunExec` function
type ExecArgs struct {
	Timeout                float64
	Wait                   bool
	Gateway                string
	Server                 string
	Command                []string
	SSHUser                string
	SSHPort                int
	EnableSSHKeyForwarding bool
}

// RunExec is the handler for 'scw exec'
func RunExec(ctx CommandContext, args ExecArgs) error {
	var fingerprints []string

	done := make(chan struct{})

	serverID, err := ctx.API.GetServerID(args.Server)
	if err != nil {
		return err
	}

	go func() {
		fingerprints = ctx.API.GetSSHFingerprintFromServer(serverID)
		close(done)
	}()
	// Resolve gateway
	if args.Gateway == "" {
		args.Gateway = ctx.Getenv("SCW_GATEWAY")
	}
	var gateway string

	if args.Gateway == serverID || args.Gateway == args.Server {
		logrus.Debugf("The server and the gateway are the same host, using direct access to the server")
		gateway = ""
	} else {
		gateway, err = api.ResolveGateway(ctx.API, args.Gateway)
		if err != nil {
			return fmt.Errorf("Cannot resolve Gateway '%s': %v", args.Gateway, err)
		}
		if gateway != "" {
			logrus.Debugf("The server will be accessed using the gateway '%s' as a SSH relay", gateway)
		}
	}

	var server *api.ScalewayServer
	if args.Wait {
		// --wait
		logrus.Debugf("Waiting for server to be ready")
		server, err = api.WaitForServerReady(ctx.API, serverID, gateway)
		if err != nil {
			return fmt.Errorf("Failed to wait for server to be ready, %v", err)
		}
	} else {
		// no --wait
		logrus.Debugf("scw won't wait for the server to be ready, if it is not, the command will fail")
		server, err = ctx.API.GetServer(serverID)
		if err != nil {
			rerr := fmt.Errorf("Failed to get server information for %s: %v", serverID, err)
			if err.Error() == `"`+serverID+`" not found` {
				return fmt.Errorf("%v\nmaybe try to flush the cache with : scw _flush-cache", rerr)
			}
			return rerr
		}
	}

	if server.PublicAddress.IP == "" && gateway == "" {
		logrus.Warn(`Your host has no public IP address, you should use '--gateway', see 'scw help exec'`)
	}

	// --timeout
	if args.Timeout > 0 {
		logrus.Debugf("Setting up a global timeout of %d seconds", args.Timeout)
		// FIXME: avoid use of log.Fatalf here
		go func() {
			time.Sleep(time.Duration(args.Timeout*1000) * time.Millisecond)
			logrus.Fatalf("Operation timed out")
		}()
	}

	<-done
	if len(fingerprints) > 0 {
		for i := range fingerprints {
			fmt.Fprintf(ctx.Stdout, "%s\n", fingerprints[i])
		}
	}
	logrus.Debugf("PublicDNS %s", serverID+api.URLPublicDNS)
	logrus.Debugf("PrivateDNS %s", serverID+api.URLPrivateDNS)
	if err = utils.SSHExec(server.PublicAddress.IP, server.PrivateIP, args.SSHUser, args.SSHPort, args.Command, !args.Wait, gateway, args.EnableSSHKeyForwarding); err != nil {
		return fmt.Errorf("Failed to run the command: %v", err)
	}

	logrus.Debugf("Command successfully executed")
	return nil
}
