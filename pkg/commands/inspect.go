// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/Sirupsen/logrus"
	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/skratchdot/open-golang/open"
)

// InspectArgs are flags for the `RunInspect` function
type InspectArgs struct {
	Format      string
	Browser     bool
	Identifiers []string
	Arch        string
}

// RunInspect is the handler for 'scw inspect'
func RunInspect(ctx CommandContext, args InspectArgs) error {
	nbInspected := 0
	ci := make(chan api.ScalewayResolvedIdentifier)
	cj := make(chan api.InspectIdentifierResult)
	go api.ResolveIdentifiers(ctx.API, args.Identifiers, ci)
	go api.InspectIdentifiers(ctx.API, ci, cj, args.Arch)

	if args.Browser {
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
					return fmt.Errorf("cannot open browser: %v", err)
				}
				nbInspected++
			case api.IdentifierImage:
				err := open.Start(fmt.Sprintf("https://cloud.scaleway.com/#/images/%s", data.Object.(*api.ScalewayImage).Identifier))
				if err != nil {
					return fmt.Errorf("cannot open browser: %v", err)
				}
				nbInspected++
			case api.IdentifierVolume:
				err := open.Start(fmt.Sprintf("https://cloud.scaleway.com/#/volumes/%s", data.Object.(*api.ScalewayVolume).Identifier))
				if err != nil {
					return fmt.Errorf("cannot open browser: %v", err)
				}
				nbInspected++
			case api.IdentifierSnapshot:
				logrus.Errorf("Cannot use '--browser' option for snapshots")
			case api.IdentifierBootscript:
				logrus.Errorf("Cannot use '--browser' option for bootscripts")
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
			if args.Format == "" {
				dataB, err := json.MarshalIndent(data.Object, "", "  ")
				if err == nil {
					if nbInspected != 0 {
						res += ",\n"
					}
					res += string(dataB)
					nbInspected++
				}
			} else {
				tmpl, err := template.New("").Funcs(api.FuncMap).Parse(args.Format)
				if err != nil {
					return fmt.Errorf("format parsing error: %v", err)
				}

				err = tmpl.Execute(ctx.Stdout, data.Object)
				if err != nil {
					return fmt.Errorf("format execution error: %v", err)
				}
				fmt.Fprint(ctx.Stdout, "\n")
				nbInspected++
			}
		}
		res += "]"

		if args.Format == "" {
			if ctx.Getenv("SCW_SENSITIVE") != "1" {
				res = ctx.API.HideAPICredentials(res)
			}
			if len(res) > 2 {
				fmt.Fprintln(ctx.Stdout, res)
			}
		}
	}

	if len(args.Identifiers) != nbInspected {
		return fmt.Errorf("at least 1 item failed to be inspected")
	}
	return nil
}
