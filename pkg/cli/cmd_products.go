// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import "github.com/scaleway/scaleway-cli/pkg/commands"

var cmdProducts = &Command{
	Exec:        runProducts,
	UsageLine:   "products [OPTIONS] PRODUCT",
	Description: "Display product PRODUCT information",
	Help:        "Display products PRODUCT information. At the moment only `servers` is supported.",
	Examples: `
    $ scw products servers
    $ scw products --short servers
`,
}

func init() {
	cmdProducts.Flag.BoolVar(&productsHelp, []string{"h", "-help"}, false, "Print usage")
	cmdProducts.Flag.BoolVar(&productsShort, []string{"s", "-short"}, false, "Print only commercial names")
}

// Flags
var productsHelp bool  // -h, --help flag
var productsShort bool // -s, --short flag

func runProducts(cmd *Command, rawArgs []string) error {
	if productsHelp {
		return cmd.PrintUsage()
	}
	if len(rawArgs) != 1 {
		return cmd.PrintShortUsage()
	}

	args := commands.ProductsArgs{
		Short:    productsShort,
		Products: rawArgs,
	}

	ctx := cmd.GetContext(rawArgs)
	return commands.RunProducts(ctx, args)
}
