// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/scaleway/scaleway-cli/pkg/commands"
)

var cmdPs = &Command{
	Exec:        runPs,
	UsageLine:   "ps [OPTIONS]",
	Description: "List servers",
	Help:        "List servers. By default, only running servers are displayed.",
	Examples: `
    $ scw ps
    $ scw ps -a
    $ scw ps -l
    $ scw ps -n=10
    $ scw ps -q
    $ scw ps --no-trunc
    $ scw ps -f state=booted
    $ scw ps -f state=running
    $ scw ps -f state=stopped
    $ scw ps -f ip=212.47.229.26
    $ scw ps -f tags=prod
    $ scw ps -f tags=boot=live
    $ scw ps -f image=docker
    $ scw ps -f image=alpine
    $ scw ps -f image=UUIDOFIMAGE
    $ scw ps -f arch=ARCH
    $ scw ps -f server-type=COMMERCIALTYPE
    $ scw ps -f "state=booted image=docker tags=prod"
`,
}

func init() {
	cmdPs.Flag.BoolVar(&psA, []string{"a", "-all"}, false, "Show all servers. Only running servers are shown by default")
	cmdPs.Flag.BoolVar(&psL, []string{"l", "-latest"}, false, "Show only the latest created server, include non-running ones")
	cmdPs.Flag.IntVar(&psN, []string{"n"}, 0, "Show n last created servers, include non-running ones")
	cmdPs.Flag.BoolVar(&psNoTrunc, []string{"-no-trunc"}, false, "Don't truncate output")
	cmdPs.Flag.BoolVar(&psQ, []string{"q", "-quiet"}, false, "Only display numeric IDs")
	cmdPs.Flag.BoolVar(&psHelp, []string{"h", "-help"}, false, "Print usage")
	cmdPs.Flag.StringVar(&psFilters, []string{"f", "-filter"}, "", "Filter output based on conditions provided")
}

// Flags
var psA bool         // -a flag
var psL bool         // -l flag
var psQ bool         // -q flag
var psNoTrunc bool   // -no-trunc flag
var psN int          // -n flag
var psHelp bool      // -h, --help flag
var psFilters string // -f, --filter flag

func runPs(cmd *Command, rawArgs []string) error {
	if psHelp {
		return cmd.PrintUsage()
	}
	if len(rawArgs) != 0 {
		return cmd.PrintShortUsage()
	}

	args := commands.PsArgs{
		All:     psA,
		Latest:  psL,
		Quiet:   psQ,
		NoTrunc: psNoTrunc,
		NLast:   psN,
		Filters: make(map[string]string, 0),
	}
	if psFilters != "" {
		for _, filter := range strings.Split(psFilters, " ") {
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
	return commands.RunPs(ctx, args)
}
