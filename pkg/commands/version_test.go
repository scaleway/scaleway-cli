// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"bytes"
	"testing"

	"github.com/scaleway/scaleway-cli/vendor/github.com/stretchr/testify/assert"
)

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
