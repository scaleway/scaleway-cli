package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var cmdLogin = &Command{
	Exec:        runLogin,
	UsageLine:   "login [options] ORGANIZATION TOKEN",
	Description: "Login generates a configuration file containing credentials",
	Help: `Login generates a configuration file in '/home/$USER/.scwrc'
containing credentials used to interact with the Scaleway API. This
configuration file is automatically used by the 'scw' command.`,
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
	scwrc_path, err := GetConfigFilePath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't get scwrc config file path: %v\n", err)
		os.Exit(1)
	}
	scwrc, err := os.OpenFile(scwrc_path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
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
