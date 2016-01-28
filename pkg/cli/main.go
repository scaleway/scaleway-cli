// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
	flag "github.com/docker/docker/pkg/mflag"

	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/commands"
	"github.com/scaleway/scaleway-cli/pkg/config"
	"github.com/scaleway/scaleway-cli/pkg/scwversion"
	"github.com/scaleway/scaleway-cli/pkg/utils"
)

// global options
var (
	flAPIEndPoint *string
	flDebug       = flag.Bool([]string{"D", "-debug"}, false, "Enable debug mode")
	flVerbose     = flag.Bool([]string{"V", "-verbose"}, false, "Enable verbose mode")
	flVersion     = flag.Bool([]string{"v", "-version"}, false, "Print version information and quit")
	flQuiet       = flag.Bool([]string{"q", "-quiet"}, false, "Enable quiet mode")
	flSensitive   = flag.Bool([]string{"-sensitive"}, false, "Show sensitive data in outputs, i.e. API Token/Organization")
)

// Start is the entrypoint
func Start(rawArgs []string, streams *commands.Streams) (int, error) {
	if streams == nil {
		streams = &commands.Streams{
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}
	}

	flag.CommandLine.Parse(rawArgs)

	config, cfgErr := config.GetConfig()
	if cfgErr != nil && !os.IsNotExist(cfgErr) {
		return 1, fmt.Errorf("unable to open .scwrc config file: %v", cfgErr)
	}

	if config != nil {
		defaultComputeAPI := os.Getenv("scaleway_api_endpoint")
		if defaultComputeAPI == "" {
			defaultComputeAPI = config.ComputeAPI
		}
		if flAPIEndPoint == nil {
			flAPIEndPoint = flag.String([]string{"-api-endpoint"}, defaultComputeAPI, "Set the API endpoint")
		}
	}

	if *flVersion {
		fmt.Fprintf(streams.Stderr, "scw version %s, build %s\n", scwversion.VERSION, scwversion.GITCOMMIT)
		return 0, nil
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

	utils.Quiet(*flQuiet)
	initLogging(os.Getenv("DEBUG") != "", *flVerbose, streams)

	args := flag.Args()
	if len(args) < 1 {
		CmdHelp.Exec(CmdHelp, []string{})
		return 1, nil
	}
	name := args[0]

	args = args[1:]

	// Apply default values
	for _, cmd := range Commands {
		cmd.streams = streams
	}

	for _, cmd := range Commands {
		if cmd.Name() == name {
			cmd.Flag.SetOutput(ioutil.Discard)
			err := cmd.Flag.Parse(args)
			if err != nil {
				return 1, fmt.Errorf("usage: scw %s", cmd.UsageLine)
			}
			switch cmd.Name() {
			case "login", "help", "version":
				// commands that don't need API
			case "_userdata":
				// commands that may need API
				api, _ := getScalewayAPI()
				cmd.API = api
			default:
				// commands that do need API
				if cfgErr != nil {
					if name != "login" && config == nil {
						logrus.Debugf("cfgErr: %v", cfgErr)
						fmt.Fprintf(streams.Stderr, "You need to login first: 'scw login'\n")
						return 1, nil
					}
				}
				api, err := getScalewayAPI()
				if err != nil {
					return 1, fmt.Errorf("unable to initialize scw api: %s", err)
				}
				cmd.API = api
			}
			err = cmd.Exec(cmd, cmd.Flag.Args())
			switch err {
			case nil:
			case ErrExitFailure:
				return 1, nil
			case ErrExitSuccess:
				return 0, nil
			default:
				return 1, fmt.Errorf("cannot execute '%s': %v", cmd.Name(), err)
			}
			if cmd.API != nil {
				cmd.API.Sync()
			}
			return 0, nil
		}
	}

	return 1, fmt.Errorf("scw: unknown subcommand %s\nRun 'scw help' for usage", name)
}

// getScalewayAPI returns a ScalewayAPI using the user config file
func getScalewayAPI() (*api.ScalewayAPI, error) {
	// We already get config globally, but whis way we can get explicit error when trying to create a ScalewayAPI object
	config, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	return api.NewScalewayAPI(os.Getenv("scaleway_api_endpoint"), config.AccountAPI, config.Organization, config.Token, scwversion.UserAgent())
}

func initLogging(debug bool, verbose bool, streams *commands.Streams) {
	logrus.SetOutput(streams.Stderr)
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else if verbose {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}
}
