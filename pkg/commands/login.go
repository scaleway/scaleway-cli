// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/config"
	"github.com/scaleway/scaleway-cli/pkg/scwversion"
)

// LoginArgs are arguments passed to `RunLogin`
type LoginArgs struct {
	Organization string
	Token        string
	SSHKey       string
	SkipSSHKey   bool
}

// selectKey allows to choice a key in ~/.ssh
func selectKey(args *LoginArgs) error {
	fmt.Println("Do you want to upload an SSH key ?")
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
		if err := promptUser("Which [id]: ", &args.SSHKey, true); err != nil {
			return err
		}
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

func getToken(connect api.ScalewayConnect) (string, error) {
	FakeConnection, err := api.NewScalewayAPI(api.ComputeAPI, api.AccountAPI, "", "", scwversion.UserAgent())
	if err != nil {
		return "", fmt.Errorf("Unable to create a fake ScalewayAPI: %s", err)
	}
	FakeConnection.EnableAccountAPI()
	FakeConnection.SetPassword(connect.Password)

	resp, err := FakeConnection.PostResponse("tokens", connect)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", fmt.Errorf("unable to connect %v", err)
	}

	// Succeed POST code
	if resp.StatusCode != 201 {
		return "", fmt.Errorf("[%d] maybe your email or your password is not valid", resp.StatusCode)
	}
	var data api.ScalewayConnectResponse

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return "", err
	}
	return data.Token.ID, nil
}

func getOrganization(token string, email string) (string, error) {
	FakeConnection, err := api.NewScalewayAPI(api.ComputeAPI, api.AccountAPI, "", token, scwversion.UserAgent())
	if err != nil {
		return "", fmt.Errorf("Unable to create a fake ScalewayAPI: %s", err)
	}
	data, err := FakeConnection.GetOrganization()
	if err != nil {
		return "", err
	}

	orgaID := ""

	for _, orga := range data.Organizations {
		for _, user := range orga.Users {
			if user.Email == email {
				for i := range user.Organizations {
					if user.Organizations[i].Name != "OCS" {
						orgaID = user.Organizations[i].ID
						goto exit
					}
				}
			}
		}
	}
	if orgaID == "" {
		return "", fmt.Errorf("Unable to find your organization")
	}
exit:
	return orgaID, nil
}

func connectAPI() (string, string, error) {
	email := ""
	password := ""
	orga := ""
	token := ""
	hostname, err := os.Hostname()
	if err != nil {
		return "", "", fmt.Errorf("unable to get your Hostname %v", err)
	}
	if err := promptUser("Login (cloud.scaleway.com): ", &email, true); err != nil {
		return "", "", err
	}
	if err := promptUser("Password: ", &password, false); err != nil {
		return "", "", err
	}

	connect := api.ScalewayConnect{
		Email:       strings.Trim(email, "\n"),
		Password:    strings.Trim(password, "\n"),
		Expires:     false,
		Description: strings.Join([]string{"scw", hostname}, "-"),
	}
	token, err = getToken(connect)
	if err != nil {
		return "", "", err
	}
	orga, err = getOrganization(token, connect.Email)
	if err != nil {
		return "", "", err
	}
	return orga, token, nil
}

// uploadSSHKeys uploads an SSH Key
func uploadSSHKeys(apiConnection *api.ScalewayAPI, newKey string) {
	user, err := apiConnection.GetUser()
	if err != nil {
		logrus.Errorf("Unable to contact ScalewayAPI: %s", err)
	} else {
		user.SSHPublicKeys = append(user.SSHPublicKeys, api.ScalewayKeyDefinition{Key: strings.Trim(newKey, "\n")})

		SSHKeys := api.ScalewayUserPatchSSHKeyDefinition{
			SSHPublicKeys: user.SSHPublicKeys,
		}
		for i := range SSHKeys.SSHPublicKeys {
			SSHKeys.SSHPublicKeys[i].Fingerprint = ""
		}

		userID, err := apiConnection.GetUserID()
		if err != nil {
			logrus.Errorf("Unable to get userID: %s", err)
		} else {
			if err = apiConnection.PatchUserSSHKey(userID, SSHKeys); err != nil {
				logrus.Errorf("Unable to patch SSHkey: %v", err)
			}
		}
	}
}

// RunLogin is the handler for 'scw login'
func RunLogin(ctx CommandContext, args LoginArgs) error {
	if args.Organization == "" || args.Token == "" {
		var err error

		args.Organization, args.Token, err = connectAPI()
		if err != nil {
			return err
		}
	}

	cfg := &config.Config{
		ComputeAPI:   api.ComputeAPI,
		AccountAPI:   api.AccountAPI,
		Organization: strings.Trim(args.Organization, "\n"),
		Token:        strings.Trim(args.Token, "\n"),
	}

	apiConnection, err := api.NewScalewayAPI(cfg.ComputeAPI, cfg.AccountAPI, cfg.Organization, cfg.Token, scwversion.UserAgent())
	if err != nil {
		return fmt.Errorf("Unable to create ScalewayAPI: %s", err)
	}
	err = apiConnection.CheckCredentials()
	if err != nil {
		return fmt.Errorf("Unable to contact ScalewayAPI: %s", err)
	}
	if !args.SkipSSHKey {
		if err := selectKey(&args); err != nil {
			logrus.Errorf("Unable to select a key: %v", err)
		} else {
			if args.SSHKey != "" {
				uploadSSHKeys(apiConnection, args.SSHKey)
			}
		}
	}
	return cfg.Save()
}

func promptUser(prompt string, output *string, echo bool) error {
	// FIXME: should use stdin/stdout from command context
	fmt.Fprintf(os.Stdout, prompt)
	os.Stdout.Sync()

	if !echo {
		b, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return fmt.Errorf("Unable to prompt for password: %s", err)
		}
		*output = string(b)
		fmt.Fprintf(os.Stdout, "\n")
	} else {
		reader := bufio.NewReader(os.Stdin)
		*output, _ = reader.ReadString('\n')
	}
	return nil
}
