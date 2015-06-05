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

func NewListOpts() CommandListOpts {
	var values []string
	return CommandListOpts{
		Values: &values,
	}
}

func (opts *CommandListOpts) String() string {
	return fmt.Sprintf("%v", []string((*opts.Values)))
}

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

func (c *Command) PrintUsage() {
	helpMessage, err := commandHelpMessage(c)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Fprintf(os.Stderr, "%s\n", helpMessage)
	os.Exit(1)
}

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

// Examples returns a string describing examples of the command
func (c *Command) ExamplesHelp() string {
	if c.Examples == "" {
		return ""
	}
	return fmt.Sprintf("Examples:\n\n%s", strings.Trim(c.Examples, "\n"))
}

// GetServer returns exactly one server matching or dies
func (cmd *Command) GetServer(needle string) string {
	servers, err := cmd.API.ResolveServer(needle)
	if err != nil {
		log.Fatalf("Unable to resolve server %s: %s", needle, err)
	}
	if len(servers) == 1 {
		return servers[0]
	}
	if len(servers) == 0 {
		log.Fatalf("No such server: %s", needle)
	}
	log.Errorf("Too many candidates for %s (%d)", needle, len(servers))
	for _, identifier := range servers {
		// FIXME: also print the name
		log.Infof("- %s", identifier)
	}
	os.Exit(1)
	return ""
}

// GetSnapshot returns exactly one snapshot matching or dies
func (cmd *Command) GetSnapshot(needle string) string {
	snapshots, err := cmd.API.ResolveSnapshot(needle)
	if err != nil {
		log.Fatalf("Unable to resolve snapshot %s: %s", needle, err)
	}
	if len(snapshots) == 1 {
		return snapshots[0]
	}
	if len(snapshots) == 0 {
		log.Fatalf("No such snapshot: %s", needle)
	}
	log.Errorf("Too many candidates for %s (%d)", needle, len(snapshots))
	for _, identifier := range snapshots {
		// FIXME: also print the name
		log.Infof("- %s", identifier)
	}
	os.Exit(1)
	return ""
}

// GetImage returns exactly one image matching or dies
func (cmd *Command) GetImage(needle string) string {
	images, err := cmd.API.ResolveImage(needle)
	if err != nil {
		log.Fatalf("Unable to resolve image %s: %s", needle, err)
	}
	if len(images) == 1 {
		return images[0]
	}
	if len(images) == 0 {
		log.Fatalf("No such image: %s", needle)
	}
	log.Errorf("Too many candidates for %s (%d)", needle, len(images))
	for _, identifier := range images {
		// FIXME: also print the name
		log.Infof("- %s", identifier)
	}
	os.Exit(1)
	return ""
}

// GetBootscript returns exactly one bootscript matching or dies
func (cmd *Command) GetBootscript(needle string) string {
	bootscripts, err := cmd.API.ResolveBootscript(needle)
	if err != nil {
		log.Fatalf("Unable to resolve bootscript %s: %s", needle, err)
	}
	if len(bootscripts) == 1 {
		return bootscripts[0]
	}
	if len(bootscripts) == 0 {
		log.Fatalf("No such bootscript: %s", needle)
	}
	log.Errorf("Too many candidates for %s (%d)", needle, len(bootscripts))
	for _, identifier := range bootscripts {
		// FIXME: also print the name
		log.Infof("- %s", identifier)
	}
	os.Exit(1)
	return ""
}

var commands = []*Command{
	cmdAttach,
	cmdCommit,
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
	var cfg_err error
	config, cfg_err = getConfig()
	if cfg_err != nil && !os.IsNotExist(cfg_err) {
		log.Fatalf("Unable to open .scwrc config file: %v", cfg_err)
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
				if cfg_err != nil {
					log.Fatalf("Unable to open .scwrc config file: %v", cfg_err)
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

	log.Fatalf("scw: unknown subcommand %s\nRun 'scw help for usage.", name)
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
	scwrc_path, err := GetConfigFilePath()
	if err != nil {
		return nil, err
	}
	file, err := ioutil.ReadFile(scwrc_path)
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
