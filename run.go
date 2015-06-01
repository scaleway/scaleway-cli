package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

var cmdRun = &Command{
	Exec:        runRun,
	UsageLine:   "run [OPTIONS] IMAGE [COMMAND] [ARG...]",
	Description: "Run a command in a new server",
	Help:        "Run a command in a new server..",
}

func init() {
	// FIXME: -h
	cmdRun.Flag.StringVar(&runCreateName, []string{"-name"}, "", "Assign a name")
	cmdRun.Flag.StringVar(&runCreateBootscript, []string{"-bootscript"}, "", "Assign a bootscript")
	cmdRun.Flag.StringVar(&runCreateEnv, []string{"e", "-env"}, "", "Provide metadata tags passed to initrd (i.e., boot=resue INITRD_DEBUG=1)")
	cmdRun.Flag.StringVar(&runCreateVolume, []string{"v", "-volume"}, "", "Attach additional volume (i.e., 50G)")
}

// Flags
var runCreateName string
var runCreateBootscript string
var runCreateEnv string
var runCreateVolume string

// FIXME: handle start --timeout

func runRun(cmd *Command, args []string) {
	if len(args) < 1 {
		log.Fatalf("usage: scw %s", cmd.UsageLine)
	}

	//image := args[0]

	// create IMAGE
	log.Debugf("Creating a new server")
	serverId, err := createServer(cmd, args[0], runCreateName, runCreateBootscript, runCreateEnv, runCreateVolume)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	log.Debugf("Created server: %s", serverId)

	// start SERVER
	log.Debugf("Starting server")
	err = startServer(cmd, serverId, false)
	if err != nil {
		log.Fatalf("Failed to start server %s: %v", serverId, err)
	}
	log.Debugf("Server is booting")

	// waiting for server to be ready
	log.Debugf("Waiting for server to be ready")
	server, err := WaitForServerReady(cmd.API, serverId)
	if err != nil {
		log.Fatalf("Cannot get access to server %s: %v", serverId, err)
	}
	log.Debugf("Server is ready: %s", server.PublicAddress.IP)

	// exec -w SERVER COMMAND ARGS...
	log.Debugf("Executing command")
	if len(args) < 2 {
		err = serverExec(server.PublicAddress.IP, []string{"if [ -x /bin/bash ]; then /bin/bash; else /bin/sh; fi"})
	} else {
		err = serverExec(server.PublicAddress.IP, args[1:])
	}
	if err != nil {
		log.Debugf("Command execution failed: %v", err)
		os.Exit(1)
	}
	log.Debugf("Command successfuly executed")
}
