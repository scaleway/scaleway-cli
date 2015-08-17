// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"strings"

	"github.com/scaleway/scaleway-cli/pkg/commands"
	"github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"
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
	cmdCreate.Flag.StringVar(&createEnv, []string{"e", "-env"}, "", "Provide metadata tags passed to initrd (i.e., boot=resue INITRD_DEBUG=1)")
	cmdCreate.Flag.StringVar(&createVolume, []string{"v", "-volume"}, "", "Attach additional volume (i.e., 50G)")
	cmdCreate.Flag.BoolVar(&createHelp, []string{"h", "-help"}, false, "Print usage")
	cmdCreate.Flag.BoolVar(&createTmpSSHKey, []string{"-tmp-ssh-key"}, false, "Access your server without uploading your SSH key to your account")
}

// Flags
var createName string       // --name flag
var createBootscript string // --bootscript flag
var createEnv string        // -e, --env flag
var createVolume string     // -v, --volume flag
var createHelp bool         // -h, --help flag
var createTmpSSHKey bool    // --tmp-ssh-key flag

func runCreate(cmd *Command, rawArgs []string) {
	if createHelp {
		cmd.PrintUsage()
	}
	if len(rawArgs) != 1 {
		cmd.PrintShortUsage()
	}

	args := commands.CreateArgs{
		Name:       createName,
		Bootscript: createBootscript,
		Image:      rawArgs[0],
		TmpSSHKey:  createTmpSSHKey,
	}

	if len(runCreateEnv) > 0 {
		args.Tags = strings.Split(runCreateEnv, " ")
	}
	if len(runCreateVolume) > 0 {
		args.Volumes = strings.Split(runCreateVolume, " ")
	}
	ctx := cmd.GetContext(rawArgs)
	err := commands.RunCreate(ctx, args)
	if err != nil {
		logrus.Fatalf("Cannot execute 'create': %v", err)
	}
}
