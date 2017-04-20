// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"strings"

	"github.com/scaleway/scaleway-cli/pkg/commands"
)

var cmdCreate = &Command{
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
    $ scw create --tmp-ssh-key 10GB
`,
}

func init() {
	cmdCreate.Flag.StringVar(&createName, []string{"-name"}, "", "Assign a name")
	cmdCreate.Flag.StringVar(&createBootscript, []string{"-bootscript"}, "", "Assign a bootscript")
	cmdCreate.Flag.StringVar(&createEnv, []string{"e", "-env"}, "", "Provide metadata tags passed to initrd (i.e., boot=rescue INITRD_DEBUG=1)")
	cmdCreate.Flag.StringVar(&createVolume, []string{"v", "-volume"}, "", "Attach additional volume (i.e., 50G)")
	cmdCreate.Flag.StringVar(&createIPAddress, []string{"-ip-address"}, "dynamic", "Assign a reserved public IP, a 'dynamic' one or 'none'")
	cmdCreate.Flag.StringVar(&createCommercialType, []string{"-commercial-type"}, "X64-2GB", "Create a server with specific commercial-type C1, C2[S|M|L], X64-[2|4|8|15|30|60|120]GB, ARM64-[2|4|8]GB")
	cmdCreate.Flag.BoolVar(&createHelp, []string{"h", "-help"}, false, "Print usage")
	cmdCreate.Flag.BoolVar(&createIPV6, []string{"-ipv6"}, false, "Enable IPV6")
	cmdCreate.Flag.BoolVar(&createTmpSSHKey, []string{"-tmp-ssh-key"}, false, "Access your server without uploading your SSH key to your account")
}

// Flags
var createName string           // --name flag
var createBootscript string     // --bootscript flag
var createEnv string            // -e, --env flag
var createVolume string         // -v, --volume flag
var createHelp bool             // -h, --help flag
var createTmpSSHKey bool        // --tmp-ssh-key flag
var createIPAddress string      // --ip-address flag
var createCommercialType string // --commercial-type flag
var createIPV6 bool             // --ipv6 flag

func runCreate(cmd *Command, rawArgs []string) error {
	if createHelp {
		return cmd.PrintUsage()
	}
	if len(rawArgs) != 1 {
		return cmd.PrintShortUsage()
	}

	args := commands.CreateArgs{
		Name:           createName,
		Bootscript:     createBootscript,
		Image:          rawArgs[0],
		TmpSSHKey:      createTmpSSHKey,
		IP:             createIPAddress,
		CommercialType: createCommercialType,
		IPV6:           createIPV6,
	}

	if len(createEnv) > 0 {
		args.Tags = strings.Split(createEnv, " ")
	}
	if len(createVolume) > 0 {
		args.Volumes = strings.Split(createVolume, " ")
	}
	ctx := cmd.GetContext(rawArgs)
	return commands.RunCreate(ctx, args)
}
