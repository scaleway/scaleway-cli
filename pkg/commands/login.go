// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/scaleway/scaleway-cli/vendor/golang.org/x/crypto/ssh/terminal"

	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/utils"
)

type LoginArgs struct {
	Organization string
	Token        string
}

// RunLogin is the handler for 'scw login'
func RunLogin(ctx CommandContext, args LoginArgs) error {
	if args.Organization == "" {
		fmt.Println("You can get your credentials on https://cloud.scaleway.com/#/credentials")
		promptUser("Organization (access key): ", &args.Organization, true)
	}
	if args.Token == "" {
		promptUser("Token: ", &args.Token, false)
	}

	cfg := &api.Config{
		APIEndPoint:  "https://account.scaleway.com/",
		Organization: strings.Trim(args.Organization, "\n"),
		Token:        strings.Trim(args.Token, "\n"),
	}

	api, err := api.NewScalewayAPI(cfg.APIEndPoint, cfg.Organization, cfg.Token)
	if err != nil {
		return fmt.Errorf("Unable to create ScalewayAPI: %s", err)
	}
	err = api.CheckCredentials()
	if err != nil {
		return fmt.Errorf("Unable to contact ScalewayAPI: %s", err)
	}

	scwrcPath, err := utils.GetConfigFilePath()
	if err != nil {
		return fmt.Errorf("Unable to get scwrc config file path: %s", err)
	}
	scwrc, err := os.OpenFile(scwrcPath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0600)
	if err != nil {
		return fmt.Errorf("Unable to create scwrc config file: %s", err)
	}
	defer scwrc.Close()
	encoder := json.NewEncoder(scwrc)
	cfg.APIEndPoint = "https://api.scaleway.com/"
	err = encoder.Encode(cfg)
	if err != nil {
		return fmt.Errorf("Unable to encode scw config file: %s", err)
	}
	return nil
}

func promptUser(prompt string, output *string, echo bool) {
	// FIXME: should use stdin/stdout from command context
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
