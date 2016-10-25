// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/go-units"
	"github.com/renstrom/fuzzysearch/fuzzy"

	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/utils"
)

// PsArgs are flags for the `RunPs` function
type PsArgs struct {
	NLast   int
	All     bool
	Latest  bool
	NoTrunc bool
	Quiet   bool
	Filters map[string]string
}

// RunPs is the handler for 'scw ps'
func RunPs(ctx CommandContext, args PsArgs) error {
	limit := args.NLast
	if args.Latest {
		limit = 1
	}

	filterState := args.Filters["state"]

	// FIXME: if filter state is defined, try to optimize the query
	all := args.All || limit > 0 || filterState != ""
	servers, err := ctx.API.GetServers(all, 0)
	if err != nil {
		return fmt.Errorf("Unable to fetch servers from the Scaleway API: %v", err)
	}

	for key, value := range args.Filters {
		switch key {
		case "state", "name", "tags", "image", "ip", "arch", "server-type", "zone":
			continue
		default:
			logrus.Warnf("Unknown filter: '%s=%s'", key, value)
		}
	}
	filtered := make([]api.ScalewayServer, 0, len(*servers))
	for _, server := range *servers {
		// filtering
		for key, value := range args.Filters {
			switch key {
			case "state":
				if value != server.State {
					goto skipServer
				}
			case "name":
				if fuzzy.RankMatch(strings.ToLower(value), strings.ToLower(server.Name)) == -1 {
					goto skipServer
				}
			case "tags":
				found := false
				for _, tag := range server.Tags {
					if tag == value {
						found = true
						continue
					}
				}
				if !found {
					goto skipServer
				}
			case "image":
				imageID, err := ctx.API.GetImageID(value, "*")
				if err != nil {
					goto skipServer
				}
				if imageID.Identifier != server.Image.Identifier {
					goto skipServer
				}
			case "ip":
				if value != server.PublicAddress.IP {
					goto skipServer
				}
			case "arch":
				if value != server.Arch {
					goto skipServer
				}
			case "server-type":
				if value != server.CommercialType {
					goto skipServer
				}
			case "zone":
				if value != server.Location.ZoneID {
					goto skipServer
				}
			}
		}
		filtered = append(filtered, server)
	skipServer:
		continue
	}
	w := tabwriter.NewWriter(ctx.Stdout, 20, 1, 3, ' ', 0)
	defer w.Flush()
	if !args.Quiet {
		fmt.Fprintf(w, "SERVER ID\tIMAGE\tZONE\tCREATED\tSTATUS\tPORTS\tNAME\tCOMMERCIAL TYPE\n")
	}
	sort.Sort(api.ScalewaySortServers(filtered))
	for i, server := range filtered {
		if limit > 0 && i >= limit {
			break
		}
		if args.Quiet {
			fmt.Fprintf(w, "%s\n", server.Identifier)
		} else {
			shortID := utils.TruncIf(server.Identifier, 8, !args.NoTrunc)
			shortImage := utils.TruncIf(utils.Wordify(server.Image.Name), 25, !args.NoTrunc)
			shortName := utils.TruncIf(utils.Wordify(server.Name), 25, !args.NoTrunc)
			creationTime, _ := time.Parse("2006-01-02T15:04:05.000000+00:00", server.CreationDate)
			shortCreationDate := units.HumanDuration(time.Now().UTC().Sub(creationTime))
			port := server.PublicAddress.IP
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", shortID, shortImage, server.Location.ZoneID, shortCreationDate, server.State, port, shortName, server.CommercialType)
		}
	}
	return nil
}
