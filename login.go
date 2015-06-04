package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	log "github.com/Sirupsen/logrus"
)

var cmdLogin = &Command{
	Exec:        runLogin,
	UsageLine:   "login [OPTIONS]",
	Description: "Log in to Scaleway API",
	Help: `Generates a configuration file in '/home/$USER/.scwrc'
containing credentials used to interact with the Scaleway API. This
configuration file is automatically used by the 'scw' commands.`,
}

func promptUser(prompt string, output *string, echo bool) {
	fmt.Fprintf(os.Stdout, prompt)
	os.Stdout.Sync()

	if !echo {
		b, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatalf("Unable to prompt for password: %s", err)
		}
		*output = string(b)
		fmt.Fprintf(os.Stdout, "\n")
	} else {
		reader := bufio.NewReader(os.Stdin)
		*output, _ = reader.ReadString('\n')
	}
}

func init() {
	cmdLogin.Flag.StringVar(&organization, []string{"o", "-organization"}, "", "Organization")
	cmdLogin.Flag.StringVar(&token, []string{"t", "-token"}, "", "Token")
	cmdLogin.Flag.BoolVar(&loginHelp, []string{"h", "-help"}, false, "Print usage")
}

// FLags
var organization string // -o flag
var token string        // -t flag
var loginHelp bool      // -h, --help flag

func runLogin(cmd *Command, args []string) {
	if loginHelp {
		cmd.PrintUsage()
	}
	if len(args) < 1 {
		cmd.PrintShortUsage()
	}

	if len(organization) == 0 && len(token) == 0 {
		promptUser("Organization: ", &organization, true)
		promptUser("Token: ", &token, false)
	}

	cfg := &Config{
		APIEndPoint:  "https://account.scaleway.com/",
		Organization: strings.Trim(organization, "\n"),
		Token:        strings.Trim(token, "\n"),
	}

	api, err := NewScalewayAPI(cfg.APIEndPoint, cfg.Organization, cfg.Token)
	if err != nil {
		log.Fatalf("Unable to create ScalewayAPI: %s", err)
	}
	err = api.CheckCredentials()
	if err != nil {
		log.Fatalf("Unable to contact ScalewayAPI: %s", err)
	}

	scwrc_path, err := GetConfigFilePath()
	if err != nil {
		log.Fatalf("Unable to get scwrc config file path: %s", err)
	}
	scwrc, err := os.OpenFile(scwrc_path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("Unable to create scwrc config file: %s", err)
	}
	defer scwrc.Close()
	encoder := json.NewEncoder(scwrc)
	cfg.APIEndPoint = "https://api.scaleway.com/"
	err = encoder.Encode(cfg)
	if err != nil {
		log.Fatalf("Unable to encode scw config file: %s", err)
	}
}
