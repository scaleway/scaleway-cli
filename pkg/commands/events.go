// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"time"

	"github.com/docker/go-units"
)

// EventsArgs are arguments passed to `RunEvents`
type EventsArgs struct{}

// RunEvents is the handler for 'scw events'
func RunEvents(ctx CommandContext, args EventsArgs) error {
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

		fmt.Fprintf(ctx.Stdout, "%s %s: %s (%s %d) %s\n", startedAt, event.HrefFrom, event.Description, event.Status, event.Progress, terminatedAt)
	}
	return nil
}
