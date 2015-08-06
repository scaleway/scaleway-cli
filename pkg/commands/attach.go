// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"log"

	types "github.com/scaleway/scaleway-cli/pkg/commands/types"
	"github.com/scaleway/scaleway-cli/pkg/utils"
)

var cmdAttach = &types.Command{
	Exec:        cmdExecAttach,
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

// AttachArgs are flags for the `RunAttach` function
type AttachArgs struct {
	NoStdin bool
	Server  string
}

func cmdExecAttach(cmd *types.Command, rawArgs []string) {
	if attachHelp {
		cmd.PrintUsage()
	}
	if len(rawArgs) != 1 {
		cmd.PrintShortUsage()
	}

	args := AttachArgs{
		NoStdin: attachNoStdin,
		Server:  rawArgs[0],
	}
	ctx := cmd.GetContext(rawArgs)
	err := RunAttach(ctx, args)
	if err != nil {
		log.Fatalf("Cannot execute 'attach': %v", err)
	}
}

// RunAttach is the handler for 'scw attach'
func RunAttach(ctx types.CommandContext, args AttachArgs) error {
	serverID := ctx.API.GetServerID(args.Server)

	return utils.AttachToSerial(serverID, ctx.API.Token, !args.NoStdin)
}
