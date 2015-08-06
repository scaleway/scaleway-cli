// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	log "github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"

	"github.com/scaleway/scaleway-cli/pkg/api"
)

var cmdRename = &Command{
	Exec:        runRename,
	UsageLine:   "rename [OPTIONS] SERVER NEW_NAME",
	Description: "Rename a server",
	Help:        "Rename a server.",
}

func init() {
	cmdRename.Flag.BoolVar(&renameHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var renameHelp bool // -h, --help flag

func runRename(cmd *Command, args []string) {
	if renameHelp {
		cmd.PrintUsage()
	}
	if len(args) != 2 {
		cmd.PrintShortUsage()
	}

	serverID := cmd.API.GetServerID(args[0])

	var server api.ScalewayServerPatchDefinition
	server.Name = &args[1]

	err := cmd.API.PatchServer(serverID, server)
	if err != nil {
		log.Fatalf("Cannot rename server: %v", err)
	} else {
		cmd.API.Cache.InsertServer(serverID, *server.Name)
	}
}
