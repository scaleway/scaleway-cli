package jobs

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	jobs "github.com/scaleway/scaleway-sdk-go/api/jobs/v1alpha1"
)

var (
	jobRunStateMarshalSpecs = human.EnumMarshalSpecs{
		jobs.JobRunStateUnknownState: &human.EnumMarshalSpec{Attribute: color.Faint},
		jobs.JobRunStateQueued:       &human.EnumMarshalSpec{Attribute: color.FgBlue},
		jobs.JobRunStateScheduled:    &human.EnumMarshalSpec{Attribute: color.FgBlue},
		jobs.JobRunStateRunning:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
		jobs.JobRunStateSucceeded:    &human.EnumMarshalSpec{Attribute: color.FgBlue},
		jobs.JobRunStateFailed:       &human.EnumMarshalSpec{Attribute: color.FgRed},
		jobs.JobRunStateCanceled:     &human.EnumMarshalSpec{Attribute: color.FgRed},
	}
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	human.RegisterMarshalerFunc(jobs.JobRunState(""), human.EnumMarshalFunc(jobRunStateMarshalSpecs))

	return cmds
}
