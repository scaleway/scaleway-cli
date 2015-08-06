// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	log "github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"
	"github.com/scaleway/scaleway-cli/vendor/github.com/skratchdot/open-golang/open"

	"github.com/scaleway/scaleway-cli/pkg/api"
)

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
}

// Flags
var inspectFormat string // -f, --format flag
var inspectBrowser bool  // -b, --browser flag
var inspectHelp bool     // -h, --help flag

func runInspect(cmd *Command, args []string) {
	if inspectHelp {
		cmd.PrintUsage()
	}
	if len(args) < 1 {
		cmd.PrintShortUsage()
	}

	nbInspected := 0
	ci := make(chan api.ScalewayResolvedIdentifier)
	cj := make(chan api.InspectIdentifierResult)
	go api.ResolveIdentifiers(cmd.API, args, ci)
	go api.InspectIdentifiers(cmd.API, ci, cj)

	if inspectBrowser {
		// --browser will open links in the browser
		for {
			data, isOpen := <-cj
			if !isOpen {
				break
			}

			switch data.Type {
			case api.IdentifierServer:
				err := open.Start(fmt.Sprintf("https://cloud.scaleway.com/#/servers/%s", data.Object.(*api.ScalewayServer).Identifier))
				if err != nil {
					log.Fatalf("Cannot open browser: %v", err)
				}
				nbInspected++
			case api.IdentifierImage:
				err := open.Start(fmt.Sprintf("https://cloud.scaleway.com/#/images/%s", data.Object.(*api.ScalewayImage).Identifier))
				if err != nil {
					log.Fatalf("Cannot open browser: %v", err)
				}
				nbInspected++
			case api.IdentifierVolume:
				err := open.Start(fmt.Sprintf("https://cloud.scaleway.com/#/volumes/%s", data.Object.(*api.ScalewayVolume).Identifier))
				if err != nil {
					log.Fatalf("Cannot open browser: %v", err)
				}
				nbInspected++
			case api.IdentifierSnapshot:
				log.Errorf("Cannot use '--browser' option for snapshots")
			case api.IdentifierBootscript:
				log.Errorf("Cannot use '--browser' option for bootscripts")
			}
		}

	} else {
		// without --browser option, inspect will print object info to the terminal
		res := "["
		for {
			data, isOpen := <-cj
			if !isOpen {
				break
			}
			if inspectFormat == "" {
				dataB, err := json.MarshalIndent(data.Object, "", "  ")
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

				err = tmpl.Execute(os.Stdout, data.Object)
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
	}

	if len(args) != nbInspected {
		os.Exit(1)
	}
}
