// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	log "github.com/Sirupsen/logrus"

	types "github.com/scaleway/scaleway-cli/commands/types"
	"github.com/scaleway/scaleway-cli/utils"
)

var cmdAttach = &types.Command{
	Exec:        runAttach,
	UsageLine:   "attach [OPTIONS] SERVER",
	Description: "Attach to a server serial console",
	Help:        "Attach to a running server serial console.",
	Examples: `
    $ scw attach my-running-server
    $ scw attach $(scw start my-stopped-server)
    $ scw attach $(scw start $(scw create ubuntu-vivid))
`,
}

func init() {
	cmdAttach.Flag.BoolVar(&attachHelp, []string{"h", "-help"}, false, "Print usage")
	cmdAttach.Flag.BoolVar(&attachNoStdin, []string{"-no-stdin"}, false, "Do not attach stdin")
}

// Flags
var attachHelp bool    // -h, --help flag
var attachNoStdin bool // --no-stdin flag

func runAttach(cmd *types.Command, args []string) {
	if attachHelp {
		cmd.PrintUsage()
	}
	if len(args) != 1 {
		cmd.PrintShortUsage()
	}

	serverID := cmd.API.GetServerID(args[0])

	err := utils.AttachToSerial(serverID, cmd.API.Token, !attachNoStdin)
	if err != nil {
		log.Fatalf("%v", err)
	}
}
