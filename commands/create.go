// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"

	log "github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"

	"github.com/scaleway/scaleway-cli/api"
	types "github.com/scaleway/scaleway-cli/commands/types"
)

var cmdCreate = &types.Command{
	Exec:        runCreate,
	UsageLine:   "create [OPTIONS] IMAGE",
	Description: "Create a new server but do not start it",
	Help:        "Create a new server but do not start it.",
	Examples: `
    $ scw create docker
    $ scw create 10GB
    $ scw create --bootscript=3.2.34 --env="boot=live rescue_image=http://j.mp/scaleway-ubuntu-trusty-tarball" 50GB
    $ scw inspect $(scw create 1GB --bootscript=rescue --volume=50GB)
    $ scw create $(scw tag my-snapshot my-image)
`,
}

func init() {
	cmdCreate.Flag.StringVar(&createName, []string{"-name"}, "", "Assign a name")
	cmdCreate.Flag.StringVar(&createBootscript, []string{"-bootscript"}, "", "Assign a bootscript")
	cmdCreate.Flag.StringVar(&createEnv, []string{"e", "-env"}, "", "Provide metadata tags passed to initrd (i.e., boot=resue INITRD_DEBUG=1)")
	cmdCreate.Flag.StringVar(&createVolume, []string{"v", "-volume"}, "", "Attach additional volume (i.e., 50G)")
	cmdCreate.Flag.BoolVar(&createHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var createName string       // --name flag
var createBootscript string // --bootscript flag
var createEnv string        // -e, --env flag
var createVolume string     // -v, --volume flag
var createHelp bool         // -h, --help flag

func runCreate(cmd *types.Command, args []string) {
	if createHelp {
		cmd.PrintUsage()
	}
	if len(args) != 1 {
		cmd.PrintShortUsage()
	}

	serverID, err := api.CreateServer(cmd.API, args[0], createName, createBootscript, createEnv, createVolume, true)

	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	fmt.Println(serverID)
}
