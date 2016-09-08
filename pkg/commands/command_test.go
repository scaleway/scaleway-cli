// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

// Golang structs for scw commands
package commands

import (
	"fmt"
	"os"
	"testing"

	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/scwversion"
	"github.com/stretchr/testify/assert"
)

func testCommandContext() CommandContext {
	apiClient, err := api.NewScalewayAPI("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", scwversion.UserAgent(), "")
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

func ExampleCommandContext() {
	apiClient, err := api.NewScalewayAPI("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", scwversion.UserAgent(), "")
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

	// Do stuff
	fmt.Println(ctx)
}

func TestCommandContext_Getenv(t *testing.T) {
	ctx := testCommandContext()
	assert.Equal(t, ctx.Getenv("HOME"), os.Getenv("HOME"))
	assert.Equal(t, ctx.Getenv("DONTEXISTS"), "")
}
