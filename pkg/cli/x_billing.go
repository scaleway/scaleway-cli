// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import "fmt"

var cmdBilling = &Command{
	Exec:        runBilling,
	UsageLine:   "_billing [OPTIONS]",
	Description: "",
	Hidden:      true,
	Help:        "Get resources billing estimation",
}

func init() {
	cmdBilling.Flag.BoolVar(&billingHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var billingHelp bool // -h, --help flag

func runBilling(cmd *Command, args []string) error {
	if billingHelp {
		return cmd.PrintUsage()
	}
	if len(args) > 0 {
		return cmd.PrintShortUsage()
	}

	fmt.Println("BILLING")

	return nil
}
