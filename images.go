package main

import (
	"fmt"
	"os"
	"sort"
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

// ScalewayImageInterface is an interface to multiple Scaleway items
type ScalewayImageInterface struct {
	CreationDate time.Time
	Identifier   string
	Name         string
	Tag          string
	VirtualSize  float64
}

type ByCreationDate []ScalewayImageInterface

func (a ByCreationDate) Len() int           { return len(a) }
func (a ByCreationDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCreationDate) Less(i, j int) bool { return a[j].CreationDate.Before(a[i].CreationDate) }

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
	var entries = []ScalewayImageInterface{}

	images, err := cmd.API.GetImages()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to fetch images from the Scaleway API: %v\n", err)
		os.Exit(1)
	}
	for _, val := range *images {
		creationDate, err := time.Parse("2006-01-02T15:04:05.000000+00:00", val.CreationDate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to parse creation date from the Scaleway API: %v\n", err)
			os.Exit(1)
		}
		entries = append(entries, ScalewayImageInterface{
			CreationDate: creationDate,
			Identifier:   val.Identifier,
			Name:         val.Name,
			Tag:          "latest",
			VirtualSize:  float64(val.RootVolume.Size),
		})
	}

	if imagesA {
		snapshots, err := cmd.API.GetSnapshots()
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to fetch snapshots from the Scaleway API: %v\n", err)
			os.Exit(1)
		}
		for _, val := range *snapshots {
			creationDate, err := time.Parse("2006-01-02T15:04:05.000000+00:00", val.CreationDate)
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to parse creation date from the Scaleway API: %v\n", err)
				os.Exit(1)
			}
			entries = append(entries, ScalewayImageInterface{
				CreationDate: creationDate,
				Identifier:   val.Identifier,
				Name:         val.Name,
				Tag:          "<none>",
				VirtualSize:  float64(val.Size),
			})
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	defer w.Flush()
	if !imagesQ {
		fmt.Fprintf(w, "REPOSITORY\tTAG\tIMAGE ID\tCREATED\tVIRTUAL SIZE\n")
	}
	sort.Sort(ByCreationDate(entries))
	for _, image := range entries {
		if imagesQ {
			fmt.Fprintf(w, "%s\n", image.Identifier)
		} else {
			tag := "latest"
			virtualSize := units.HumanSize(image.VirtualSize)
			short_id := truncIf(image.Identifier, 8, !imagesNoTrunc)
			short_name := truncIf(wordify(image.Name), 25, !imagesNoTrunc)
			creationDate := units.HumanDuration(time.Now().UTC().Sub(image.CreationDate))
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", short_name, tag, short_id, creationDate, virtualSize)
		}
	}
}
