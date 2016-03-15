// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"

	"github.com/Sirupsen/logrus"
)

// RmiArgs are flags for the `RunRmi` function
type RmiArgs struct {
	Identifier []string // images/volumes/snapshots
}

// RunRmi is the handler for 'scw rmi'
func RunRmi(ctx CommandContext, args RmiArgs) error {
	hasError := false
	for _, needle := range args.Identifier {
		if image, err := ctx.API.GetImageID(needle, "*"); err == nil {
			if err = ctx.API.DeleteImage(image.Identifier); err != nil {
				logrus.Errorf("failed to delete image %s: %s", image.Identifier, err)
				hasError = true
			} else {
				fmt.Fprintln(ctx.Stdout, needle)
			}
			continue
		}
		if snapshotID, err := ctx.API.GetSnapshotID(needle); err == nil {
			if err = ctx.API.DeleteSnapshot(snapshotID); err != nil {
				logrus.Errorf("failed to delete snapshot %s: %s", snapshotID, err)
				hasError = true
			} else {
				fmt.Fprintln(ctx.Stdout, needle)
			}
			continue
		}
		if volumeID, err := ctx.API.GetVolumeID(needle); err == nil {
			if err = ctx.API.DeleteVolume(volumeID); err != nil {
				logrus.Errorf("failed to delete volume %s: %s", volumeID, err)
				hasError = true
			} else {
				fmt.Fprintln(ctx.Stdout, needle)
			}
			continue
		}
		hasError = true
	}
	if hasError {
		return fmt.Errorf("at least 1 image/snapshot/volume failed to be removed")
	}
	return nil
}
