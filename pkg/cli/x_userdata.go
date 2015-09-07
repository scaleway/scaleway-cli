// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"fmt"
	"strings"

	"github.com/scaleway/scaleway-cli/pkg/api"
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

	ctx := cmd.GetContext(args)
	var Api *api.ScalewayAPI
	var err error
	var serverID string
	if args[0] == "local" {
		Api, err = api.NewScalewayAPI("", "", "", "")
		if err != nil {
			return err
		}
		Api.EnableMetadataAPI()
	} else {
		if ctx.API == nil {
			return fmt.Errorf("You need to login first: 'scw login'")
		}
		serverID = ctx.API.GetServerID(args[0])
		Api = ctx.API
	}

	switch len(args) {
	case 1:
		// List userdata
		res, err := Api.GetUserdatas(serverID)
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
			res, err := Api.GetUserdata(serverID, key)
			if err != nil {
				return err
			}
			fmt.Fprintf(ctx.Stdout, "%s\n", res.String())
		default:
			value := parts[1]
			if value != "" {
				// Set userdata
				err := Api.PatchUserdata(serverID, key, []byte(value))
				if err != nil {
					return err
				}
				fmt.Fprintln(ctx.Stdout, key)
			} else {
				// Delete userdata
				err := Api.DeleteUserdata(serverID, key)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
