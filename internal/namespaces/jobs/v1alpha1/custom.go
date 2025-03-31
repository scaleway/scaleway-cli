package jobs

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	jobs "github.com/scaleway/scaleway-sdk-go/api/jobs/v1alpha1"
)

var jobRunStateMarshalSpecs = human.EnumMarshalSpecs{
	jobs.JobRunStateUnknownState: &human.EnumMarshalSpec{Attribute: color.Faint},
	jobs.JobRunStateQueued:       &human.EnumMarshalSpec{Attribute: color.FgBlue},
	jobs.JobRunStateScheduled:    &human.EnumMarshalSpec{Attribute: color.FgBlue},
	jobs.JobRunStateRunning:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
	jobs.JobRunStateSucceeded:    &human.EnumMarshalSpec{Attribute: color.FgBlue},
	jobs.JobRunStateFailed:       &human.EnumMarshalSpec{Attribute: color.FgRed},
	jobs.JobRunStateCanceled:     &human.EnumMarshalSpec{Attribute: color.FgRed},
}

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("jobs").Groups = []string{"serverless"}

	human.RegisterMarshalerFunc(
		jobs.JobRunState(""),
		human.EnumMarshalFunc(jobRunStateMarshalSpecs),
	)

	cmds.Merge(core.NewCommands(
		jobsRunWait(),
	))

	definitionStartBuilder(cmds.MustFind("jobs", "definition", "start"))

	return cmds
}
