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
	serverID, err := ctx.API.GetServerID(args.Server)
	if err != nil {
		return err
	}
	if err = ctx.API.PatchServer(serverID,
		api.ScalewayServerPatchDefinition{
			Name: &args.NewName,
		}); err != nil {
		return fmt.Errorf("cannot rename server: %v", err)
	}
	if server, err := ctx.API.GetServer(serverID); err == nil {
		ctx.API.Cache.InsertServer(serverID, server.Location.ZoneID, server.Arch, server.Organization, server.Name)
	}
	return nil
}
