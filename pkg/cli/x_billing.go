// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"fmt"
	"math/big"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/go-units"
	"github.com/scaleway/scaleway-cli/pkg/commands"
	"github.com/scaleway/scaleway-cli/pkg/pricing"
	"github.com/scaleway/scaleway-cli/pkg/utils"
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

	logrus.Warn("")
	logrus.Warn("Warning: 'scw _billing' is a work-in-progress price estimation tool")
	logrus.Warn("For real usage, visit https://cloud.scaleway.com/#/billing")
	logrus.Warn("")

	// table
	w := tabwriter.NewWriter(ctx.Stdout, 20, 1, 3, ' ', 0)
	defer w.Flush()
	fmt.Fprintf(w, "ID\tNAME\tSTARTED\tMONTH PRICE\n")

	// servers
	servers, err := cmd.API.GetServers(true, 0)
	if err != nil {
		return err
	}

	totalMonthPrice := new(big.Rat)

	for _, server := range *servers {
		if server.State != "running" {
			continue
		}
		commercialType := strings.ToLower(server.CommercialType)
		shortID := utils.TruncIf(server.Identifier, 8, !args.NoTrunc)
		shortName := utils.TruncIf(utils.Wordify(server.Name), 25, !args.NoTrunc)
		modificationTime, _ := time.Parse("2006-01-02T15:04:05.000000+00:00", server.ModificationDate)
		modificationAgo := time.Now().UTC().Sub(modificationTime)
		shortModificationDate := units.HumanDuration(modificationAgo)
		usage := pricing.NewUsageByPath(fmt.Sprintf("/compute/%s/run", commercialType))
		usage.SetStartEnd(modificationTime, time.Now().UTC())

		totalMonthPrice = totalMonthPrice.Add(totalMonthPrice, usage.Total())

		fmt.Fprintf(w, "server/%s/%s\t%s\t%s\t%s\n", commercialType, shortID, shortName, shortModificationDate, usage.TotalString())
	}

	fmt.Fprintf(w, "TOTAL\t\t\t%s\n", pricing.PriceString(totalMonthPrice, "EUR"))

	return nil
}
