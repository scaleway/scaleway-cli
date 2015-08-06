// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"time"

	log "github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"
	"github.com/scaleway/scaleway-cli/vendor/github.com/docker/docker/pkg/units"

	types "github.com/scaleway/scaleway-cli/pkg/commands/types"
)

var cmdEvents = &types.Command{
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

type EventsArgs struct{}

func cmdExecEvents(cmd *types.Command, rawArgs []string) {
	if eventsHelp {
		cmd.PrintUsage()
	}
	if len(rawArgs) != 0 {
		cmd.PrintShortUsage()
	}

	args := EventsArgs{}
	ctx := cmd.GetContext(rawArgs)
	err := RunEvents(ctx, args)
	if err != nil {
		log.Fatalf("Cannot execute 'events': %v", err)
	}
}

// RunEvents is the handler for 'scw events'
func RunEvents(ctx types.CommandContext, args EventsArgs) error {
	events, err := ctx.API.GetTasks()
	if err != nil {
		return fmt.Errorf("unable to fetch tasks from the Scaleway API: %v", err)
	}

	for _, event := range *events {
		startedAt, err := time.Parse("2006-01-02T15:04:05.000000+00:00", event.StartDate)
		if err != nil {
			return fmt.Errorf("unable to parse started date from the Scaleway API: %v", err)
		}

		terminatedAt := ""
		if event.TerminationDate != "" {
			terminatedAtTime, err := time.Parse("2006-01-02T15:04:05.000000+00:00", event.TerminationDate)
			if err != nil {
				return fmt.Errorf("unable to parse terminated date from the Scaleway API: %v", err)
			}
			terminatedAt = units.HumanDuration(time.Now().UTC().Sub(terminatedAtTime))
		}

		fmt.Printf("%s %s: %s (%s %d) %s\n", startedAt, event.HrefFrom, event.Description, event.Status, event.Progress, terminatedAt)
	}
	return nil
}
