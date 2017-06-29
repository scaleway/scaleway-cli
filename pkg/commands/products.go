// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"sort"

	"github.com/dustin/go-humanize"
	"github.com/scaleway/scaleway-cli/pkg/api"
)

// ProductsArgs are flags for the `RunProducts` function
type ProductsArgs struct {
	Short    bool
	Products []string
}

// DisplayServerFunc
type displayServerFunc func(CommandContext, *api.ScalewayProductsServers)

// RunProducts is the handler for 'scw products'
func RunProducts(ctx CommandContext, args ProductsArgs) error {
	for _, product := range args.Products {
		switch product {
		case "servers":
			products, err := ctx.API.GetProductsServers()
			if err != nil {
				return fmt.Errorf("Unable to fetch products from the Scaleway API: %v", err)
			}

			var displayFunc displayServerFunc = DisplayServerFull
			if args.Short {
				displayFunc = DisplayServerShort
			}
			displayFunc(ctx, products)
		default:
			return fmt.Errorf("Unknow product '%v'", product)
		}
	}

	return nil
}

// DisplayServerShort only display the product name
func DisplayServerShort(ctx CommandContext, products *api.ScalewayProductsServers) {
	for name := range products.Servers {
		fmt.Fprintf(ctx.Stdout, "%v\n", name)
	}
}

// DisplayServerFull only display the server product information
func DisplayServerFull(ctx CommandContext, products *api.ScalewayProductsServers) {
	fmt.Fprintf(ctx.Stdout, "%-12s %8s %8s %8s %10s\n",
		"COMMERCIAL TYPE", "ARCH", "CPUs", "RAM", "BAREMETAL")

	names := make([]string, len(products.Servers))
	i := 0
	for k := range products.Servers {
		names[i] = k
		i++
	}
	sort.Strings(names)

	for _, name := range names {
		offer := products.Servers[name]

		fmt.Fprintf(ctx.Stdout, "%-15s %8s %8d %8s %10t\n",
			name, offer.Arch, offer.Ncpus, humanize.Bytes(offer.Ram), offer.Baremetal)
	}
}
