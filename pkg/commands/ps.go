// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"text/tabwriter"
	"time"

	"github.com/scaleway/scaleway-cli/vendor/github.com/docker/docker/pkg/units"

	"github.com/scaleway/scaleway-cli/pkg/utils"
)

// PsArgs are flags for the `RunPs` function
type PsArgs struct {
	All     bool
	Latest  bool
	NLast   int
	NoTrunc bool
	Quiet   bool
}

// RunPs is the handler for 'scw ps'
func RunPs(ctx CommandContext, args PsArgs) error {
	limit := args.NLast
	if args.Latest {
		limit = 1
	}
	all := args.All || args.NLast > 0 || args.Latest
	servers, err := ctx.API.GetServers(all, limit)
	if err != nil {
		return fmt.Errorf("Unable to fetch servers from the Scaleway API: %v", err)
	}

	w := tabwriter.NewWriter(ctx.Stdout, 20, 1, 3, ' ', 0)
	defer w.Flush()
	if !args.Quiet {
		fmt.Fprintf(w, "SERVER ID\tIMAGE\tCOMMAND\tCREATED\tSTATUS\tPORTS\tNAME\n")
	}
	for _, server := range *servers {
		if args.Quiet {
			fmt.Fprintf(w, "%s\n", server.Identifier)
		} else {
			shortID := utils.TruncIf(server.Identifier, 8, !args.NoTrunc)
			shortImage := utils.TruncIf(utils.Wordify(server.Image.Name), 25, !args.NoTrunc)
			shortName := utils.TruncIf(utils.Wordify(server.Name), 25, !args.NoTrunc)
			creationTime, _ := time.Parse("2006-01-02T15:04:05.000000+00:00", server.CreationDate)
			shortCreationDate := units.HumanDuration(time.Now().UTC().Sub(creationTime))
			port := server.PublicAddress.IP
			fmt.Fprintf(w, "%s\t%s\t\t%s\t%s\t%s\t%s\n", shortID, shortImage, shortCreationDate, server.State, port, shortName)
		}
	}
	return nil
}
