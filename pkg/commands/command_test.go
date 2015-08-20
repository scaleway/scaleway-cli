// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

// Golang structs for scw commands
package commands

import (
	"bytes"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/config"
	"github.com/scaleway/scaleway-cli/vendor/github.com/stretchr/testify/assert"
)

func ExampleCommandContext() CommandContext {
	apiClient, err := api.NewScalewayAPI("https://example.org/", "https://example.org/", "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx")
	if err != nil {
		panic(err)
	}

	ctx := CommandContext{
		Streams: Streams{
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		},
		Env: []string{
			"HOME" + os.Getenv("HOME"),
		},
		RawArgs: []string{},
		API:     apiClient,
	}
	return ctx
}

func TestCommandContext_Getenv(t *testing.T) {
	ctx := ExampleCommandContext()
	assert.Equal(t, ctx.Getenv("HOME"), os.Getenv("HOME"))
	assert.Equal(t, ctx.Getenv("DONTEXISTS"), "")
}

func RealAPIContext() *CommandContext {
	config, err := config.GetConfig()
	if err != nil {
		logrus.Warnf("RealAPIContext: failed to call config.GetConfig(): %v", err)
		return nil
	}

	apiClient, err := api.NewScalewayAPI(config.ComputeAPI, config.AccountAPI, config.Organization, config.Token)
	if err != nil {
		logrus.Warnf("RealAPIContext: failed to call api.NewScalewayAPI(): %v", err)
		return nil
	}

	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	ctx := CommandContext{
		Streams: Streams{
			Stdin:  os.Stdin,
			Stdout: &stdout,
			Stderr: &stderr,
		},
		Env: []string{
			"HOME" + os.Getenv("HOME"),
		},
		RawArgs: []string{},
		API:     apiClient,
	}
	return &ctx
}
