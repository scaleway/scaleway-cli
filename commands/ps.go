// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	log "github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"
	"github.com/scaleway/scaleway-cli/vendor/github.com/docker/docker/pkg/units"

	types "github.com/scaleway/scaleway-cli/commands/types"
	"github.com/scaleway/scaleway-cli/utils"
)

var cmdPs = &types.Command{
	Exec:        runPs,
	UsageLine:   "ps [OPTIONS]",
	Description: "List servers",
	Help:        "List servers. By default, only running servers are displayed.",
}

func init() {
	cmdPs.Flag.BoolVar(&psA, []string{"a", "-all"}, false, "Show all servers. Only running servers are shown by default")
	cmdPs.Flag.BoolVar(&psL, []string{"l", "-latest"}, false, "Show only the latest created server, include non-running ones")
	cmdPs.Flag.IntVar(&psN, []string{"n"}, 0, "Show n last created servers, include non-running ones")
	cmdPs.Flag.BoolVar(&psNoTrunc, []string{"-no-trunc"}, false, "Don't truncate output")
	cmdPs.Flag.BoolVar(&psQ, []string{"q", "-quiet"}, false, "Only display numeric IDs")
	cmdPs.Flag.BoolVar(&psHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var psA bool       // -a flag
var psL bool       // -l flag
var psQ bool       // -q flag
var psNoTrunc bool // -no-trunc flag
var psN int        // -n flag
var psHelp bool    // -h, --help flag

func runPs(cmd *types.Command, args []string) {
	if psHelp {
		cmd.PrintUsage()
	}
	if len(args) != 0 {
		cmd.PrintShortUsage()
	}

	limit := psN
	if psL {
		limit = 1
	}
	servers, err := cmd.API.GetServers(psA || psN > 0 || psL, limit)
	if err != nil {
		log.Fatalf("Unable to fetch servers from the Scaleway API: %v", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	defer w.Flush()
	if !psQ {
		fmt.Fprintf(w, "SERVER ID\tIMAGE\tCOMMAND\tCREATED\tSTATUS\tPORTS\tNAME\n")
	}
	for _, server := range *servers {
		if psQ {
			fmt.Fprintf(w, "%s\n", server.Identifier)
		} else {
			shortID := utils.TruncIf(server.Identifier, 8, !psNoTrunc)
			shortImage := utils.TruncIf(utils.Wordify(server.Image.Name), 25, !psNoTrunc)
			shortName := utils.TruncIf(utils.Wordify(server.Name), 25, !psNoTrunc)
			creationTime, _ := time.Parse("2006-01-02T15:04:05.000000+00:00", server.CreationDate)
			shortCreationDate := units.HumanDuration(time.Now().UTC().Sub(creationTime))
			port := server.PublicAddress.IP
			fmt.Fprintf(w, "%s\t%s\t\t%s\t%s\t%s\t%s\n", shortID, shortImage, shortCreationDate, server.State, port, shortName)
		}
	}
}
