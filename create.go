package main

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/namesgenerator"
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
	cmdCreate.Flag.StringVar(&createName, []string{"-name"}, "", "Assign a name")
	cmdCreate.Flag.StringVar(&createBootscript, []string{"-bootscript"}, "", "Assign a bootscript")
	cmdCreate.Flag.StringVar(&createEnv, []string{"e", "-env"}, "", "Provide metadata tags passed to initrd (i.e., boot=resue INITRD_DEBUG=1)")
	cmdCreate.Flag.StringVar(&createVolume, []string{"v", "-volume"}, "", "Attach additional volume (i.e., 50G)")
}

// Flags
var createName string
var createBootscript string
var createEnv string
var createVolume string

func CreateVolumeFromHumanSize(cmd *Command, size string) (*string, error) {
	bytes, err := humanize.ParseBytes(size)
	if err != nil {
		return nil, err
	}

	var newVolume ScalewayVolumeDefinition
	newVolume.Name = size
	newVolume.Size = bytes
	newVolume.Type = "l_ssd"

	volumeId, err := cmd.API.PostVolume(newVolume)
	if err != nil {
		return nil, err
	}

	return &volumeId, nil
}

func runCreate(cmd *Command, args []string) {
	if len(args) != 1 {
		log.Fatalf("usage: scw %s", cmd.UsageLine)
	}

	if createName == "" {
		createName = strings.Replace(namesgenerator.GetRandomName(0), "_", "-", -1)
	}

	var server ScalewayServerDefinition
	server.Volumes = make(map[string]string)

	server.Tags = []string{}
	if createEnv != "" {
		server.Tags = strings.Split(createEnv, " ")
	}
	if createVolume != "" {
		volumes := strings.Split(createVolume, " ")
		for i := range volumes {
			volumeId, err := CreateVolumeFromHumanSize(cmd, volumes[i])
			if err != nil {
				log.Fatalf("Failed to create volume: %v", err)
			}

			volumeIdx := fmt.Sprintf("%d", i+1)
			server.Volumes[volumeIdx] = *volumeId
		}
	}
	server.Name = createName
	if createBootscript != "" {
		bootscript := cmd.GetBootscript(createBootscript)
		server.Bootscript = &bootscript
	}

	_, err := humanize.ParseBytes(args[0])
	if err == nil {
		// Create a new root volume
		volumeId, err := CreateVolumeFromHumanSize(cmd, args[0])
		if err != nil {
			log.Fatalf("Failed to create volume: %v", err)
		}
		server.Volumes["0"] = *volumeId
	} else {
		// Use an existing image
		// FIXME: handle snapshots
		image := cmd.GetImage(args[0])
		server.Image = &image
	}

	serverId, err := cmd.API.PostServer(server)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	fmt.Println(serverId)
}
