// scw interract with Scaleway from the command line
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	flag "github.com/docker/docker/pkg/mflag"

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
	flVersion     = flag.Bool([]string{"v", "--version"}, false, "Print version information and quit")
)

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
