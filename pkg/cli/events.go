// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"github.com/Sirupsen/logrus"
	"github.com/scaleway/scaleway-cli/pkg/commands"
)

var cmdEvents = &Command{
	Exec:        cmdExecEvents,
	UsageLine:   "events [OPTIONS]",
	Description: "Get real time events from the API",
	Help:        "Get real time events from the API.",
}

func init() {
	cmdEvents.Flag.BoolVar(&eventsHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var eventsHelp bool // -h, --help flag

func cmdExecEvents(cmd *Command, rawArgs []string) {
	if eventsHelp {
		cmd.PrintUsage()
	}
	if len(rawArgs) != 0 {
		cmd.PrintShortUsage()
	}

	args := commands.EventsArgs{}
	ctx := cmd.GetContext(rawArgs)
	err := commands.RunEvents(ctx, args)
	if err != nil {
		logrus.Fatalf("Cannot execute 'events': %v", err)
	}
}
