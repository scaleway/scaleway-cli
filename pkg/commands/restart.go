// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"os"

	log "github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"

	types "github.com/scaleway/scaleway-cli/pkg/commands/types"
)

var cmdRestart = &types.Command{
	Exec:        runRestart,
	UsageLine:   "restart [OPTIONS] SERVER [SERVER...]",
	Description: "Restart a running server",
	Help:        "Restart a running server.",
}

func init() {
	cmdRestart.Flag.BoolVar(&restartHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var restartHelp bool // -h, --help flag

func runRestart(cmd *types.Command, args []string) {
	if restartHelp {
		cmd.PrintUsage()
	}
	if len(args) < 1 {
		cmd.PrintShortUsage()
	}

	hasError := false
	for _, needle := range args {
		server := cmd.API.GetServerID(needle)
		err := cmd.API.PostServerAction(server, "reboot")
		if err != nil {
			if err.Error() != "server is being stopped or rebooted" {
				log.Errorf("failed to restart server %s: %s", server, err)
				hasError = true
			}
		} else {
			fmt.Println(needle)
		}
		if hasError {
			os.Exit(1)
		}
	}
}
