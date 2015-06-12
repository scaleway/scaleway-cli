package main

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	flag "github.com/docker/docker/pkg/mflag"

	cmds "github.com/scaleway/scaleway-cli/commands"
)

func main() {
	config, cfgErr := getConfig()
	if cfgErr != nil && !os.IsNotExist(cfgErr) {
		log.Fatalf("Unable to open .scwrc config file: %v", cfgErr)
	}

	if config != nil {
		flAPIEndPoint = flag.String([]string{"-api-endpoint"}, config.APIEndPoint, "Set the API endpoint")
	}
	flag.Parse()

	if *flVersion {
		showVersion()
		return
	}

	if flAPIEndPoint != nil {
		os.Setenv("scaleway_api_endpoint", *flAPIEndPoint)
	}

	if *flDebug {
		os.Setenv("DEBUG", "1")
	}

	initLogging(os.Getenv("DEBUG") != "")

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}
	name := args[0]

	args = args[1:]

	for _, cmd := range cmds.Commands {
		if cmd.Name() == name {
			cmd.Flag.SetOutput(ioutil.Discard)
			err := cmd.Flag.Parse(args)
			if err != nil {
				log.Fatalf("usage: scw %s", cmd.UsageLine)
			}
			if cmd.Name() != "login" && cmd.Name() != "help" {
				if cfgErr != nil {
					if name != "login" && config == nil {
						fmt.Fprintf(os.Stderr, "You need to login first: 'scw login'\n")
						os.Exit(1)
					}
				}
				api, err := getScalewayAPI()
				if err != nil {
					log.Fatalf("unable to initialize scw api: %s", err)
				}
				cmd.API = api
			}
			cmd.Exec(cmd, cmd.Flag.Args())
			if cmd.API != nil {
				cmd.API.Sync()
			}
			os.Exit(0)
		}
	}

	log.Fatalf("scw: unknown subcommand %s\nRun 'scw help' for usage.", name)
}

func initLogging(debug bool) {
	log.SetOutput(os.Stderr)
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}
