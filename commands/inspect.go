// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	log "github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"

	"github.com/scaleway/scaleway-cli/api"
	types "github.com/scaleway/scaleway-cli/commands/types"
)

var cmdInspect = &types.Command{
	Exec:        runInspect,
	UsageLine:   "inspect [OPTIONS] IDENTIFIER [IDENTIFIER...]",
	Description: "Return low-level information on a server, image, snapshot, volume or bootscript",
	Help:        "Return low-level information on a server, image, snapshot, volume or bootscript.",
	Examples: `
    $ scw inspect my-server
    $ scw inspect server:my-server
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
	cmdInspect.Flag.StringVar(&inspectFormat, []string{"f", "-format"}, "", "Format the output using the given go template.")
}

// Flags
var inspectFormat string // -f, --format flat
var inspectHelp bool     // -h, --help flag

func runInspect(cmd *types.Command, args []string) {
	if inspectHelp {
		cmd.PrintUsage()
	}
	if len(args) < 1 {
		cmd.PrintShortUsage()
	}

	res := "["
	nbInspected := 0
	ci := make(chan api.ScalewayResolvedIdentifier)
	cj := make(chan interface{})
	go api.ResolveIdentifiers(cmd.API, args, ci)
	go api.InspectIdentifiers(cmd.API, ci, cj)
	for {
		data, open := <-cj
		if !open {
			break
		}
		if inspectFormat == "" {
			dataB, err := json.MarshalIndent(data, "", "  ")
			if err == nil {
				if nbInspected != 0 {
					res += ",\n"
				}
				res += string(dataB)
				nbInspected++
			}
		} else {
			tmpl, err := template.New("").Funcs(api.FuncMap).Parse(inspectFormat)
			if err != nil {
				log.Fatalf("Format parsing error: %v", err)
			}

			err = tmpl.Execute(os.Stdout, data)
			if err != nil {
				log.Fatalf("Format execution error: %v", err)
			}
			fmt.Fprint(os.Stdout, "\n")
			nbInspected++
		}
	}
	res += "]"

	if inspectFormat == "" {
		if os.Getenv("SCW_SENSITIVE") != "1" {
			res = cmd.API.HideAPICredentials(res)
		}
		fmt.Println(res)
	}

	if len(args) != nbInspected {
		os.Exit(1)
	}
}
