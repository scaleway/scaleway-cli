package jobs

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	jobs "github.com/scaleway/scaleway-sdk-go/api/jobs/v1alpha2"
)

var jobRunStateFormat = humanv2.EnumFormats[jobs.JobRunState]{
	jobs.JobRunStateUnknownState: humanv2.EnumFormat{Color: color.Faint},
	jobs.JobRunStateQueued:       humanv2.EnumFormat{Color: color.FgBlue},
	jobs.JobRunStateInitialized:  humanv2.EnumFormat{Color: color.FgBlue},
	jobs.JobRunStateValidated:    humanv2.EnumFormat{Color: color.FgBlue},
	jobs.JobRunStateRunning:      humanv2.EnumFormat{Color: color.FgBlue},
	jobs.JobRunStateSucceeded:    humanv2.EnumFormat{Color: color.FgGreen},
	jobs.JobRunStateFailed:       humanv2.EnumFormat{Color: color.FgRed},
	jobs.JobRunStateInterrupting: humanv2.EnumFormat{Color: color.FgBlue},
	jobs.JobRunStateInterrupted:  humanv2.EnumFormat{Color: color.FgRed},
}

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("jobs").Groups = []string{"serverless"}

	humanv2.RegisterEnumFormat(jobRunStateFormat)
	cmds.Merge(core.NewCommands(
		jobsRunWait(),
	))

	definitionStartBuilder(cmds.MustFind("jobs", "definition", "start"))

	return cmds
}
