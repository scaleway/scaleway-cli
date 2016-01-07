// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
)

// CommitArgs are flags for the `RunCommit` function
type CommitArgs struct {
	Volume int
	Server string
	Name   string
}

// RunCommit is the handler for 'scw commit'
func RunCommit(ctx CommandContext, args CommitArgs) error {
	serverID, err := ctx.API.GetServerID(args.Server)
	if err != nil {
		return err
	}
	server, err := ctx.API.GetServer(serverID)
	if err != nil {
		return fmt.Errorf("Cannot fetch server: %v", err)
	}
	var volume = server.Volumes[fmt.Sprintf("%d", args.Volume)]
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
	fmt.Fprintln(ctx.Stdout, snapshot)
	return nil
}
