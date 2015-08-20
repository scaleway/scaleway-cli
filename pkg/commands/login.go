// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"
	"github.com/scaleway/scaleway-cli/vendor/golang.org/x/crypto/ssh/terminal"

	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/config"
)

// LoginArgs are arguments passed to `RunLogin`
type LoginArgs struct {
	Organization string
	Token        string
	SSHKey       string
}

// selectKey allows to choice a key in ~/.ssh
func selectKey(args *LoginArgs) error {
	fmt.Println("Do you want to upload a SSH key ?")
	home, err := config.GetHomeDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(home, ".ssh")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("Unable to open your ~/.ssh: %v", err)
	}
	var pubs []string

	for i := range files {
		if filepath.Ext(files[i].Name()) == ".pub" {
			pubs = append(pubs, files[i].Name())
		}
	}
	if len(pubs) == 0 {
		return nil
	}
	fmt.Println("[0] I don't want to upload a key !")
	for i := range pubs {
		fmt.Printf("[%d] %s\n", i+1, pubs[i])
	}
	for {
		promptUser("Which [id]: ", &args.SSHKey, true)
		id, err := strconv.ParseUint(strings.TrimSpace(args.SSHKey), 10, 32)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if int(id) > len(pubs) {
			fmt.Println("Out of range id must be lower than", len(pubs))
			continue
		}
		args.SSHKey = ""
		if id == 0 {
			break
		}
		buff, err := ioutil.ReadFile(filepath.Join(dir, pubs[id-1]))
		if err != nil {
			return fmt.Errorf("Unable to open your key: %v", err)
		}
		args.SSHKey = string(buff[:len(buff)])
		break
	}
	return nil
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

	cfg := &config.Config{
		ComputeAPI:   "https://api.scaleway.com/",
		AccountAPI:   "https://account.scaleway.com/",
		Organization: strings.Trim(args.Organization, "\n"),
		Token:        strings.Trim(args.Token, "\n"),
	}

	apiConnection, err := api.NewScalewayAPI(cfg.ComputeAPI, cfg.AccountAPI, cfg.Organization, cfg.Token)
	if err != nil {
		return fmt.Errorf("Unable to create ScalewayAPI: %s", err)
	}
	err = apiConnection.CheckCredentials()
	if err != nil {
		return fmt.Errorf("Unable to contact ScalewayAPI: %s", err)
	}
	if err := selectKey(&args); err != nil {
		logrus.Errorf("Unable to select a key: %v", err)
	} else {
		if args.SSHKey != "" {
			userID, err := apiConnection.GetUserID()
			if err != nil {
				logrus.Errorf("Unable to contact ScalewayAPI: %s", err)
			} else {

				SSHKey := api.ScalewayUserPatchDefinition{
					SSHPublicKeys: []api.ScalewayUserPatchKeyDefinition{{
						Key: strings.Trim(args.SSHKey, "\n"),
					}},
				}

				if err = apiConnection.PatchUser(userID, SSHKey); err != nil {
					logrus.Errorf("Unable to patch SSHkey: %v", err)
				}
			}
		}
	}
	return cfg.Save()
}

func promptUser(prompt string, output *string, echo bool) {
	// FIXME: should use stdin/stdout from command context
	fmt.Fprintf(os.Stdout, prompt)
	os.Stdout.Sync()

	if !echo {
		b, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			logrus.Fatalf("Unable to prompt for password: %s", err)
		}
		*output = string(b)
		fmt.Fprintf(os.Stdout, "\n")
	} else {
		reader := bufio.NewReader(os.Stdin)
		*output, _ = reader.ReadString('\n')
	}
}
