// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"os"
	"runtime"

	"github.com/kardianos/osext"

	"github.com/scaleway/scaleway-cli/pkg/config"
	"github.com/scaleway/scaleway-cli/pkg/utils"
)

// InfoArgs are flags for the `RunInfo` function
type InfoArgs struct{}

// RunInfo is the handler for 'scw info'
func RunInfo(ctx CommandContext, args InfoArgs) error {
	// FIXME: fmt.Fprintf(ctx.Stdout, "Servers: %s\n", "quantity")
	// FIXME: fmt.Fprintf(ctx.Stdout, "Images: %s\n", "quantity")
	fmt.Fprintf(ctx.Stdout, "Debug mode (client):\t%v\n", ctx.Getenv("DEBUG") != "")

	fmt.Fprintf(ctx.Stdout, "Organization:\t\t%s\n", ctx.API.Organization)
	// FIXME: add partially-masked token
	fmt.Fprintf(ctx.Stdout, "API Endpoint:\t\t%s\n", ctx.Getenv("scaleway_api_endpoint"))
	configPath, _ := config.GetConfigFilePath()
	fmt.Fprintf(ctx.Stdout, "RC file:\t\t%s\n", configPath)
	fmt.Fprintf(ctx.Stdout, "User:\t\t\t%s\n", ctx.Getenv("USER"))
	fmt.Fprintf(ctx.Stdout, "CPUs:\t\t\t%d\n", runtime.NumCPU())
	hostname, _ := os.Hostname()
	fmt.Fprintf(ctx.Stdout, "Hostname:\t\t%s\n", hostname)
	cliPath, _ := osext.Executable()
	fmt.Fprintf(ctx.Stdout, "CLI Path:\t\t%s\n", cliPath)

	fmt.Fprintln(ctx.Stdout, "")
	fmt.Fprintf(ctx.Stdout, "Cache:\t\t\t%s\n", ctx.API.Cache.Path)
	fmt.Fprintf(ctx.Stdout, "  Servers:\t\t%d\n", ctx.API.Cache.GetNbServers())
	fmt.Fprintf(ctx.Stdout, "  Images:\t\t%d\n", ctx.API.Cache.GetNbImages())
	fmt.Fprintf(ctx.Stdout, "  Snapshots:\t\t%d\n", ctx.API.Cache.GetNbSnapshots())
	fmt.Fprintf(ctx.Stdout, "  Volumes:\t\t%d\n", ctx.API.Cache.GetNbVolumes())
	fmt.Fprintf(ctx.Stdout, "  Bootscripts:\t\t%d\n", ctx.API.Cache.GetNbBootscripts())

	user, err := ctx.API.GetUser()
	if err != nil {
		return fmt.Errorf("Unable to get your SSH Keys")
	}

	if len(user.SSHPublicKeys) == 0 {
		fmt.Fprintln(ctx.Stdout, "You have no ssh keys")
	} else {
		fmt.Fprintln(ctx.Stdout, "")
		fmt.Fprintln(ctx.Stdout, "SSH Keys:")
		for id, key := range user.SSHPublicKeys {
			fingerprint, err := utils.SSHGetFingerprint([]byte(key.Key))
			if err != nil {
				return err
			}
			fmt.Fprintf(ctx.Stdout, "  [%d] %s\n", id, fingerprint)
		}
		fmt.Fprintf(ctx.Stdout, "\n")
	}

	dashboard, err := ctx.API.GetDashboard()
	if err != nil {
		return fmt.Errorf("Unable to get your dashboard")
	}
	fmt.Fprintln(ctx.Stdout, "Dashboard:")
	fmt.Fprintf(ctx.Stdout, "  Volumes:\t\t%d\n", dashboard.VolumesCount)
	fmt.Fprintf(ctx.Stdout, "  Running servers:\t%d\n", dashboard.RunningServersCount)
	fmt.Fprintf(ctx.Stdout, "  Images:\t\t%d\n", dashboard.ImagesCount)
	fmt.Fprintf(ctx.Stdout, "  Snapshots:\t\t%d\n", dashboard.SnapshotsCount)
	fmt.Fprintf(ctx.Stdout, "  Servers:\t\t%d\n", dashboard.ServersCount)
	fmt.Fprintf(ctx.Stdout, "  Ips:\t\t\t%d\n", dashboard.IPsCount)

	fmt.Fprintf(ctx.Stdout, "\n")
	permissions, err := ctx.API.GetPermissions()
	if err != nil {
		return fmt.Errorf("Unable to get your permisssions")
	}
	fmt.Fprintln(ctx.Stdout, "Permissions:")
	for _, service := range permissions.Permissions {
		for key, serviceName := range service {
			fmt.Fprintf(ctx.Stdout, "  %s\n", key)
			for _, perm := range serviceName {
				fmt.Fprintf(ctx.Stdout, "    %s\n", perm)
			}
		}
	}
	fmt.Fprintf(ctx.Stdout, "\n")
	quotas, err := ctx.API.GetQuotas()
	if err != nil {
		return fmt.Errorf("Unable to get your quotas")
	}
	fmt.Fprintln(ctx.Stdout, "Quotas:")
	for key, value := range quotas.Quotas {
		fmt.Fprintf(ctx.Stdout, "  %-20s: %d\n", key, value)
	}
	return nil
}
