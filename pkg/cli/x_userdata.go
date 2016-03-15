// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/scwversion"
)

var cmdUserdata = &Command{
	Exec:        runUserdata,
	UsageLine:   "_userdata [OPTIONS] SERVER [FIELD[=VALUE]]",
	Description: "",
	Hidden:      true,
	Help:        "List, read and write and delete server's userdata",
	Examples: `
    $ scw _userdata myserver
    $ scw _userdata myserver key
    $ scw _userdata myserver key=value
    $ scw _userdata myserver key=""
`,
}

func init() {
	cmdUserdata.Flag.BoolVar(&userdataHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var userdataHelp bool // -h, --help flag

func runUserdata(cmd *Command, args []string) error {
	if userdataHelp {
		return cmd.PrintUsage()
	}
	if len(args) < 1 {
		return cmd.PrintShortUsage()
	}
	metadata := false
	ctx := cmd.GetContext(args)
	var API *api.ScalewayAPI
	var err error
	var serverID string
	if args[0] == "local" {
		API, err = api.NewScalewayAPI("", "", scwversion.UserAgent())
		if err != nil {
			return err
		}
		metadata = true
	} else {
		if ctx.API == nil {
			return fmt.Errorf("You need to login first: 'scw login'")
		}
		serverID, err = ctx.API.GetServerID(args[0])
		if err != nil {
			return err
		}
		API = ctx.API
	}

	switch len(args) {
	case 1:
		// List userdata
		res, err := API.GetUserdatas(serverID, metadata)
		if err != nil {
			return err
		}
		for _, key := range res.UserData {
			fmt.Fprintln(ctx.Stdout, key)
		}
	default:
		parts := strings.Split(args[1], "=")
		key := parts[0]
		switch len(parts) {
		case 1:
			// Get userdatas
			res, err := API.GetUserdata(serverID, key, metadata)
			if err != nil {
				return err
			}
			fmt.Fprintf(ctx.Stdout, "%s\n", res.String())
		default:
			value := parts[1]
			if value != "" {
				var data []byte
				// Set userdata
				if value[0] == '@' {
					data, err = ioutil.ReadFile(value[1:])
					if err != nil {
						return err
					}
				} else {
					data = []byte(value)
				}
				err := API.PatchUserdata(serverID, key, data, metadata)
				if err != nil {
					return err
				}
				fmt.Fprintln(ctx.Stdout, key)
			} else {
				// Delete userdata
				err := API.DeleteUserdata(serverID, key, metadata)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
