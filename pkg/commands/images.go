// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/go-units"
	"github.com/renstrom/fuzzysearch/fuzzy"
	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/utils"
)

// ImagesArgs are flags for the `RunImages` function
type ImagesArgs struct {
	All     bool
	NoTrunc bool
	Quiet   bool
	Filters map[string]string
}

// RunImages is the handler for 'scw images'
func RunImages(ctx CommandContext, args ImagesArgs) error {
	wg := sync.WaitGroup{}
	chEntries := make(chan api.ScalewayImageInterface)
	errChan := make(chan error, 10)
	var entries = []api.ScalewayImageInterface{}

	filterType := args.Filters["type"]

	if filterType == "" || filterType == "image" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			images, err := ctx.API.GetImages()
			if err != nil {
				errChan <- fmt.Errorf("unable to fetch images from the Scaleway API: %v", err)
				return
			}
			for _, val := range *images {
				creationDate, err := time.Parse("2006-01-02T15:04:05.000000+00:00", val.CreationDate)
				if err != nil {
					errChan <- fmt.Errorf("unable to parse creation date from the Scaleway API: %v", err)
					return
				}
				archAvailable := make(map[string]struct{})
				zoneAvailable := make(map[string]struct{})

				for _, version := range val.Versions {
					if val.CurrentPublicVersion == version.ID {
						for _, local := range version.LocalImages {
							archAvailable[local.Arch] = struct{}{}
							zoneAvailable[local.Zone] = struct{}{}
						}
						break
					}
				}
				regions := []string{}
				for k := range zoneAvailable {
					regions = append(regions, k)
				}
				archs := []string{}
				for k := range archAvailable {
					archs = append(archs, k)
				}
				chEntries <- api.ScalewayImageInterface{
					Type:         "image",
					CreationDate: creationDate,
					Identifier:   val.CurrentPublicVersion,
					Name:         val.Name,
					Tag:          "latest",
					Organization: val.Organization.ID,
					Public:       val.Public,
					Region:       regions,
					Archs:        archs,
				}
			}
		}()
	}

	if args.All || filterType != "" {
		if filterType == "" || filterType == "snapshot" {
			wg.Add(1)
			go func() {
				defer wg.Done()
				snapshots, err := ctx.API.GetSnapshots()
				if err != nil {
					errChan <- fmt.Errorf("unable to fetch snapshots from the Scaleway API: %v", err)
					return
				}
				for _, val := range *snapshots {
					creationDate, err := time.Parse("2006-01-02T15:04:05.000000+00:00", val.CreationDate)
					if err != nil {
						errChan <- fmt.Errorf("unable to parse creation date from the Scaleway API: %v", err)
						return
					}
					chEntries <- api.ScalewayImageInterface{
						Type:         "snapshot",
						CreationDate: creationDate,
						Identifier:   val.Identifier,
						Name:         val.Name,
						Tag:          "<snapshot>",
						VirtualSize:  val.Size,
						Public:       false,
						Organization: val.Organization,
						// FIXME the region should not be hardcoded
						Region: []string{"par1"},
					}
				}
			}()
		}

		if filterType == "" || filterType == "bootscript" {
			wg.Add(1)
			go func() {
				defer wg.Done()
				bootscripts, err := ctx.API.GetBootscripts()
				if err != nil {
					errChan <- fmt.Errorf("unable to fetch bootscripts from the Scaleway API: %v", err)
					return
				}
				for _, val := range *bootscripts {
					chEntries <- api.ScalewayImageInterface{
						Type:       "bootscript",
						Identifier: val.Identifier,
						Name:       val.Title,
						Tag:        "<bootscript>",
						Public:     false,
						Region:     []string{""},
						Archs:      []string{val.Arch},
					}
				}
			}()
		}

		if filterType == "" || filterType == "volume" {
			wg.Add(1)
			go func() {
				defer wg.Done()
				volumes, err := ctx.API.GetVolumes()
				if err != nil {
					errChan <- fmt.Errorf("unable to fetch volumes from the Scaleway API: %v", err)
					return
				}
				for _, val := range *volumes {
					creationDate, err := time.Parse("2006-01-02T15:04:05.000000+00:00", val.CreationDate)
					if err != nil {
						errChan <- fmt.Errorf("unable to parse creation date from the Scaleway API: %v", err)
						return
					}
					chEntries <- api.ScalewayImageInterface{
						Type:         "volume",
						CreationDate: creationDate,
						Identifier:   val.Identifier,
						Name:         val.Name,
						Tag:          "<volume>",
						VirtualSize:  val.Size,
						Public:       false,
						Organization: val.Organization,
						// FIXME the region should not be hardcoded
						Region: []string{"par1"},
					}
				}
			}()
		}
	}

	go func() {
		wg.Wait()
		close(chEntries)
	}()

	for {
		entry, ok := <-chEntries
		if !ok {
			break
		}
		entries = append(entries, entry)
	}
	select {
	case err := <-errChan:
		return err
	default:
		break
	}
	for key, value := range args.Filters {
		switch key {
		case "organization", "type", "name", "public":
			continue
		default:
			logrus.Warnf("Unknown filter: '%s=%s'", key, value)
		}
	}

	w := tabwriter.NewWriter(ctx.Stdout, 20, 1, 3, ' ', 0)
	defer w.Flush()
	if !args.Quiet {
		fmt.Fprintf(w, "REPOSITORY\tTAG\tIMAGE ID\tCREATED\tREGION\tARCH\n")
	}
	sort.Sort(api.ByCreationDate(entries))
	for _, image := range entries {
		if image.Identifier == "" {
			continue
		}
		for key, value := range args.Filters {
			switch key {
			case "type":
				if value != image.Type {
					goto skipimage
				}
			case "organization":
				switch value {
				case "me":
					value = ctx.API.Organization
				case "official-distribs":
					value = "a283af0b-d13e-42e1-a43f-855ffbf281ab"
				case "official-apps":
					value = "c3884e19-7a3e-4b69-9db8-50e7f902aafc"
				}
				if image.Organization != value {
					goto skipimage
				}
			case "name":
				if fuzzy.RankMatch(strings.ToLower(value), strings.ToLower(image.Name)) == -1 {
					goto skipimage
				}
			case "public":
				if (value == "true" && !image.Public) || (value == "false" && image.Public) {
					goto skipimage
				}
			}
		}

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
			var creationDate string
			if image.CreationDate.IsZero() {
				creationDate = "n/a"
			} else {
				creationDate = units.HumanDuration(time.Now().UTC().Sub(image.CreationDate))
			}
			if len(image.Archs) == 0 {
				image.Archs = []string{"n/a"}
			}
			sort.Strings(image.Archs)
			if len(image.Region) == 1 {
				image.Region = append(image.Region, "    ")
			}
			sort.Strings(image.Region)
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%v\n", shortName, tag, shortID, creationDate, image.Region, image.Archs)
		}

	skipimage:
		continue
	}
	return nil
}
