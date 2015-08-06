// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	log "github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"

	"github.com/scaleway/scaleway-cli/pkg/api"
	types "github.com/scaleway/scaleway-cli/pkg/commands/types"
	"github.com/scaleway/scaleway-cli/pkg/utils"
)

var cmdSearch = &types.Command{
	Exec:        runSearch,
	UsageLine:   "search [OPTIONS] TERM",
	Description: "Search the Scaleway Hub for images",
	Help:        "Search the Scaleway Hub for images.",
}

func init() {
	cmdSearch.Flag.BoolVar(&searchNoTrunc, []string{"-no-trunc"}, false, "Don't truncate output")
	cmdSearch.Flag.BoolVar(&searchHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var searchNoTrunc bool // --no-trunc flag
var searchHelp bool    // -h, --help flag

func runSearch(cmd *types.Command, args []string) {
	if searchHelp {
		cmd.PrintUsage()
	}
	if len(args) != 1 {
		cmd.PrintShortUsage()
	}

	w := tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', 0)
	defer w.Flush()
	fmt.Fprintf(w, "NAME\tDESCRIPTION\tSTARS\tOFFICIAL\tAUTOMATED\n")

	var entries = []api.ScalewayImageInterface{}

	images, err := cmd.API.GetImages()
	if err != nil {
		log.Fatalf("unable to fetch images from the Scaleway API: %v", err)
	}
	for _, val := range *images {
		entries = append(entries, api.ScalewayImageInterface{
			Type:   "image",
			Name:   val.Name,
			Public: val.Public,
		})
	}

	snapshots, err := cmd.API.GetSnapshots()
	if err != nil {
		log.Fatalf("unable to fetch snapshots from the Scaleway API: %v", err)
	}
	for _, val := range *snapshots {
		entries = append(entries, api.ScalewayImageInterface{
			Type:   "snapshot",
			Name:   val.Name,
			Public: false,
		})
	}

	for _, image := range entries {
		// name field
		name := utils.TruncIf(utils.Wordify(image.Name), 45, !searchNoTrunc)

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
		description = utils.TruncIf(utils.Wordify(description), 45, !searchNoTrunc)

		// official field
		var official string
		if image.Public {
			official = "[OK]"
		} else {
			official = ""
		}

		fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\n", name, description, 0, official, "")
	}
}
