// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"encoding/json"
	"fmt"
	minio "github.com/minio/mc/cmd"
	"github.com/scaleway/scaleway-cli/pkg/config"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type clientHost struct {
	Url       string `json:"url"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	Api       string `json:"api"`
	Lookup    string `json:"lookup"`
}

type clientConfig struct {
	Version string `json:"version"`

	Hosts map[string]clientHost `json:"hosts"`
}

// S3Args are flags for the `runS3` function
type S3Args struct {
	Command []string // command line for minio client
}

// getClientConfigPath returns the path of the client configuration directory
func getClientConfigPath() (string, error) {
	path, err := config.GetConfigFilePath()
	if err != nil {
		return "", err
	}
	name := filepath.Base(path)
	if name[0] == '.' {
		name = name[1:]
	}
	path = filepath.Dir(path)
	path = filepath.Join(path, ".scw")
	path = filepath.Join(path, name)
	return path, nil
}

// s3configure will prompt for accessKey and secretKey of the user and
// save it for ams and par regions
func s3Configure() error {
	path, err := getClientConfigPath()
	if err != nil {
		return err
	}
	err = os.MkdirAll(path, 0600)
	if err != nil {
		return err
	}
	var config clientConfig
	configPath := filepath.Join(path, "config.json")
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		// Set default config values
		config.Version = "9"
		config.Hosts = make(map[string]clientHost)
	} else {
		// Retrieve config file
		err = json.Unmarshal(file, &config)
		if err != nil {
			return err
		}
	}
	// Get credentials from client
	var accessKey string
	var secretKey string
	fmt.Print("AccessKey: ")
	fmt.Scanln(&accessKey)
	fmt.Print("SecretKey: ")
	fmt.Scanln(&secretKey)
	ams := config.Hosts["ams"]
	par := config.Hosts["par"]
	ams.AccessKey = accessKey
	par.AccessKey = accessKey
	ams.SecretKey = secretKey
	par.SecretKey = secretKey
	ams.Url = "https://s3.nl-ams.scw.cloud"
	par.Url = "https://s3.fr-par.scw.cloud"
	ams.Lookup = "auto"
	par.Lookup = "auto"
	ams.Api = "S3v2"
	par.Api = "S3v2"
	config.Hosts["ams"] = ams
	config.Hosts["par"] = par
	configFile, err := os.OpenFile(configPath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer configFile.Close()
	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(config)
	if err != nil {
		return err
	}
	// fmt.Printf("%s\n", config)
	return nil
}

// removeClientEnv will remove any environment variable related to the client
func removeClientEnv() {
	envs := os.Environ()
	for i := 0; i < len(envs); i++ {
		if strings.HasPrefix(envs[i], "MC_") {
			os.Unsetenv(strings.Split(envs[i], "=")[0])
		}
	}
}

// S3 is the handler for 'scw s3'
func S3(ctx CommandContext, args S3Args) error {
	// Client own the config command, but we override it to configure Only
	// the host we want
	if len(args.Command) > 0 && args.Command[0] == "config" {
		return s3Configure()
	}
	// Remove any client configuration environment variable
	removeClientEnv()
	// Get the scw cli configuration path
	confPath, err := getClientConfigPath()
	if err != nil {
		return err
	}
	newArgs := make([]string, 0)
	newArgs = append(newArgs, os.Args[:1]...)
	// Override the client config directory
	newArgs = append(newArgs, "--config-folder")
	newArgs = append(newArgs, confPath)
	newArgs = append(newArgs, "--quiet")
	// Remove any scw cli options
	newArgs = append(newArgs, args.Command...)
	os.Args = newArgs
	// Override args[0] to display `scw s3` in the client help pages
	os.Args[0] = os.Args[0] + " s3"
	minio.Main()
	return nil
}
