// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package api

import (
	"encoding/json"
	"fmt"
	"github.com/scaleway/scaleway-cli/pkg/utils"
	"os"
)

// Config is a Scaleway CLI configuration file
type Config struct {
	// ComputeAPI is the endpoint to the Scaleway API
	ComputeAPI string `json:"api_endpoint"`

	// AccountAPI is the endpoint to the Scaleway Account API
	AccountAPI string `json:"account_endpoint"`

	// Organization is the identifier of the Scaleway orgnization
	Organization string `json:"organization"`

	// Token is the authentication token for the Scaleway organization
	Token string `json:"token"`
}

// Save write the config file
func (c *Config) Save() error {
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
	err = encoder.Encode(c)
	if err != nil {
		return fmt.Errorf("Unable to encode scw config file: %s", err)
	}
	return nil
}
