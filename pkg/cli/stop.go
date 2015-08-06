// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"fmt"
	"os"
	"time"

	log "github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"

	"github.com/scaleway/scaleway-cli/pkg/api"
)

var cmdStop = &Command{
	Exec:        runStop,
	UsageLine:   "stop [OPTIONS] SERVER [SERVER...]",
	Description: "Stop a running server",
	Help:        "Stop a running server.",
	Examples: `
    $ scw stop my-running-server my-second-running-server
    $ scw stop -t my-running-server my-second-running-server
    $ scw stop $(scw ps -q)
    $ scw stop $(scw ps | grep mysql | awk '{print $1}')
    $ scw stop server && stop wait server
    $ scw stop -w server
`,
}

func init() {
	cmdStop.Flag.BoolVar(&stopT, []string{"t", "-terminate"}, false, "Stop and trash a server with its volumes")
	cmdStop.Flag.BoolVar(&stopHelp, []string{"h", "-help"}, false, "Print usage")
	cmdStop.Flag.BoolVar(&stopW, []string{"w", "-wait"}, false, "Synchronous stop. Wait for SSH to be ready")
}

// Flags
var stopT bool    // -t flag
var stopHelp bool // -h, --help flag
var stopW bool    // -w, --wait flat

// FIXME: parallelize stop when stopping multiple servers
func runStop(cmd *Command, args []string) {
	if stopHelp {
		cmd.PrintUsage()
	}
	if len(args) < 1 {
		cmd.PrintShortUsage()
	}

	hasError := false
	for _, needle := range args {
		serverID := cmd.API.GetServerID(needle)
		action := "poweroff"
		if stopT {
			action = "terminate"
		}
		err := cmd.API.PostServerAction(serverID, action)
		if err != nil {
			if err.Error() != "server should be running" && err.Error() != "server is being stopped or rebooted" {
				log.Warningf("failed to stop server %s: %s", serverID, err)
				hasError = true
			}
		} else {
			if stopW {
				// We wait for 10 seconds which is the minimal amount of time needed for a server to stop
				time.Sleep(10 * time.Second)
				_, err = api.WaitForServerStopped(cmd.API, serverID)
				if err != nil {
					log.Errorf("failed to wait for server %s: %v", serverID, err)
					hasError = true
				}
			}
			if stopT {
				cmd.API.Cache.RemoveServer(serverID)
			}
			fmt.Println(needle)
		}
	}

	if hasError {
		os.Exit(1)
	}
}
