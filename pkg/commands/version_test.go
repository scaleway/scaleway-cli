// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"bytes"
	"os"
	"testing"

	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/stretchr/testify/assert"
)

func ExampleCommandContext() CommandContext {
	apiClient, err := api.NewScalewayAPI("https://example.org/", "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx")
	if err != nil {
		panic(err)
	}

	ctx := CommandContext{
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Env:     []string{},
		RawArgs: []string{},
		API:     apiClient,
	}
	return ctx
}

func TestToto(t *testing.T) {
	ctx := ExampleCommandContext()
	var buf bytes.Buffer
	ctx.Stdout = &buf

	args := VersionArgs{}

	err := Version(ctx, args)

	assert.Nil(t, err)
	assert.Contains(t, buf.String(), "Client version: ")
	assert.Contains(t, buf.String(), "Go version (client): ")
	assert.Contains(t, buf.String(), "Git commit (client): ")
	assert.Contains(t, buf.String(), "OS/Arch (client): ")
}

func ExampleVersion() {
	ctx := ExampleCommandContext()
	args := VersionArgs{}
	Version(ctx, args)
}
