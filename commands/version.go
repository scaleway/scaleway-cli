// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"runtime"

	types "github.com/scaleway/scaleway-cli/commands/types"
	"github.com/scaleway/scaleway-cli/scwversion"
)

var cmdVersion = &types.Command{
	Exec:        runVersion,
	UsageLine:   "version [OPTIONS]",
	Description: "Show the version information",
	Help:        "Show the version information.",
}

func init() {
	cmdVersion.Flag.BoolVar(&versionHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var versionHelp bool // -h, --help flag

func runVersion(cmd *types.Command, args []string) {
	if versionHelp {
		cmd.PrintUsage()
	}
	if len(args) != 0 {
		cmd.PrintShortUsage()
	}

	fmt.Printf("Client version: %s\n", scwversion.VERSION)
	fmt.Printf("Go version (client): %s\n", runtime.Version())
	fmt.Printf("Git commit (client): %s\n", scwversion.GITCOMMIT)
	fmt.Printf("OS/Arch (client): %s/%s\n", runtime.GOOS, runtime.GOARCH)
	// FIXME: API version information
}
