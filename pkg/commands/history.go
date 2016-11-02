// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"text/tabwriter"
	"time"

	"github.com/docker/go-units"
	"github.com/scaleway/scaleway-cli/pkg/utils"
)

// HistoryArgs are flags for the `RunHistory` function
type HistoryArgs struct {
	NoTrunc bool
	Quiet   bool
	Image   string
	Arch    string
}

// RunHistory is the handler for 'scw history'
func RunHistory(ctx CommandContext, args HistoryArgs) error {
	imageID, err := ctx.API.GetImageID(args.Image, args.Arch)
	if err != nil {
		return err
	}
	image, err := ctx.API.GetImage(imageID.Identifier)
	if err != nil {
		return fmt.Errorf("cannot get image %s: %v", imageID.Identifier, err)
	}

	if args.Quiet {
		fmt.Fprintln(ctx.Stdout, imageID.Identifier)
		return nil
	}

	w := tabwriter.NewWriter(ctx.Stdout, 10, 1, 3, ' ', 0)
	defer w.Flush()
	fmt.Fprintf(w, "IMAGE\tCREATED\tCREATED BY\tSIZE\n")

	identifier := utils.TruncIf(image.Identifier, 8, !args.NoTrunc)

	creationDate, err := time.Parse("2006-01-02T15:04:05.000000+00:00", image.CreationDate)
	if err != nil {
		return fmt.Errorf("unable to parse creation date from the Scaleway API: %v", err)
	}
	creationDateStr := units.HumanDuration(time.Now().UTC().Sub(creationDate))

	volumeName := utils.TruncIf(image.RootVolume.Name, 25, !args.NoTrunc)
	size := units.HumanSize(float64(image.RootVolume.Size))

	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", identifier, creationDateStr, volumeName, size)
	return nil
}
