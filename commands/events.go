package commands

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/units"

	types "github.com/scaleway/scaleway-cli/commands/types"
)

var cmdEvents = &types.Command{
	Exec:        runEvents,
	UsageLine:   "events [OPTIONS]",
	Description: "Get real time events from the API",
	Help:        "Get real time events from the API.",
}

func init() {
	cmdEvents.Flag.BoolVar(&eventsHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var eventsHelp bool // -h, --help flag

func runEvents(cmd *types.Command, args []string) {
	if eventsHelp {
		cmd.PrintUsage()
	}
	if len(args) != 0 {
		cmd.PrintShortUsage()
	}

	events, err := cmd.API.GetTasks()
	if err != nil {
		log.Fatalf("unable to fetch tasks from the Scaleway API: %v", err)
	}

	for _, event := range *events {
		startedAt, err := time.Parse("2006-01-02T15:04:05.000000+00:00", event.StartDate)
		if err != nil {
			log.Fatalf("unable to parse started date from the Scaleway API: %v", err)
		}

		terminatedAt := ""
		if event.TerminationDate != "" {
			terminatedAtTime, err := time.Parse("2006-01-02T15:04:05.000000+00:00", event.TerminationDate)
			if err != nil {
				log.Fatalf("unable to parse terminated date from the Scaleway API: %v", err)
			}
			terminatedAt = units.HumanDuration(time.Now().UTC().Sub(terminatedAtTime))
		}

		fmt.Printf("%s %s: %s (%s %s) %s\n", startedAt, event.HrefFrom, event.Description, event.Status, event.Progress, terminatedAt)
	}
}
