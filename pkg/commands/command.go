// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

// Package commands contains the workflows behind each commands of the CLI (run, attach, start, exec, commit, ...)
package commands

import (
	"io"
	"os"

	"github.com/scaleway/scaleway-cli/pkg/api"
)

// Streams is used to redirects the streams
type Streams struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// CommandContext is passed to all commands and contains streams, environment, api and arguments
type CommandContext struct {
	Streams

	Env     []string
	RawArgs []string
	API     *api.ScalewayAPI
}

// Getenv returns the equivalent of os.Getenv for the CommandContext.Env
func (c *CommandContext) Getenv(key string) string {
	// FIXME: parse c.Env instead
	return os.Getenv(key)
}
