package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
)

var cmdLogin = &Command{
	Exec:        runLogin,
	UsageLine:   "login organization token",
	Description: "Login generates a configuration file containing credentials",
	Help: `
Login generates a configuration file in '/home/$USER/.scwrc'
containing credentials used to interact with the Scaleway API. This
configuration file is automatically used by the 'scw' command.
`,
}

type Config struct {
	// APIEndpoint is the endpoint to the Scaleway API
	APIEndPoint string

	// Organization is the identifier of the Scaleway orgnization
	Organization string

	// Token is the authentication token for the Scaleway organization
	Token string
}

func runLogin(cmd *Command, args []string) {
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: scw %s\n", cmd.UsageLine)
		os.Exit(1)
	}
	cfg := &Config{
		APIEndPoint:  "https://api.cloud.online.net/",
		Organization: args[0],
		Token:        args[1],
	}
	u, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't get current user: %v\n", err)
		os.Exit(1)
	}
	scwrc, err := os.OpenFile(fmt.Sprintf("%s/.scwrc", u.HomeDir), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't create scwrc config file: %v\n", err)
		os.Exit(1)
	}
	defer scwrc.Close()
	encoder := json.NewEncoder(scwrc)
	err = encoder.Encode(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't encode scw config file: %v\n", err)
		os.Exit(1)
	}
}
