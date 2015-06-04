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
	cmdCreate.Flag.StringVar(&createName, []string{"-name"}, "", "Assign a name")
	cmdCreate.Flag.StringVar(&createBootscript, []string{"-bootscript"}, "", "Assign a bootscript")
	cmdCreate.Flag.StringVar(&createEnv, []string{"e", "-env"}, "", "Provide metadata tags passed to initrd (i.e., boot=resue INITRD_DEBUG=1)")
	cmdCreate.Flag.StringVar(&createVolume, []string{"v", "-volume"}, "", "Attach additional volume (i.e., 50G)")
	cmdCreate.Flag.BoolVar(&createHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var createName string       // --name flag
var createBootscript string // --bootscript flag
var createEnv string        // -e, --env flag
var createVolume string     // -v, --volume flag
var createHelp bool         // -h, --help flag

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

func createServer(cmd *Command, imageName string, name string, bootscript string, env string, additionalVolumes string) (string, error) {
	if name == "" {
		name = strings.Replace(namesgenerator.GetRandomName(0), "_", "-", -1)
	}

	var server ScalewayServerDefinition
	server.Volumes = make(map[string]string)

	server.Tags = []string{}
	if env != "" {
		server.Tags = strings.Split(env, " ")
	}
	if additionalVolumes != "" {
		volumes := strings.Split(additionalVolumes, " ")
		for i := range volumes {
			volumeId, err := CreateVolumeFromHumanSize(cmd, volumes[i])
			if err != nil {
				return "", err
			}

			volumeIdx := fmt.Sprintf("%d", i+1)
			server.Volumes[volumeIdx] = *volumeId
		}
	}
	server.Name = name
	if bootscript != "" {
		bootscript := cmd.GetBootscript(bootscript)
		server.Bootscript = &bootscript
	}

	_, err := humanize.ParseBytes(imageName)
	if err == nil {
		// Create a new root volume
		volumeId, err := CreateVolumeFromHumanSize(cmd, imageName)
		if err != nil {
			return "", err
		}
		server.Volumes["0"] = *volumeId
	} else {
		// Use an existing image
		// FIXME: handle snapshots
		image := cmd.GetImage(imageName)
		server.Image = &image
	}

	serverId, err := cmd.API.PostServer(server)
	if err != nil {
		return "", nil
	}

	return serverId, nil

}

func runCreate(cmd *Command, args []string) {
	if createHelp {
		cmd.PrintUsage()
	}
	if len(args) != 1 {
		cmd.PrintShortUsage()
	}

	serverId, err := createServer(cmd, args[0], createName, createBootscript, createEnv, createVolume)

	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	fmt.Println(serverId)
}
