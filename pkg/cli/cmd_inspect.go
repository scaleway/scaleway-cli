// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import "github.com/scaleway/scaleway-cli/pkg/commands"

var cmdInspect = &Command{
	Exec:        runInspect,
	UsageLine:   "inspect [OPTIONS] IDENTIFIER [IDENTIFIER...]",
	Description: "Return low-level information on a server, image, snapshot, volume or bootscript",
	Help:        "Return low-level information on a server, image, snapshot, volume or bootscript.",
	Examples: `
    $ scw inspect my-server
    $ scw inspect server:my-server
    $ scw inspect --browser my-server
    $ scw inspect a-public-image
    $ scw inspect image:a-public-image
    $ scw inspect my-snapshot
    $ scw inspect snapshot:my-snapshot
    $ scw inspect my-volume
    $ scw inspect volume:my-volume
    $ scw inspect my-image
    $ scw inspect image:my-image
    $ scw inspect my-server | jq '.[0].public_ip.address'
    $ scw inspect $(scw inspect my-image | jq '.[0].root_volume.id')
    $ scw inspect -f "{{ .PublicAddress.IP }}" my-server
    $ scw --sensitive inspect my-server
`,
}

func init() {
	cmdInspect.Flag.BoolVar(&inspectHelp, []string{"h", "-help"}, false, "Print usage")
	cmdInspect.Flag.StringVar(&inspectFormat, []string{"f", "-format"}, "", "Format the output using the given go template")
	cmdInspect.Flag.BoolVar(&inspectBrowser, []string{"b", "-browser"}, false, "Inspect object in browser")
	cmdInspect.Flag.StringVar(&inspectArch, []string{"-arch"}, "*", "Specify architecture")
}

// Flags
var inspectFormat string // -f, --format flag
var inspectBrowser bool  // -b, --browser flag
var inspectHelp bool     // -h, --help flag
var inspectArch string   // --arch flag

func runInspect(cmd *Command, rawArgs []string) error {
	if inspectHelp {
		return cmd.PrintUsage()
	}
	if len(rawArgs) < 1 {
		return cmd.PrintShortUsage()
	}

	args := commands.InspectArgs{
		Format:      inspectFormat,
		Browser:     inspectBrowser,
		Identifiers: rawArgs,
		Arch:        inspectArch,
	}
	ctx := cmd.GetContext(rawArgs)
	return commands.RunInspect(ctx, args)
}
