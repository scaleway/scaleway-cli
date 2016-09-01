// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

// Package config contains helpers to manage '~/.scwrc'
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/scaleway/scaleway-cli/pkg/scwversion"
)

// Config is a Scaleway CLI configuration file
type Config struct {
	// Organization is the identifier of the Scaleway orgnization
	Organization string `json:"organization"`

	// Token is the authentication token for the Scaleway organization
	Token string `json:"token"`

	// Version is the actual version of scw
	Version string `json:"version"`
}

// Save write the config file
func (c *Config) Save() error {
	scwrcPath, err := GetConfigFilePath()
	if err != nil {
		return fmt.Errorf("Unable to get scwrc config file path: %s", err)
	}
	scwrc, err := os.OpenFile(scwrcPath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0600)
	if err != nil {
		return fmt.Errorf("Unable to create scwrc config file: %s", err)
	}
	defer scwrc.Close()
	c.Version = scwversion.VERSION
	err = json.NewEncoder(scwrc).Encode(c)
	if err != nil {
		return fmt.Errorf("Unable to encode scw config file: %s", err)
	}
	return nil
}

// GetConfig returns the Scaleway CLI config file for the current user
func GetConfig() (*Config, error) {
	scwrcPath, err := GetConfigFilePath()
	if err != nil {
		return nil, err
	}

	// Don't check permissions on Windows, Go knows nothing about them on this platform
	// User profile is to be assumed safe anyway
	if runtime.GOOS != "windows" {
		stat, errStat := os.Stat(scwrcPath)
		// we don't care if it fails, the user just won't see the warning
		if errStat == nil {
			perm := stat.Mode().Perm()
			if perm&0066 != 0 {
				return nil, fmt.Errorf("permissions %#o for .scwrc are too open", perm)
			}
		}
	}

	file, err := ioutil.ReadFile(scwrcPath)
	if err != nil {
		return nil, err
	}
	var config Config

	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// GetConfigFilePath returns the path to the Scaleway CLI config file
func GetConfigFilePath() (string, error) {
	path, err := GetHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(path, ".scwrc"), nil
}

// GetHomeDir returns the path to your home
func GetHomeDir() (string, error) {
	homeDir := os.Getenv("HOME") // *nix
	if homeDir == "" {           // Windows
		homeDir = os.Getenv("USERPROFILE")
	}
	if homeDir == "" {
		return "", errors.New("user home directory not found")
	}
	return homeDir, nil
}
