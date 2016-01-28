// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/scaleway/scaleway-cli/pkg/commands"
)

var cmdImages = &Command{
	Exec:        runImages,
	UsageLine:   "images [OPTIONS]",
	Description: "List images",
	Help:        "List images.",
	Examples: `
    $ scw images
    $ scw images -a
    $ scw images -q
    $ scw images --no-trunc
    $ scw images -f organization=me
    $ scw images -f organization=official-distribs
    $ scw images -f organization=official-apps
    $ scw images -f organization=UUIDOFORGANIZATION
    $ scw images -f name=ubuntu
    $ scw images -f type=image
    $ scw images -f type=bootscript
    $ scw images -f type=snapshot
    $ scw images -f type=volume
    $ scw images -f public=true
    $ scw images -f public=false
    $ scw images -f "organization=me type=volume" -q
`,
}

func init() {
	cmdImages.Flag.BoolVar(&imagesA, []string{"a", "-all"}, false, "Show all images")
	cmdImages.Flag.BoolVar(&imagesNoTrunc, []string{"-no-trunc"}, false, "Don't truncate output")
	cmdImages.Flag.BoolVar(&imagesQ, []string{"q", "-quiet"}, false, "Only show numeric IDs")
	cmdImages.Flag.BoolVar(&imagesHelp, []string{"h", "-help"}, false, "Print usage")
	cmdImages.Flag.StringVar(&imagesFilters, []string{"f", "-filter"}, "", "Filter output based on conditions provided")
}

// Flags
var imagesA bool         // -a flag
var imagesQ bool         // -q flag
var imagesNoTrunc bool   // -no-trunc flag
var imagesHelp bool      // -h, --help flag
var imagesFilters string // -f, --filters

func runImages(cmd *Command, rawArgs []string) error {
	if imagesHelp {
		return cmd.PrintUsage()
	}
	if len(rawArgs) != 0 {
		return cmd.PrintShortUsage()
	}

	args := commands.ImagesArgs{
		All:     imagesA,
		Quiet:   imagesQ,
		NoTrunc: imagesNoTrunc,
		Filters: make(map[string]string, 0),
	}
	if imagesFilters != "" {
		for _, filter := range strings.Split(imagesFilters, " ") {
			parts := strings.SplitN(filter, "=", 2)
			if len(parts) != 2 {
				logrus.Warnf("Invalid filter '%s', should be in the form 'key=value'", filter)
				continue
			}
			if _, ok := args.Filters[parts[0]]; ok {
				logrus.Warnf("Duplicated filter: %q", parts[0])
			} else {
				args.Filters[parts[0]] = parts[1]
			}
		}
	}
	ctx := cmd.GetContext(rawArgs)
	return commands.RunImages(ctx, args)
}
