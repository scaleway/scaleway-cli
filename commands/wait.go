// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"

	types "github.com/scaleway/scaleway-cli/commands/types"
)

var cmdWait = &types.Command{
	Exec:        runWait,
	UsageLine:   "wait [OPTIONS] SERVER [SERVER...]",
	Description: "Block until a server stops",
	Help:        "Block until a server stops.",
}

func init() {
	cmdWait.Flag.BoolVar(&waitHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var waitHelp bool // -h, --help flag

func runWait(cmd *types.Command, args []string) {
	if waitHelp {
		cmd.PrintUsage()
	}
	if len(args) < 1 {
		cmd.PrintShortUsage()
	}

	hasError := false
	for _, needle := range args {
		serverIdentifier := cmd.API.GetServerID(needle)
		for {
			server, err := cmd.API.GetServer(serverIdentifier)
			if err != nil {
				log.Errorf("failed to retrieve information from server %s: %s", serverIdentifier, err)
				hasError = true
				break
			}
			if server.State == "stopped" {
				break
			}
			time.Sleep(1 * time.Second)
		}
	}
	if hasError {
		os.Exit(1)
	}
}
