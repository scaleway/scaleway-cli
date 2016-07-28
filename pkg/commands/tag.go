// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"

	"github.com/moul/anonuuid"
)

// TagArgs are flags for the `RunTag` function
type TagArgs struct {
	Snapshot   string
	Bootscript string
	Name       string
	Arch       string
}

// RunTag is the handler for 'scw tag'
func RunTag(ctx CommandContext, args TagArgs) error {
	snapshotID, err := ctx.API.GetSnapshotID(args.Snapshot)
	if err != nil {
		return err
	}
	snapshot, err := ctx.API.GetSnapshot(snapshotID)
	if err != nil {
		return fmt.Errorf("cannot fetch snapshot: %v", err)
	}

	bootscriptID := ""
	if args.Bootscript != "" {
		if anonuuid.IsUUID(args.Bootscript) == nil {
			bootscriptID = args.Bootscript
		} else {
			bootscriptID, err = ctx.API.GetBootscriptID(args.Bootscript, args.Arch)
			if err != nil {
				return err
			}
		}
	}
	image, err := ctx.API.PostImage(snapshot.Identifier, args.Name, bootscriptID, args.Arch)
	if err != nil {
		return fmt.Errorf("cannot create image: %v", err)
	}
	fmt.Fprintln(ctx.Stdout, image)
	return nil
}
