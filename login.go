package main

import (
	"encoding/json"
	"os"

	log "github.com/Sirupsen/logrus"
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
		log.Fatalf("usage: scw %s", cmd.UsageLine)
	}
	cfg := &Config{
		APIEndPoint:  "https://api.scaleway.com/",
		Organization: args[0],
		Token:        args[1],
	}
	scwrc_path, err := GetConfigFilePath()
	if err != nil {
		log.Fatalf("can't get scwrc config file path: %v", err)
	}
	scwrc, err := os.OpenFile(scwrc_path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("can't create scwrc config file: %v", err)
	}
	defer scwrc.Close()
	encoder := json.NewEncoder(scwrc)
	err = encoder.Encode(cfg)
	if err != nil {
		log.Fatalf("can't encode scw config file: %v", err)
	}
}
