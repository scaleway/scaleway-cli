package main

import (
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/pkg/units"
)

var cmdEvents = &Command{
	Exec:        runEvents,
	UsageLine:   "events [OPTIONS]",
	Description: "Get real time events from the API",
	Help:        "Get real time events from the API.",
}

func runEvents(cmd *Command, args []string) {
	events, err := cmd.API.GetTasks()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to fetch tasks from the Scaleway API: %v\n", err)
		os.Exit(1)
	}

	for _, event := range *events {
		startedAt, err := time.Parse("2006-01-02T15:04:05.000000+00:00", event.StartDate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to parse started date from the Scaleway API: %v\n", err)
			os.Exit(1)
		}

		terminatedAt := ""
		if event.TerminationDate != "" {
			terminatedAtTime, err := time.Parse("2006-01-02T15:04:05.000000+00:00", event.TerminationDate)
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to parse terminated date from the Scaleway API: %v\n", err)
				os.Exit(1)
			}
			terminatedAt = units.HumanDuration(time.Now().UTC().Sub(terminatedAtTime))
		}

		fmt.Printf("%s %s: %s (%s %s) %s\n", startedAt, event.HrefFrom, event.Description, event.Status, event.Progress, terminatedAt)
	}
}
