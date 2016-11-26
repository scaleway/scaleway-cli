// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import "github.com/scaleway/scaleway-cli/pkg/commands"

var cmdCp = &Command{
	Exec:        runCp,
	UsageLine:   "cp [OPTIONS] SERVER:PATH|HOSTPATH|- SERVER:PATH|HOSTPATH|-",
	Description: "Copy files/folders from a PATH on the server to a HOSTDIR on the host",
	Help:        "Copy files/folders from a PATH on the server to a HOSTDIR on the host\nrunning the command. Use '-' to write the data as a tar file to STDOUT.",
	Examples: `
    $ scw cp path/to/my/local/file myserver:path
    $ scw cp --gateway=myotherserver path/to/my/local/file myserver:path
    $ scw cp myserver:path/to/file path/to/my/local/dir
    $ scw cp myserver:path/to/file myserver2:path/to/dir
    $ scw cp myserver:path/to/file - > myserver-pathtofile-backup.tar
    $ scw cp myserver:path/to/file - | tar -tvf -
    $ scw cp path/to/my/local/dir  myserver:path
    $ scw cp myserver:path/to/dir  path/to/my/local/dir
    $ scw cp myserver:path/to/dir  myserver2:path/to/dir
    $ scw cp myserver:path/to/dir  - > myserver-pathtodir-backup.tar
    $ scw cp myserver:path/to/dir  - | tar -tvf -
    $ cat archive.tar | scw cp - myserver:/path
    $ tar -cvf - . | scw cp - myserver:path
`,
}

func init() {
	cmdCp.Flag.BoolVar(&cpHelp, []string{"h", "-help"}, false, "Print usage")
	cmdCp.Flag.StringVar(&cpGateway, []string{"g", "-gateway"}, "", "Use a SSH gateway")
	cmdCp.Flag.StringVar(&cpSSHUser, []string{"-user"}, "root", "Specify SSH user")
	cmdCp.Flag.IntVar(&cpSSHPort, []string{"p", "-port"}, 22, "Specify SSH port")
}

// Flags
var cpHelp bool      // -h, --help flag
var cpGateway string // -g, --gateway flag
var cpSSHUser string // --user flag
var cpSSHPort int    // -p, --port flag

func runCp(cmd *Command, rawArgs []string) error {
	if cpHelp {
		return cmd.PrintUsage()
	}
	if len(rawArgs) != 2 {
		return cmd.PrintShortUsage()
	}

	args := commands.CpArgs{
		Gateway:     cpGateway,
		Source:      rawArgs[0],
		Destination: rawArgs[1],
		SSHUser:     cpSSHUser,
		SSHPort:     cpSSHPort,
	}
	ctx := cmd.GetContext(rawArgs)
	return commands.RunCp(ctx, args)
}
