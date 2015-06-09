// scw interract with Scaleway from the command line
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	log "github.com/Sirupsen/logrus"
	flag "github.com/docker/docker/pkg/mflag"
	"github.com/scaleway/scaleway-cli/scwversion"
)

var (
	config *Config
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

// Command is a Scaleway command
type Command struct {
	// Exec executes the command
	Exec func(cmd *Command, args []string)

	// Usage is the one-line usage message.
	UsageLine string

	// Description is the description of the command
	Description string

	// Help is the full description of the command
	Help string

	// Examples are some examples of the command
	Examples string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet

	// Hidden is a flat to hide command from global help commands listing
	Hidden bool

	// API is the interface used to communicate with Scaleway's API
	API *ScalewayAPI
}

// Name returns the command's name
func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

var fullHelpTemplate = `
Usage: scw {{.UsageLine}}

{{.Help}}

{{.Options}}
{{.ExamplesHelp}}
`

func commandHelpMessage(cmd *Command) (string, error) {
	t := template.New("full")
	template.Must(t.Parse(fullHelpTemplate))
	var output bytes.Buffer
	err := t.Execute(&output, cmd)
	if err != nil {
		return "", err
	}
	return strings.Trim(output.String(), "\n"), nil
}

func commandUsage(name string) {
}

// PrintUsage prints a full command usage and exits
func (c *Command) PrintUsage() {
	helpMessage, err := commandHelpMessage(c)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Fprintf(os.Stderr, "%s\n", helpMessage)
	os.Exit(1)
}

// PrintShortUsage prints a short command usage and exits
func (c *Command) PrintShortUsage() {
	fmt.Fprintf(os.Stderr, "usage: scw %s. See 'scw %s --help'.\n", c.UsageLine, c.Name())
	os.Exit(1)
}

// Options returns a string describing options of the command
func (c *Command) Options() string {
	var options string
	visitor := func(flag *flag.Flag) {
		name := strings.Join(flag.Names, ", -")
		var optionUsage string
		if flag.DefValue == "" {
			optionUsage = fmt.Sprintf("%s=\"\"", name)
		} else {
			optionUsage = fmt.Sprintf("%s=%s", name, flag.DefValue)
		}
		options += fmt.Sprintf("  -%-20s %s\n", optionUsage, flag.Usage)
	}
	c.Flag.VisitAll(visitor)
	if len(options) == 0 {
		return ""
	}
	return fmt.Sprintf("Options:\n\n%s", options)
}

// ExamplesHelp returns a string describing examples of the command
func (c *Command) ExamplesHelp() string {
	if c.Examples == "" {
		return ""
	}
	return fmt.Sprintf("Examples:\n\n%s", strings.Trim(c.Examples, "\n"))
}

var commands = []*Command{
	cmdAttach,
	cmdCommit,
	cmdCompletion,
	cmdCp,
	cmdCreate,
	cmdEvents,
	cmdExec,
	cmdHelp,
	cmdHistory,
	cmdImages,
	cmdInfo,
	cmdInspect,
	cmdKill,
	cmdLogin,
	cmdLogs,
	cmdPort,
	cmdPs,
	cmdRename,
	cmdRestart,
	cmdRm,
	cmdRmi,
	cmdRun,
	cmdSearch,
	cmdStart,
	cmdStop,
	cmdTag,
	cmdTop,
	cmdVersion,
	cmdWait,
}

var (
	flAPIEndPoint *string
	flDebug       = flag.Bool([]string{"D", "-debug"}, false, "Enable debug mode")
	flVersion     = flag.Bool([]string{"v", "--version"}, false, "Print version information and quit")
)

func main() {
	var cfgErr error
	config, cfgErr = getConfig()
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

	for _, cmd := range commands {
		if cmd.Name() == name {
			cmd.Flag.SetOutput(ioutil.Discard)
			err := cmd.Flag.Parse(args)
			if err != nil {
				log.Fatalf("usage: scw %s", cmd.UsageLine)
			}
			if cmd.Name() != "login" {
				if cfgErr != nil {
					log.Fatalf("Unable to open .scwrc config file: %v", cfgErr)
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
	cmdHelp.Exec(cmdHelp, []string{})
	os.Exit(1)
}

// Config is a Scaleway CLI configuration file
type Config struct {
	// APIEndpoint is the endpoint to the Scaleway API
	APIEndPoint string `json:"api_endpoint"`

	// Organization is the identifier of the Scaleway orgnization
	Organization string `json:"organization"`

	// Token is the authentication token for the Scaleway organization
	Token string `json:"token"`
}

// GetConfigFilePath returns the path to the Scaleway CLI config file
func GetConfigFilePath() (string, error) {
	homeDir := os.Getenv("HOME") // *nix
	if homeDir == "" {           // Windows
		homeDir = os.Getenv("USERPROFILE")
	}
	if homeDir == "" {
		return "", errors.New("User home directory not found.")
	}

	return filepath.Join(homeDir, ".scwrc"), nil
}

// getConfig returns the Scaleway CLI config file for the current user
func getConfig() (*Config, error) {
	scwrcPath, err := GetConfigFilePath()
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
	var config Config
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
func getScalewayAPI() (*ScalewayAPI, error) {
	// We already get config globally, but whis way we can get explicit error when trying to create a ScalewayAPI object
	config, err := getConfig()
	if err != nil {
		return nil, err
	}
	return NewScalewayAPI(os.Getenv("scaleway_api_endpoint"), config.Organization, config.Token)
}

func showVersion() {
	fmt.Printf("scw version %s, build %s\n", scwversion.VERSION, scwversion.GITCOMMIT)
}
