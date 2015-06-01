package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
)

var cmdSearch = &Command{
	Exec:        runSearch,
	UsageLine:   "search [OPTIONS] TERM",
	Description: "Search the Scaleway Hub for images",
	Help:        "Search the Scaleway Hub for images.",
}

func init() {
	// FIXME: -h
	cmdSearch.Flag.BoolVar(&searchNoTrunc, []string{"-no-trunc"}, false, "Don't truncate output")
}

// Flags
var searchNoTrunc bool // --no-trunc flag

func runSearch(cmd *Command, args []string) {
	if len(args) != 1 {
		log.Fatalf("usage: scw %s", cmd.UsageLine)
	}

	w := tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', 0)
	defer w.Flush()
	fmt.Fprintf(w, "NAME\tDESCRIPTION\tSTARS\tOFFICIAL\tAUTOMATED\n")

	var entries = []ScalewayImageInterface{}

	images, err := cmd.API.GetImages()
	if err != nil {
		log.Fatalf("unable to fetch images from the Scaleway API: %v", err)
	}
	for _, val := range *images {
		entries = append(entries, ScalewayImageInterface{
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
		entries = append(entries, ScalewayImageInterface{
			Type:   "snapshot",
			Name:   val.Name,
			Public: false,
		})
	}

	for _, image := range entries {
		// name field
		name := truncIf(wordify(image.Name), 45, !searchNoTrunc)

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
		description = truncIf(wordify(description), 45, !searchNoTrunc)

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
