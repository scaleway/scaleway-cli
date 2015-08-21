// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"fmt"
	"strings"

	"github.com/scaleway/scaleway-cli/pkg/commands"
)

var cmdRun = &Command{
	Exec:        runRun,
	UsageLine:   "run [OPTIONS] IMAGE [COMMAND] [ARG...]",
	Description: "Run a command in a new server",
	Help:        "Run a command in a new server.",
	Examples: `
    $ scw run ubuntu-trusty
    $ scw run --rm ubuntu-trusty
    $ scw run -a --rm ubuntu-trusty
    $ scw run --gateway=myotherserver ubuntu-trusty
    $ scw run ubuntu-trusty bash
    $ scw run --name=mydocker docker docker run moul/nyancat:armhf
    $ scw run --bootscript=3.2.34 --env="boot=live rescue_image=http://j.mp/scaleway-ubuntu-trusty-tarball" 50GB bash
    $ scw run --attach alpine
    $ scw run --detach alpine
    $ scw run --tmp-ssh-key alpine
`,
}

func init() {
	cmdRun.Flag.StringVar(&runCreateName, []string{"-name"}, "", "Assign a name")
	cmdRun.Flag.StringVar(&runCreateBootscript, []string{"-bootscript"}, "", "Assign a bootscript")
	cmdRun.Flag.StringVar(&runCreateEnv, []string{"e", "-env"}, "", "Provide metadata tags passed to initrd (i.e., boot=resue INITRD_DEBUG=1)")
	cmdRun.Flag.StringVar(&runCreateVolume, []string{"v", "-volume"}, "", "Attach additional volume (i.e., 50G)")
	cmdRun.Flag.BoolVar(&runHelpFlag, []string{"h", "-help"}, false, "Print usage")
	cmdRun.Flag.BoolVar(&runAttachFlag, []string{"a", "-attach"}, false, "Attach to serial console")
	cmdRun.Flag.BoolVar(&runDetachFlag, []string{"d", "-detach"}, false, "Run server in background and print server ID")
	cmdRun.Flag.StringVar(&runGateway, []string{"g", "-gateway"}, "", "Use a SSH gateway")
	cmdRun.Flag.BoolVar(&runAutoRemove, []string{"-rm"}, false, "Automatically remove the server when it exits")
	cmdRun.Flag.BoolVar(&runTmpSSHKey, []string{"-tmp-ssh-key"}, false, "Access your server without uploading your SSH key to your account")
	// FIXME: handle start --timeout
}

// Flags
var runCreateName string       // --name flag
var runAutoRemove bool         // --rm flag
var runCreateBootscript string // --bootscript flag
var runCreateEnv string        // -e, --env flag
var runCreateVolume string     // -v, --volume flag
var runHelpFlag bool           // -h, --help flag
var runAttachFlag bool         // -a, --attach flag
var runDetachFlag bool         // -d, --detach flag
var runGateway string          // -g, --gateway flag
var runTmpSSHKey bool          // --tmp-ssh-key flag

func runRun(cmd *Command, rawArgs []string) error {
	if runHelpFlag {
		return cmd.PrintUsage()
	}
	if len(rawArgs) < 1 {
		return cmd.PrintShortUsage()
	}
	if runAttachFlag && len(rawArgs) > 1 {
		return fmt.Errorf("conflicting options: -a and COMMAND")
	}
	if runAttachFlag && runDetachFlag {
		return fmt.Errorf("conflicting options: -a and -d")
	}
	if runDetachFlag && len(rawArgs) > 1 {
		return fmt.Errorf("conflicting options: -d and COMMAND")
	}
	if runAutoRemove && runDetachFlag {
		return fmt.Errorf("conflicting options: --detach and --rm")
	}

	args := commands.RunArgs{
		Attach:     runAttachFlag,
		Bootscript: runCreateBootscript,
		Command:    rawArgs[1:],
		Detach:     runDetachFlag,
		Gateway:    runGateway,
		Image:      rawArgs[0],
		Name:       runCreateName,
		AutoRemove: runAutoRemove,
		TmpSSHKey:  runTmpSSHKey,
		// FIXME: DynamicIPRequired
		// FIXME: Timeout
	}

	if len(runCreateEnv) > 0 {
		args.Tags = strings.Split(runCreateEnv, " ")
	}
	if len(runCreateVolume) > 0 {
		args.Volumes = strings.Split(runCreateVolume, " ")
	}
	ctx := cmd.GetContext(rawArgs)
	return commands.Run(ctx, args)
}
