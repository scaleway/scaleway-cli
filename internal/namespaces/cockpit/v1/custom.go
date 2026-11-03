package cockpit

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("cockpit").Groups = []string{"monitoring"}

	cmds.MustFind("cockpit", "token", "get").Override(cockpitTokenGetBuilder)
	cmds.Add(cockpitConfigRoot())
	cmds.Add(cockpitConfigGetCommand())
	addCockpitConfigGetSeeAlso(cmds)

	cmds.MustFind("cockpit", "contact-point", "list").
		Override(core.ListViewBuilder("Region", "SendResolvedNotifications"))
	cmds.MustFind("cockpit", "data-source", "list").
		Override(core.ListViewBuilder("ID", "ProjectID", "Name", "URL", "Type", "Origin", "CreatedAt", "UpdatedAt", "SynchronizedWithGrafana", "RetentionDays", "Region", "CurrentMonthUsage"))
	cmds.MustFind("cockpit", "grafana-user", "list").
		Override(core.ListViewBuilder("ID", "Login", "Role", "Password"))
	cmds.MustFind("cockpit", "plan", "list").
		Override(core.ListViewBuilder("Name", "RetentionMetricsInterval", "RetentionLogsInterval", "RetentionTracesInterval", "SampleIngestionPrice", "LogsIngestionPrice", "TracesIngestionPrice", "MonthlyPrice"))
	cmds.MustFind("cockpit", "product-dashboards", "list").
		Override(core.ListViewBuilder("Name", "Title", "URL", "Tags", "Variables"))
	cmds.MustFind("cockpit", "token", "list").
		Override(core.ListViewBuilder("ID", "ProjectID", "Name", "CreatedAt", "UpdatedAt", "Scopes", "SecretKey", "Region"))

	return cmds
}

func addCockpitConfigGetSeeAlso(cmds *core.Commands) {
	for _, cmd := range cmds.GetAll() {
		if cmd.Namespace != "cockpit" || cmd.Run == nil || cmd.Resource == "config" {
			continue
		}

		cmd.SeeAlsos = append(cmd.SeeAlsos, &core.SeeAlso{
			Command: "scw cockpit config get",
			Short:   "Generate a data source configuration snippet",
		})
	}
}
