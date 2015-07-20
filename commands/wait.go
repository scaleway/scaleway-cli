// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"os"

	log "github.com/Sirupsen/logrus"

	"github.com/scaleway/scaleway-cli/api"
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

		_, err := api.WaitForServerStopped(cmd.API, serverIdentifier)
		if err != nil {
			log.Errorf("failed to wait for server %s: %v", serverIdentifier, err)
			hasError = true
		}
	}

	if hasError {
		os.Exit(1)
	}
}
