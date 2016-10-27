// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	flag "github.com/docker/docker/pkg/mflag"

	version "github.com/hashicorp/go-version"
	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/clilogger"
	"github.com/scaleway/scaleway-cli/pkg/commands"
	"github.com/scaleway/scaleway-cli/pkg/config"
	"github.com/scaleway/scaleway-cli/pkg/scwversion"
	"github.com/scaleway/scaleway-cli/pkg/utils"
)

// global options
var (
	flDebug     = flag.Bool([]string{"D", "-debug"}, false, "Enable debug mode")
	flVerbose   = flag.Bool([]string{"V", "-verbose"}, false, "Enable verbose mode")
	flVersion   = flag.Bool([]string{"v", "-version"}, false, "Print version information and quit")
	flQuiet     = flag.Bool([]string{"q", "-quiet"}, false, "Enable quiet mode")
	flSensitive = flag.Bool([]string{"-sensitive"}, false, "Show sensitive data in outputs, i.e. API Token/Organization")
	flRegion    = flag.String([]string{"-region"}, "par1", "Change the default region (e.g. ams1)")
)

// Start is the entrypoint
func Start(rawArgs []string, streams *commands.Streams) (int, error) {
	checkVersion()
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

	if *flVersion {
		fmt.Fprintf(streams.Stderr, "scw version %s, build %s\n", scwversion.VERSION, scwversion.GITCOMMIT)
		return 0, nil
	}

	if *flSensitive {
		os.Setenv("SCW_SENSITIVE", "1")
	}

	if *flDebug {
		os.Setenv("DEBUG", "1")
	}

	if *flVerbose {
		os.Setenv("SCW_VERBOSE_API", "1")
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
				api, _ := getScalewayAPI(*flRegion)
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
				api, errGet := getScalewayAPI(*flRegion)
				if errGet != nil {
					return 1, fmt.Errorf("unable to initialize scw api: %v", errGet)
				}
				cmd.API = api
			}
			// clean cache between versions
			if cmd.API != nil && config.Version != scwversion.VERSION {
				cmd.API.ClearCache()
				config.Save()
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
func getScalewayAPI(region string) (*api.ScalewayAPI, error) {
	// We already get config globally, but whis way we can get explicit error when trying to create a ScalewayAPI object
	config, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	return api.NewScalewayAPI(config.Organization, config.Token, scwversion.UserAgent(), region, clilogger.SetupLogger)
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

func checkVersion() {
	if os.Getenv("SCW_NOCHECKVERSION") == "1" {
		return
	}
	homeDir, err := config.GetHomeDir()
	if err != nil {
		return
	}
	updateFiles := []string{"/var/run/.scw-update", "/tmp/.scw-update", filepath.Join(homeDir, ".scw-update")}
	updateFile := ""

	callAPI := false
	for _, file := range updateFiles {
		if stat, err := os.Stat(file); err == nil {
			updateFile = file
			callAPI = stat.ModTime().Before(time.Now().AddDate(0, 0, -1))
			break
		}
	}
	if updateFile == "" {
		for _, file := range updateFiles {
			if scwupdate, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0600); err == nil {
				scwupdate.Close()
				updateFile = file
				callAPI = true
				break
			}
		}
	}
	if callAPI {
		scwupdate, err := os.OpenFile(updateFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0600)
		if err != nil {
			return
		}
		scwupdate.Close()
		req := http.Client{
			Timeout: 1 * time.Second,
		}
		resp, err := req.Get("https://fr-1.storage.online.net/scaleway/scaleway-cli/VERSION")
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		if scwversion.VERSION == "" {
			return
		}
		ver := scwversion.VERSION
		if ver[0] == 'v' {
			ver = string([]byte(ver)[1:])
		}
		actual, err1 := version.NewVersion(ver)
		update, err2 := version.NewVersion(strings.Trim(string(body), "\n"))
		if err1 != nil || err2 != nil {
			return
		}
		if actual.LessThan(update) {
			logrus.Infof("A new version of scw is available (%v), beware that you are currently running %v", update, actual)
		}
	}
}
