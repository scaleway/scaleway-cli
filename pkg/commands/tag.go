// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import "fmt"

// TagArgs are flags for the `RunTag` function
type TagArgs struct {
	Snapshot string
	Name     string
}

// RunTag is the handler for 'scw tag'
func RunTag(ctx CommandContext, args TagArgs) error {
	snapshotID := ctx.API.GetSnapshotID(args.Snapshot)
	snapshot, err := ctx.API.GetSnapshot(snapshotID)
	if err != nil {
		return fmt.Errorf("cannot fetch snapshot: %v", err)
	}

	image, err := ctx.API.PostImage(snapshot.Identifier, args.Name)
	if err != nil {
		return fmt.Errorf("cannot create image: %v", err)
	}
	fmt.Fprintln(ctx.Stdout, image)
	return nil
}
