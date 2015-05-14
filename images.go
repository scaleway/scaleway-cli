package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/docker/docker/pkg/units"
)

var cmdImages = &Command{
	Exec:        runImages,
	UsageLine:   "images [OPTIONS]",
	Description: "List images",
	Help:        "List images.",
}

func init() {
	// FIXME: -h
	cmdImages.Flag.BoolVar(&imagesA, []string{"a", "-all"}, false, "Show all iamges")
	cmdImages.Flag.BoolVar(&imagesNoTrunc, []string{"-no-trunc"}, false, "Don't truncate output")
	cmdImages.Flag.BoolVar(&imagesQ, []string{"q", "-quiet"}, false, "Only show numeric IDs")
}

// Flags
var imagesA bool       // -a flag
var imagesQ bool       // -q flag
var imagesNoTrunc bool // -no-trunc flag

func runImages(cmd *Command, args []string) {
	images, err := cmd.API.GetImages()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to fetch images from the Scaleway API: %v\n", err)
		os.Exit(1)
	}

	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	defer w.Flush()
	if !imagesQ {
		fmt.Fprintf(w, "REPOSITORY\tTAG\tIMAGE ID\tCREATED\tVIRTUAL SIZE\n")
	}
	for _, image := range *images {
		if imagesQ {
			fmt.Fprintf(w, "%s\n", image.Identifier)
		} else {
			tag := "latest"
			virtualSize := units.HumanSize(float64(image.RootVolume.Size))
			short_id := truncIf(image.Identifier, 8, !imagesNoTrunc)
			short_name := truncIf(wordify(image.Name), 25, !imagesNoTrunc)
			creationTime, _ := time.Parse("2006-01-02T15:04:05.000000+00:00", image.CreationDate)
			creationDate := units.HumanDuration(time.Now().UTC().Sub(creationTime))
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", short_name, tag, short_id, creationDate, virtualSize)
		}
	}
}
