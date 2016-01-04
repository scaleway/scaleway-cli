// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

// Manage BareMetal Servers from Command Line (as easily as with Docker)

// +build go1.5

package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/scaleway/scaleway-cli/pkg/cli"
)

func main() {
	ec, err := cli.Start(os.Args[1:], nil)
	if err != nil {
		logrus.Fatalf("%s", err)
	}
	os.Exit(ec)
}
