// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import "fmt"

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

	fmt.Println("Not implemented")

	return nil
}
