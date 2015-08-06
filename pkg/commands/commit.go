// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"

	log "github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"

	types "github.com/scaleway/scaleway-cli/pkg/commands/types"
)

var cmdCommit = &types.Command{
	Exec:        runCommit,
	UsageLine:   "commit [OPTIONS] SERVER [NAME]",
	Description: "Create a new snapshot from a server's volume",
	Help:        "Create a new snapshot from a server's volume.",
	Examples: `
    $ scw commit my-stopped-server
    $ scw commit -v 1 my-stopped-server
`,
}

func init() {
	cmdCommit.Flag.IntVar(&commitVolume, []string{"v", "-volume"}, 0, "Volume slot")
	cmdCommit.Flag.BoolVar(&commitHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var commitVolume int // -v, --volume flag
var commitHelp bool  // -h, --help flag

func runCommit(cmd *types.Command, args []string) {
	if commitHelp {
		cmd.PrintUsage()
	}
	if len(args) < 1 {
		cmd.PrintShortUsage()
	}

	serverID := cmd.API.GetServerID(args[0])
	server, err := cmd.API.GetServer(serverID)
	if err != nil {
		log.Fatalf("Cannot fetch server: %v", err)
	}
	var volume = server.Volumes[fmt.Sprintf("%d", commitVolume)]
	var name string
	if len(args) > 1 {
		name = args[1]
	} else {
		name = volume.Name + "-snapshot"
	}
	snapshot, err := cmd.API.PostSnapshot(volume.Identifier, name)
	if err != nil {
		log.Fatalf("Cannot create snapshot: %v", err)
	}
	fmt.Println(snapshot)
}
