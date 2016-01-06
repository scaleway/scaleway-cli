// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/pkg/api"
)

// RenameArgs are flags for the `RunRename` function
type RenameArgs struct {
	Server  string
	NewName string
}

// RunRename is the handler for 'scw rename'
func RunRename(ctx CommandContext, args RenameArgs) error {
	serverID := ctx.API.GetServerID(args.Server)

	var server api.ScalewayServerPatchDefinition
	server.Name = &args.NewName

	err := ctx.API.PatchServer(serverID, server)
	if err != nil {
		return fmt.Errorf("cannot rename server: %v", err)
	}
	// FIXME region, arch, owner, title
	ctx.API.Cache.InsertServer(serverID, "fr-1", *server.Arch, *server.Organization, *server.Name)
	return nil
}
