package main

import (
	"fmt"
	"os"
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

func runCreate(cmd *Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: scw %s\n", cmd.UsageLine)
		os.Exit(1)
	}

	image := cmd.GetImage(args[0])
	var server ScalewayServerDefinition
	server.Name = createName

	// FIXME: handle snapshots
	// FIXME: handle creation of volumes
	// FIXME: handle tags
	server.Image = image

	if createBootscript != "" {
		bootscript := cmd.GetBootscript(createBootscript)
		server.Bootscript = &bootscript
	}

	serverId, err := cmd.API.PostServer(server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot create server: %v\n", err)
	}
	fmt.Println(serverId)
}
