// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

// Manage BareMetal Servers from Command Line (as easily as with Docker)
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"
	flag "github.com/scaleway/scaleway-cli/vendor/github.com/docker/docker/pkg/mflag"

	"github.com/scaleway/scaleway-cli/api"
	cmds "github.com/scaleway/scaleway-cli/commands"
	"github.com/scaleway/scaleway-cli/scwversion"
	"github.com/scaleway/scaleway-cli/utils"
)

// CommandListOpts holds a list of parameters
type CommandListOpts struct {
	Values *[]string
}

// NewListOpts create an empty CommandListOpts
func NewListOpts() CommandListOpts {
	var values []string
	return CommandListOpts{
		Values: &values,
	}
}

// String returns a string representation of a CommandListOpts object
func (opts *CommandListOpts) String() string {
	return fmt.Sprintf("%v", []string((*opts.Values)))
}

// Set appends a new value to a CommandListOpts
func (opts *CommandListOpts) Set(value string) error {
	(*opts.Values) = append((*opts.Values), value)
	return nil
}

func commandUsage(name string) {
}

var (
	flAPIEndPoint *string
	flDebug       = flag.Bool([]string{"D", "-debug"}, false, "Enable debug mode")
	flVerbose     = flag.Bool([]string{"V", "-verbose"}, false, "Enable verbose mode")
	flVersion     = flag.Bool([]string{"v", "-version"}, false, "Print version information and quit")
	flSensitive   = flag.Bool([]string{"-sensitive"}, false, "Show sensitive data in outputs, i.e. API Token/Organization")
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

	if *flSensitive {
		os.Setenv("SCW_SENSITIVE", "1")
	}

	if *flDebug {
		os.Setenv("DEBUG", "1")
	}

	initLogging(os.Getenv("DEBUG") != "", *flVerbose)

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
			if cmd.Name() != "login" && cmd.Name() != "help" && cmd.Name() != "version" {
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

func usage() {
	cmds.CmdHelp.Exec(cmds.CmdHelp, []string{})
	os.Exit(1)
}

// getConfig returns the Scaleway CLI config file for the current user
func getConfig() (*api.Config, error) {
	scwrcPath, err := utils.GetConfigFilePath()
	if err != nil {
		return nil, err
	}

	stat, err := os.Stat(scwrcPath)
	// we don't care if it fails, the user just won't see the warning
	if err == nil {
		mode := stat.Mode()
		if mode&0066 != 0 {
			log.Fatalf("Permissions %#o for .scwrc are too open.", mode)
		}
	}

	file, err := ioutil.ReadFile(scwrcPath)
	if err != nil {
		return nil, err
	}
	var config api.Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	if os.Getenv("scaleway_api_endpoint") == "" {
		os.Setenv("scaleway_api_endpoint", config.APIEndPoint)
	}
	return &config, nil
}

// getScalewayAPI returns a ScalewayAPI using the user config file
func getScalewayAPI() (*api.ScalewayAPI, error) {
	// We already get config globally, but whis way we can get explicit error when trying to create a ScalewayAPI object
	config, err := getConfig()
	if err != nil {
		return nil, err
	}
	return api.NewScalewayAPI(os.Getenv("scaleway_api_endpoint"), config.Organization, config.Token)
}

func showVersion() {
	fmt.Printf("scw version %s, build %s\n", scwversion.VERSION, scwversion.GITCOMMIT)
}

func initLogging(debug bool, verbose bool) {
	log.SetOutput(os.Stderr)
	if debug {
		log.SetLevel(log.DebugLevel)
	} else if verbose {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}
}
