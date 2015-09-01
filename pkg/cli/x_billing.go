// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"fmt"
	"text/tabwriter"
	"time"

	"github.com/scaleway/scaleway-cli/pkg/commands"
	"github.com/scaleway/scaleway-cli/pkg/pricing"
	"github.com/scaleway/scaleway-cli/pkg/utils"
	"github.com/scaleway/scaleway-cli/vendor/github.com/docker/docker/pkg/units"
)

var cmdBilling = &Command{
	Exec:        runBilling,
	UsageLine:   "_billing [OPTIONS]",
	Description: "",
	Hidden:      true,
	Help:        "Get resources billing estimation",
}

func init() {
	cmdBilling.Flag.BoolVar(&billingHelp, []string{"h", "-help"}, false, "Print usage")
	cmdBilling.Flag.BoolVar(&billingNoTrunc, []string{"-no-trunc"}, false, "Don't truncate output")
}

// BillingArgs are flags for the `RunBilling` function
type BillingArgs struct {
	NoTrunc bool
}

// Flags
var billingHelp bool    // -h, --help flag
var billingNoTrunc bool // --no-trunc flag

func runBilling(cmd *Command, rawArgs []string) error {
	if billingHelp {
		return cmd.PrintUsage()
	}
	if len(rawArgs) > 0 {
		return cmd.PrintShortUsage()
	}

	// cli parsing
	args := commands.PsArgs{
		NoTrunc: billingNoTrunc,
	}
	ctx := cmd.GetContext(rawArgs)

	// table
	w := tabwriter.NewWriter(ctx.Stdout, 20, 1, 3, ' ', 0)
	defer w.Flush()
	fmt.Fprintf(w, "ID\tNAME\tSTARTED\tMONTH PRICE\n")

	// servers
	servers, err := cmd.API.GetServers(true, 0)
	if err != nil {
		return err
	}
	for _, server := range *servers {
		if server.State != "running" {
			continue
		}
		shortID := utils.TruncIf(server.Identifier, 8, !args.NoTrunc)
		shortName := utils.TruncIf(utils.Wordify(server.Name), 25, !args.NoTrunc)
		modificationTime, _ := time.Parse("2006-01-02T15:04:05.000000+00:00", server.ModificationDate)
		modificationAgo := time.Now().UTC().Sub(modificationTime)
		shortModificationDate := units.HumanDuration(modificationAgo)
		usage := pricing.NewUsageByPath("/compute/c1/run")
		usage.SetStartEnd(modificationTime, time.Now().UTC())

		fmt.Fprintf(w, "server/%s\t%s\t%s\t%s\n", shortID, shortName, shortModificationDate, usage.TotalString())
	}

	return nil
}
