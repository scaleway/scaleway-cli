package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	humanize "github.com/dustin/go-humanize"
)

var cmdCreate = &Command{
	Exec:        runCreate,
	UsageLine:   "create [OPTIONS] IMAGE",
	Description: "Create a new server but do not create it",
	Help:        "Create a new server but do not create it.",
}

func init() {
	// FIXME: -h
	cmdCreate.Flag.StringVar(&createName, []string{"-name"}, "noname", "Assign a name")
	cmdCreate.Flag.StringVar(&createBootscript, []string{"-bootscript"}, "", "Assign a bootscript")
}

// Flags
var createName string
var createBootscript string

func CreateServerCommonFields(cmd *Command, server interface{}) {
}

func runCreate(cmd *Command, args []string) {
	if len(args) == 0 {
		log.Fatalf("usage: scw %s", cmd.UsageLine)
	}

	var serverId string

	// FIXME: use an interface to remove duplicates

	bytes, err := humanize.ParseBytes(args[0])
	if err == nil {
		var server ScalewayServerWithVolumeDefinition
		// Create a new root volume
		var newVolume ScalewayVolumeDefinition
		newVolume.Name = args[0]
		newVolume.Size = bytes
		newVolume.Type = "l_ssd"
		volumeId, err := cmd.API.PostVolume(newVolume)
		if err != nil {
			log.Fatalf("Failed to create volume: %v", err)
		}
		server.Volumes = make(map[string]string)
		server.Volumes["0"] = volumeId

		// Common fields
		server.Organization = cmd.API.Organization
		server.Name = createName
		if createBootscript != "" {
			bootscript := cmd.GetBootscript(createBootscript)
			server.Bootscript = &bootscript
		}
		// FIXME: handle tags
		// End of common fields

		serverId, err = cmd.API.PostServer(server)
		if err != nil {
			log.Fatalf("Failed to create server: %v", err)
		}
	} else {
		var server ScalewayServerWithImageDefinition
		// Use an existing image
		// FIXME: handle snapshots
		image := cmd.GetImage(args[0])
		server.Image = image

		// Common fields
		server.Organization = cmd.API.Organization
		server.Name = createName
		if createBootscript != "" {
			bootscript := cmd.GetBootscript(createBootscript)
			server.Bootscript = &bootscript
		}
		// FIXME: handle tags
		// End of common fields

		serverId, err = cmd.API.PostServer(server)
		if err != nil {
			log.Fatalf("Failed to create server: %v", err)
		}
	}

	fmt.Println(serverId)
}
