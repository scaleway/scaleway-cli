package jobs

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	jobs "github.com/scaleway/scaleway-sdk-go/api/jobs/v1alpha2"
)

var jobRunStateMarshalSpecs = human.EnumMarshalSpecs{
	jobs.JobRunStateUnknownState: &human.EnumMarshalSpec{Attribute: color.Faint},
	jobs.JobRunStateQueued:       &human.EnumMarshalSpec{Attribute: color.FgBlue},
	jobs.JobRunStateInitialized:  &human.EnumMarshalSpec{Attribute: color.FgBlue},
	jobs.JobRunStateValidated:    &human.EnumMarshalSpec{Attribute: color.FgBlue},
	jobs.JobRunStateRunning:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
	jobs.JobRunStateSucceeded:    &human.EnumMarshalSpec{Attribute: color.FgGreen},
	jobs.JobRunStateFailed:       &human.EnumMarshalSpec{Attribute: color.FgRed},
	jobs.JobRunStateInterrupting: &human.EnumMarshalSpec{Attribute: color.FgBlue},
	jobs.JobRunStateInterrupted:  &human.EnumMarshalSpec{Attribute: color.FgRed},
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
