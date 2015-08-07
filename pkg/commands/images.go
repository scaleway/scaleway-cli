// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"sort"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/utils"
	"github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"
	"github.com/scaleway/scaleway-cli/vendor/github.com/docker/docker/pkg/units"
)

// ImagesArgs are flags for the `RunImages` function
type ImagesArgs struct {
	All     bool
	NoTrunc bool
	Quiet   bool
}

// RunImages is the handler for 'scw images'
func RunImages(ctx CommandContext, args ImagesArgs) error {
	wg := sync.WaitGroup{}
	chEntries := make(chan api.ScalewayImageInterface)
	var entries = []api.ScalewayImageInterface{}

	// FIXME: remove log.Fatalf in routines

	wg.Add(1)
	go func() {
		defer wg.Done()
		images, err := ctx.API.GetImages()
		if err != nil {
			logrus.Fatalf("unable to fetch images from the Scaleway API: %v", err)
		}
		for _, val := range *images {
			creationDate, err := time.Parse("2006-01-02T15:04:05.000000+00:00", val.CreationDate)
			if err != nil {
				logrus.Fatalf("unable to parse creation date from the Scaleway API: %v", err)
			}
			chEntries <- api.ScalewayImageInterface{
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

	if args.All {
		wg.Add(1)
		go func() {
			defer wg.Done()
			snapshots, err := ctx.API.GetSnapshots()
			if err != nil {
				logrus.Fatalf("unable to fetch snapshots from the Scaleway API: %v", err)
			}
			for _, val := range *snapshots {
				creationDate, err := time.Parse("2006-01-02T15:04:05.000000+00:00", val.CreationDate)
				if err != nil {
					logrus.Fatalf("unable to parse creation date from the Scaleway API: %v", err)
				}
				chEntries <- api.ScalewayImageInterface{
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
			bootscripts, err := ctx.API.GetBootscripts()
			if err != nil {
				logrus.Fatalf("unable to fetch bootscripts from the Scaleway API: %v", err)
			}
			for _, val := range *bootscripts {
				chEntries <- api.ScalewayImageInterface{
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
			volumes, err := ctx.API.GetVolumes()
			if err != nil {
				logrus.Fatalf("unable to fetch volumes from the Scaleway API: %v", err)
			}
			for _, val := range *volumes {
				creationDate, err := time.Parse("2006-01-02T15:04:05.000000+00:00", val.CreationDate)
				if err != nil {
					logrus.Fatalf("unable to parse creation date from the Scaleway API: %v", err)
				}
				chEntries <- api.ScalewayImageInterface{
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

	w := tabwriter.NewWriter(ctx.Stdout, 20, 1, 3, ' ', 0)
	defer w.Flush()
	if !args.Quiet {
		fmt.Fprintf(w, "REPOSITORY\tTAG\tIMAGE ID\tCREATED\tVIRTUAL SIZE\n")
	}
	sort.Sort(api.ByCreationDate(entries))
	for _, image := range entries {
		if args.Quiet {
			fmt.Fprintf(ctx.Stdout, "%s\n", image.Identifier)
		} else {
			tag := image.Tag
			shortID := utils.TruncIf(image.Identifier, 8, !args.NoTrunc)
			name := utils.Wordify(image.Name)
			if !image.Public && image.Type == "image" {
				name = "user/" + name
			}
			shortName := utils.TruncIf(name, 25, !args.NoTrunc)
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
	return nil
}
