// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"os"
	"runtime"

	"github.com/scaleway/scaleway-cli/vendor/github.com/kardianos/osext"

	"github.com/scaleway/scaleway-cli/pkg/config"
)

// InfoArgs are flags for the `RunInfo` function
type InfoArgs struct{}

// RunInfo is the handler for 'scw info'
func RunInfo(ctx CommandContext, args InfoArgs) error {
	// FIXME: fmt.Fprintf(ctx.Stdout, "Servers: %s\n", "quantity")
	// FIXME: fmt.Fprintf(ctx.Stdout, "Images: %s\n", "quantity")
	fmt.Fprintf(ctx.Stdout, "Debug mode (client): %v\n", ctx.Getenv("DEBUG") != "")

	fmt.Fprintf(ctx.Stdout, "Organization: %s\n", ctx.API.Organization)
	// FIXME: add partially-masked token
	fmt.Fprintf(ctx.Stdout, "API Endpoint: %s\n", ctx.Getenv("scaleway_api_endpoint"))
	configPath, _ := config.GetConfigFilePath()
	fmt.Fprintf(ctx.Stdout, "RC file: %s\n", configPath)
	fmt.Fprintf(ctx.Stdout, "User: %s\n", ctx.Getenv("USER"))
	fmt.Fprintf(ctx.Stdout, "CPUs: %d\n", runtime.NumCPU())
	hostname, _ := os.Hostname()
	fmt.Fprintf(ctx.Stdout, "Hostname: %s\n", hostname)
	cliPath, _ := osext.Executable()
	fmt.Fprintf(ctx.Stdout, "CLI Path: %s\n", cliPath)

	fmt.Fprintf(ctx.Stdout, "Cache: %s\n", ctx.API.Cache.Path)
	fmt.Fprintf(ctx.Stdout, "  Servers: %d\n", ctx.API.Cache.GetNbServers())
	fmt.Fprintf(ctx.Stdout, "  Images: %d\n", ctx.API.Cache.GetNbImages())
	fmt.Fprintf(ctx.Stdout, "  Snapshots: %d\n", ctx.API.Cache.GetNbSnapshots())
	fmt.Fprintf(ctx.Stdout, "  Volumes: %d\n", ctx.API.Cache.GetNbVolumes())
	fmt.Fprintf(ctx.Stdout, "  Bootscripts: %d\n", ctx.API.Cache.GetNbBootscripts())
	return nil
}
