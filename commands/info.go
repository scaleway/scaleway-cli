// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"os"
	"runtime"

	"github.com/scaleway/scaleway-cli/vendor/github.com/kardianos/osext"

	types "github.com/scaleway/scaleway-cli/commands/types"
	"github.com/scaleway/scaleway-cli/utils"
)

var cmdInfo = &types.Command{
	Exec:        runInfo,
	UsageLine:   "info [OPTIONS]",
	Description: "Display system-wide information",
	Help:        "Display system-wide information.",
}

func init() {
	cmdInfo.Flag.BoolVar(&infoHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var infoHelp bool // -h, --help flag

func runInfo(cmd *types.Command, args []string) {
	if infoHelp {
		cmd.PrintUsage()
	}
	if len(args) != 0 {
		cmd.PrintShortUsage()
	}

	// FIXME: fmt.Printf("Servers: %s\n", "quantity")
	// FIXME: fmt.Printf("Images: %s\n", "quantity")
	fmt.Printf("Debug mode (client): %v\n", os.Getenv("DEBUG") != "")

	fmt.Printf("Organization: %s\n", cmd.API.Organization)
	// FIXME: add partially-masked token
	fmt.Printf("API Endpoint: %s\n", os.Getenv("scaleway_api_endpoint"))
	configPath, _ := utils.GetConfigFilePath()
	fmt.Printf("RC file: %s\n", configPath)
	fmt.Printf("User: %s\n", os.Getenv("USER"))
	fmt.Printf("CPUs: %d\n", runtime.NumCPU())
	hostname, _ := os.Hostname()
	fmt.Printf("Hostname: %s\n", hostname)
	cliPath, _ := osext.Executable()
	fmt.Printf("CLI Path: %s\n", cliPath)

	fmt.Printf("Cache: %s\n", cmd.API.Cache.Path)
	fmt.Printf("  Servers: %d\n", cmd.API.Cache.GetNbServers())
	fmt.Printf("  Images: %d\n", cmd.API.Cache.GetNbImages())
	fmt.Printf("  Snapshots: %d\n", cmd.API.Cache.GetNbSnapshots())
	fmt.Printf("  Volumes: %d\n", cmd.API.Cache.GetNbVolumes())
	fmt.Printf("  Bootscripts: %d\n", cmd.API.Cache.GetNbBootscripts())
}
