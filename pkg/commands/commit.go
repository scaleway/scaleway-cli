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
	Exec:        cmdExecCommit,
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

// CommitArgs are flags for the `RunCommit` function
type CommitArgs struct {
	Volume int
	Server string
	Name   string
}

func cmdExecCommit(cmd *types.Command, rawArgs []string) {
	if commitHelp {
		cmd.PrintUsage()
	}
	if len(rawArgs) < 1 {
		cmd.PrintShortUsage()
	}

	args := CommitArgs{
		Volume: commitVolume,
		Server: rawArgs[0],
		Name:   "",
	}
	if len(rawArgs) > 1 {
		args.Name = rawArgs[1]
	}

	ctx := cmd.GetContext(rawArgs)
	err := RunCommit(ctx, args)
	if err != nil {
		log.Fatalf("Cannot execute 'commit': %v", err)
	}
}

// RunCommit is the handler for 'scw commit'
func RunCommit(ctx types.CommandContext, args CommitArgs) error {
	serverID := ctx.API.GetServerID(args.Server)
	server, err := ctx.API.GetServer(serverID)
	if err != nil {
		return fmt.Errorf("Cannot fetch server: %v", err)
	}
	var volume = server.Volumes[fmt.Sprintf("%d", commitVolume)]
	var name string
	if args.Name != "" {
		name = args.Name
	} else {
		name = volume.Name + "-snapshot"
	}
	snapshot, err := ctx.API.PostSnapshot(volume.Identifier, name)
	if err != nil {
		return fmt.Errorf("Cannot create snapshot: %v", err)
	}
	fmt.Println(snapshot)
	return nil
}
