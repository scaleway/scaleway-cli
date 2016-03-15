// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/renstrom/fuzzysearch/fuzzy"
	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/utils"
)

// SearchArgs are flags for the `RunSearch` function
type SearchArgs struct {
	Term    string
	NoTrunc bool
}

// RunSearch is the handler for 'scw search'
func RunSearch(ctx CommandContext, args SearchArgs) error {
	// FIXME: parallelize API calls

	term := strings.ToLower(args.Term)
	w := tabwriter.NewWriter(ctx.Stdout, 10, 1, 3, ' ', 0)
	defer w.Flush()
	fmt.Fprintf(w, "NAME\tDESCRIPTION\tSTARS\tOFFICIAL\tAUTOMATED\n")

	var entries = []api.ScalewayImageInterface{}

	images, err := ctx.API.GetImages()
	if err != nil {
		return fmt.Errorf("unable to fetch images from the Scaleway API: %v", err)
	}
	for _, val := range *images {
		if fuzzy.Match(term, strings.ToLower(val.Name)) {
			entries = append(entries, api.ScalewayImageInterface{
				Type:   "image",
				Name:   val.Name,
				Public: val.Public,
			})
		}
	}

	snapshots, err := ctx.API.GetSnapshots()
	if err != nil {
		return fmt.Errorf("unable to fetch snapshots from the Scaleway API: %v", err)
	}
	for _, val := range *snapshots {
		if fuzzy.Match(term, strings.ToLower(val.Name)) {
			entries = append(entries, api.ScalewayImageInterface{
				Type:   "snapshot",
				Name:   val.Name,
				Public: false,
			})
		}
	}

	for _, image := range entries {
		// name field
		name := utils.TruncIf(utils.Wordify(image.Name), 45, !args.NoTrunc)

		// description field
		var description string
		switch image.Type {
		case "image":
			if image.Public {
				description = "public image"
			} else {
				description = "user image"
			}

		case "snapshot":
			description = "user snapshot"
		}
		description = utils.TruncIf(utils.Wordify(description), 45, !args.NoTrunc)

		// official field
		var official string
		if image.Public {
			official = "[OK]"
		} else {
			official = ""
		}

		fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\n", name, description, 0, official, "")
	}
	return nil
}
