package main

import (
	"fmt"
	"os"
	"sort"
	"sync"
	"text/tabwriter"
	"time"

	log "github.com/Sirupsen/logrus"
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
	Public       bool
	Type         string
}

// ByCreationDate sorts images by CreationDate field
type ByCreationDate []ScalewayImageInterface

func (a ByCreationDate) Len() int           { return len(a) }
func (a ByCreationDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCreationDate) Less(i, j int) bool { return a[j].CreationDate.Before(a[i].CreationDate) }

func init() {
	cmdImages.Flag.BoolVar(&imagesA, []string{"a", "-all"}, false, "Show all iamges")
	cmdImages.Flag.BoolVar(&imagesNoTrunc, []string{"-no-trunc"}, false, "Don't truncate output")
	cmdImages.Flag.BoolVar(&imagesQ, []string{"q", "-quiet"}, false, "Only show numeric IDs")
	cmdImages.Flag.BoolVar(&imagesHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var imagesA bool       // -a flag
var imagesQ bool       // -q flag
var imagesNoTrunc bool // -no-trunc flag
var imagesHelp bool    // -h, --help flag

func runImages(cmd *Command, args []string) {
	if imagesHelp {
		cmd.PrintUsage()
	}
	if len(args) != 0 {
		cmd.PrintShortUsage()
	}

	wg := sync.WaitGroup{}
	chEntries := make(chan ScalewayImageInterface)
	var entries = []ScalewayImageInterface{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		images, err := cmd.API.GetImages()
		if err != nil {
			log.Fatalf("unable to fetch images from the Scaleway API: %v", err)
		}
		for _, val := range *images {
			creationDate, err := time.Parse("2006-01-02T15:04:05.000000+00:00", val.CreationDate)
			if err != nil {
				log.Fatalf("unable to parse creation date from the Scaleway API: %v", err)
			}
			chEntries <- ScalewayImageInterface{
				Type:         "image",
				CreationDate: creationDate,
				Identifier:   val.Identifier,
				Name:         val.Name,
				Public:       val.Public,
				Tag:          "latest",
				VirtualSize:  float64(val.RootVolume.Size),
			}
		}
	}()

	if imagesA {
		wg.Add(1)
		go func() {
			defer wg.Done()
			snapshots, err := cmd.API.GetSnapshots()
			if err != nil {
				log.Fatalf("unable to fetch snapshots from the Scaleway API: %v", err)
			}
			for _, val := range *snapshots {
				creationDate, err := time.Parse("2006-01-02T15:04:05.000000+00:00", val.CreationDate)
				if err != nil {
					log.Fatalf("unable to parse creation date from the Scaleway API: %v", err)
				}
				chEntries <- ScalewayImageInterface{
					Type:         "snapshot",
					CreationDate: creationDate,
					Identifier:   val.Identifier,
					Name:         val.Name,
					Tag:          "<snapshot>",
					VirtualSize:  float64(val.Size),
					Public:       false,
				}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			bootscripts, err := cmd.API.GetBootscripts()
			if err != nil {
				log.Fatalf("unable to fetch bootscripts from the Scaleway API: %v", err)
			}
			for _, val := range *bootscripts {
				chEntries <- ScalewayImageInterface{
					Type:       "bootscript",
					Identifier: val.Identifier,
					Name:       val.Title,
					Tag:        "<bootscript>",
					Public:     false,
				}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			volumes, err := cmd.API.GetVolumes()
			if err != nil {
				log.Fatalf("unable to fetch volumes from the Scaleway API: %v", err)
			}
			for _, val := range *volumes {
				creationDate, err := time.Parse("2006-01-02T15:04:05.000000+00:00", val.CreationDate)
				if err != nil {
					log.Fatalf("unable to parse creation date from the Scaleway API: %v", err)
				}
				chEntries <- ScalewayImageInterface{
					Type:         "volume",
					CreationDate: creationDate,
					Identifier:   val.Identifier,
					Name:         val.Name,
					Tag:          "<volume>",
					VirtualSize:  float64(val.Size),
					Public:       false,
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(chEntries)
	}()

	done := false
	for {
		select {
		case entry, ok := <-chEntries:
			if !ok {
				done = true
				break
			}
			entries = append(entries, entry)
		}
		if done {
			break
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
			tag := image.Tag
			shortID := truncIf(image.Identifier, 8, !imagesNoTrunc)
			name := wordify(image.Name)
			if !image.Public {
				name = "user/" + name
			}
			shortName := truncIf(name, 25, !imagesNoTrunc)
			var creationDate, virtualSize string
			if image.CreationDate.IsZero() {
				creationDate = "n/a"
			} else {
				creationDate = units.HumanDuration(time.Now().UTC().Sub(image.CreationDate))
			}
			if image.VirtualSize == 0 {
				virtualSize = "n/a"
			} else {
				virtualSize = units.HumanSize(image.VirtualSize)
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", shortName, tag, shortID, creationDate, virtualSize)
		}
	}
}
