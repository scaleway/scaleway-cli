package edge_services

import "github.com/scaleway/scaleway-cli/v2/core"

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("edge-services").Groups = []string{"network"}

	listCmd := cmds.MustFind("edge-services", "plan", "list")
	listCmd.Override(planListBuilder)
	listCmd.Override(
		core.ListViewBuilder(
			"PlanName",
			"PackageGb",
			"PipelineLimit",
			"WafRequests",
			"BackendLimit",
			"WildcardDomain",
		),
	)

	cmds.Merge(core.NewCommands(
		edgeServicesRouteRulesEditCommand(),
	))

	cmds.MustFind("edge-services", "pipeline", "list").
		Override(core.ListViewBuilder("ID", "Name", "Description", "Status", "Errors", "ProjectID", "OrganizationID", "CreatedAt", "UpdatedAt", "VpcEndpointIDs"))
	cmds.MustFind("edge-services", "purge-request", "list").
		Override(core.ListViewBuilder("ID", "PipelineID", "Status", "Assets", "All", "CreatedAt", "UpdatedAt"))

	cmds.MustFind("edge-services", "backend-stage", "list").
		Override(core.ListViewBuilder("ID", "PipelineID", "Status", "CreatedAt", "UpdatedAt"))
	cmds.MustFind("edge-services", "cache-stage", "list").
		Override(core.ListViewBuilder("ID", "PipelineID", "FallbackTTL", "IncludeCookies", "Status", "CreatedAt", "UpdatedAt", "WafStageID", "BackendStageID", "RouteStageID"))
	cmds.MustFind("edge-services", "dns-stage", "list").
		Override(core.ListViewBuilder("ID", "DefaultFqdn", "DefaultPrivateFqdn", "Fqdns", "Type", "PipelineID", "Status", "CreatedAt", "UpdatedAt", "TLSStageID", "CacheStageID", "BackendStageID", "WildcardDomain", "FullPrivate"))
	cmds.MustFind("edge-services", "route-stage", "list").
		Override(core.ListViewBuilder("ID", "PipelineID", "WafStageID", "BackendStageID", "Status", "CreatedAt", "UpdatedAt"))
	cmds.MustFind("edge-services", "tls-stage", "list").
		Override(core.ListViewBuilder("ID", "Secrets", "ManagedCertificate", "PipelineID", "CertificateExpiresAt", "Status", "CreatedAt", "UpdatedAt", "CacheStageID", "BackendStageID", "WafStageID", "RouteStageID"))
	cmds.MustFind("edge-services", "waf-stage", "list").
		Override(core.ListViewBuilder("ID", "PipelineID", "Mode", "ParanoiaLevel", "Status", "CreatedAt", "UpdatedAt", "BackendStageID", "ExclusionRules"))

	return cmds
}
