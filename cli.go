// scw interract with Scaleway from the command line
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	flag "github.com/docker/docker/pkg/mflag"
)

var (
	config *Config
)

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

// Options returns a string describing options of the command
func (c *Command) Options() string {
	var options string
	visitor := func(flag *flag.Flag) {
		name := strings.Join(flag.Names, ", -")
		options += fmt.Sprintf("  -%-12s %s (%s)\n", name, flag.Usage, flag.DefValue)
	}
	c.Flag.VisitAll(visitor)
	if len(options) == 0 {
		options = ""
	}
	return options
}

// GetServer returns exactly one server matching or dies
func (cmd *Command) GetServer(needle string) string {
	servers, err := cmd.API.ResolveServer(needle)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to resolve server %s: %s\n", needle, err)
		os.Exit(1)
	}
	if len(servers) == 1 {
		return servers[0]
	}
	if len(servers) == 0 {
		fmt.Fprintf(os.Stderr, "No such server: %s\n", needle)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Too many candidates for %s (%d)\n", needle, len(servers))
	for _, identifier := range servers {
		// FIXME: also print the name
		fmt.Fprintf(os.Stderr, "%s\n", identifier)
	}
	os.Exit(1)
	return ""
}

// GetImage returns exactly one image matching or dies
func (cmd *Command) GetImage(needle string) string {
	images, err := cmd.API.ResolveImage(needle)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to resolve image %s: %s\n", needle, err)
		os.Exit(1)
	}
	if len(images) == 1 {
		return images[0]
	}
	if len(images) == 0 {
		fmt.Fprintf(os.Stderr, "No such image: %s\n", needle)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Too many candidates for %s (%d)\n", needle, len(images))
	for _, identifier := range images {
		// FIXME: also print the name
		fmt.Fprintf(os.Stderr, "%s\n", identifier)
	}
	os.Exit(1)
	return ""
}

// GetBootscript returns exactly one bootscript matching or dies
func (cmd *Command) GetBootscript(needle string) string {
	bootscripts, err := cmd.API.ResolveBootscript(needle)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to resolve bootscript %s: %s\n", needle, err)
		os.Exit(1)
	}
	if len(bootscripts) == 1 {
		return bootscripts[0]
	}
	if len(bootscripts) == 0 {
		fmt.Fprintf(os.Stderr, "No such bootscript: %s\n", needle)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Too many candidates for %s (%d)\n", needle, len(bootscripts))
	for _, identifier := range bootscripts {
		// FIXME: also print the name
		fmt.Fprintf(os.Stderr, "%s\n", identifier)
	}
	os.Exit(1)
	return ""
}

var commands = []*Command{
	cmdCommit,
	cmdCreate,
	cmdExec,
	cmdHelp,
	cmdImages,
	cmdInfo,
	cmdInspect,
	cmdLogin,
	cmdPs,
	cmdRestart,
	cmdRm,
	cmdRmi,
	cmdStart,
	cmdStop,
	cmdVersion,
	cmdWait,
}

var (
	flAPIEndPoint *string
	flDebug       = flag.Bool([]string{"D", "-debug"}, false, "Enable debug mode")
	flVersion     = flag.Bool([]string{"v", "--version"}, false, "Print version information and quit")
)

func main() {
	config, _ = getConfig()

	flAPIEndPoint = flag.String([]string{"-api-endpoint"}, config.APIEndPoint, "Set the API endpoint")
	flag.Parse()

	if *flVersion {
		showVersion()
		return
	}

	os.Setenv("scaleway_api_endpoint", *flAPIEndPoint)

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
				fmt.Fprintf(os.Stderr, "usage: scw %s\n", cmd.UsageLine)
				os.Exit(1)
			}
			api, err := getScalewayAPI()
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to initialize scw api: %s\n", err)
				os.Exit(1)
			}
			cmd.API = api
			cmd.Exec(cmd, cmd.Flag.Args())
			cmd.API.Sync()
			os.Exit(0)
		}
	}

	fmt.Fprintf(os.Stderr, "scw: unknown subcommand %s\nRun 'scw help for usage.\n", name)
	os.Exit(1)
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
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/.scwrc", u.HomeDir), nil
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
	fmt.Printf("scw version %s, build %s\n", "FIXME:VERSION", "FIXME:GITCOMMIT")
}
