// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/scaleway/scaleway-cli/pkg/api"
)

// StopArgs are flags for the `RunStop` function
type StopArgs struct {
	Terminate bool
	Wait      bool
	Servers   []string
}

// RunStop is the handler for 'scw stop'
func RunStop(ctx CommandContext, args StopArgs) error {
	// FIXME: parallelize stop when stopping multiple servers
	hasError := false
	for _, needle := range args.Servers {
		serverID, err := ctx.API.GetServerID(needle)
		if err != nil {
			return err
		}
		action := "poweroff"
		if args.Terminate {
			action = "terminate"
		}
		if err = ctx.API.PostServerAction(serverID, action); err != nil {
			if err.Error() != "server should be running" && err.Error() != "server is being stopped or rebooted" {
				logrus.Warningf("failed to stop server %s: %s", serverID, err)
				hasError = true
			}
		} else {
			if args.Wait {
				// We wait for 10 seconds which is the minimal amount of time needed for a server to stop
				time.Sleep(10 * time.Second)
				if _, err = api.WaitForServerStopped(ctx.API, serverID); err != nil {
					logrus.Errorf("failed to wait for server %s: %v", serverID, err)
					hasError = true
				}
			}
			if args.Terminate {
				ctx.API.Cache.RemoveServer(serverID)
			}
			fmt.Fprintln(ctx.Stdout, needle)
		}
	}

	if hasError {
		return fmt.Errorf("at least 1 server failed to be stopped")
	}
	return nil
}
