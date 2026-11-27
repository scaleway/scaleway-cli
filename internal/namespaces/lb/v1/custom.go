package lb

import (
	"reflect"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
)

const (
	warningKapsuleTaggedMessage = "This resource is auto managed by Kapsule, all your modifications will be overwritten."
	kapsuleTag                  = "kapsule"
)

func warningKapsuleTaggedMessageView() string {
	return terminal.Style("Warning: ", color.Bold, color.FgRed) + warningKapsuleTaggedMessage
}

func GetCommands() *core.Commands {
	humanv2.RegisterEnumFormat(humanv2.EnumFormats[lb.LBTypeStock]{
		lb.LBTypeStockLowStock:   humanv2.EnumFormat{Color: color.FgYellow, Value: "low stock"},
		lb.LBTypeStockAvailable:  humanv2.EnumFormat{Color: color.FgGreen, Value: "available"},
		lb.LBTypeStockOutOfStock: humanv2.EnumFormat{Color: color.FgRed, Value: "out of stock"},
	})

	humanv2.RegisterEnumFormat(humanv2.EnumFormats[lb.LBStatus]{
		lb.LBStatusError:     humanv2.EnumFormat{Color: color.FgRed, Value: "error"},
		lb.LBStatusLocked:    humanv2.EnumFormat{Color: color.FgRed, Value: "locked"},
		lb.LBStatusMigrating: humanv2.EnumFormat{Color: color.FgBlue, Value: "migrating"},
		lb.LBStatusPending:   humanv2.EnumFormat{Color: color.FgBlue, Value: "pending"},
		lb.LBStatusReady:     humanv2.EnumFormat{Color: color.FgGreen, Value: "ready"},
		lb.LBStatusStopped:   humanv2.EnumFormat{Color: color.Faint, Value: "stopped"},
		lb.LBStatusUnknown:   humanv2.EnumFormat{Color: color.Faint, Value: "unknown"},
	})

	humanv2.RegisterEnumFormat(humanv2.EnumFormats[lb.CertificateStatus]{
		lb.CertificateStatusError:   humanv2.EnumFormat{Color: color.FgRed, Value: "error"},
		lb.CertificateStatusPending: humanv2.EnumFormat{Color: color.FgBlue, Value: "pending"},
		lb.CertificateStatusReady:   humanv2.EnumFormat{Color: color.FgGreen, Value: "ready"},
	})

	humanv2.RegisterEnumFormat(humanv2.EnumFormats[lb.ACLActionType]{
		lb.ACLActionTypeAllow:    humanv2.EnumFormat{Color: color.FgGreen, Value: "allow"},
		lb.ACLActionTypeDeny:     humanv2.EnumFormat{Color: color.FgRed, Value: "deny"},
		lb.ACLActionTypeRedirect: humanv2.EnumFormat{Value: "redirect"},
	})

	humanv2.RegisterEnumFormat(humanv2.EnumFormats[lb.BackendServerStatsHealthCheckStatus]{
		lb.BackendServerStatsHealthCheckStatusPassed:   humanv2.EnumFormat{Color: color.FgGreen, Value: "passed"},
		lb.BackendServerStatsHealthCheckStatusFailed:   humanv2.EnumFormat{Color: color.FgRed, Value: "failed"},
		lb.BackendServerStatsHealthCheckStatusUnknown:  humanv2.EnumFormat{Color: color.Faint, Value: "unknown"},
		lb.BackendServerStatsHealthCheckStatusNeutral:  humanv2.EnumFormat{Color: color.Faint, Value: "neutral"},
		lb.BackendServerStatsHealthCheckStatusCondpass: humanv2.EnumFormat{Color: color.FgBlue, Value: "condition passed"},
	})
	humanv2.RegisterEnumFormat(humanv2.EnumFormats[lb.BackendServerStatsServerState]{
		lb.BackendServerStatsServerStateStopped:  humanv2.EnumFormat{Color: color.FgRed, Value: "stopped"},
		lb.BackendServerStatsServerStateStarting: humanv2.EnumFormat{Color: color.FgBlue, Value: "starting"},
		lb.BackendServerStatsServerStateRunning:  humanv2.EnumFormat{Color: color.FgGreen, Value: "running"},
		lb.BackendServerStatsServerStateStopping: humanv2.EnumFormat{Color: color.FgBlue, Value: "stopping"},
	})

	humanv2.RegisterLayout(reflect.TypeFor[lb.LB](), humanv2.Layout{
		Sections: []humanv2.Section{
			{Path: "IP", Title: "IPs"},
			{Path: "Instances", Title: "LB Instances"},
		},
	})
	humanv2.RegisterLayout(reflect.TypeFor[lb.Backend](), humanv2.Layout{
		Sections: []humanv2.Section{
			{Path: "HealthCheck"},
			{Path: "Pool"},
			{Path: "LB"},
			{Path: "LB.IP", Title: "IPs"},
			{Path: "LB.Instances", Title: "LB Instances"},
		},
	})
	humanv2.RegisterLayout(reflect.TypeFor[lb.Frontend](), humanv2.Layout{
		Sections: []humanv2.Section{
			{Path: "LB"},
			{Path: "LB.IP", Title: "IPs"},
			{Path: "LB.Instances", Title: "LB Instances"},
			{Path: "Backend"},
			{Path: "Backend.HealthCheck", Title: "Health Check"},
			{Path: "Backend.Pool", Title: "Pool"},
			{Path: "Backend.LB", Title: "Lb"},
			{Path: "Backend.LB.IP", Title: "IPs"},
			{Path: "Backend.LB.Instances", Title: "LB Instances"},
		},
	})
	humanv2.RegisterLayout(reflect.TypeFor[lb.Certificate](), humanv2.Layout{
		Sections: []humanv2.Section{
			{Path: "LB"},
		},
	})
	humanv2.RegisterLayout(reflect.TypeFor[lb.ACL](), humanv2.Layout{
		Sections: []humanv2.Section{
			{Path: "Frontend"},
		},
	})
	humanv2.RegisterLayout(reflect.TypeFor[[]*lb.LB](), humanv2.Layout{
		Sections: []humanv2.Section{{
			Path: "", Main: true,
			Fields: []humanv2.Field{
				{Path: "ID"},
				{Path: "Name"},
				{Path: "Description"},
				{Path: "Status"},
				{Path: "Instances"},
				{Path: "OrganizationID"},
				{Path: "ProjectID"},
				{Path: "IP"},
				{Path: "IpamIPs"},
				{Path: "Tags"},
				{Path: "FrontendCount"},
				{Path: "BackendCount"},
				{Path: "Type"},
				{Path: "SslCompatibilityLevel"},
				{Path: "CreatedAt"},
				{Path: "UpdatedAt"},
				{Path: "PrivateNetworkCount"},
				{Path: "RouteCount"},
				{Path: "Region"},
				{Path: "Zone"},
				{Path: "Srn"},
			},
		}},
	})
	humanv2.RegisterLayout(reflect.TypeFor[[]*lb.PrivateNetwork](), humanv2.Layout{
		Sections: []humanv2.Section{{
			Path: "", Main: true,
			Fields: []humanv2.Field{
				{Path: "IpamIDs"},
				{Path: "DHCPConfig.IPID"},
				{Path: "StaticConfig.IPAddress"},
				{Path: "PrivateNetworkID"},
				{Path: "Status"},
				{Path: "CreatedAt"},
				{Path: "UpdatedAt"},
			},
		}},
	})

	cmds := GetGeneratedCommands()

	cmds.MustFind("lb").Groups = []string{"network"}

	cmds.Add(
		lbWaitCommand(),
	)

	cmds.MustFind("lb", "lb", "create").Override(lbCreateBuilder)
	cmds.MustFind("lb", "lb", "get").Override(lbGetBuilder)
	cmds.MustFind("lb", "lb", "migrate").Override(lbMigrateBuilder)
	cmds.MustFind("lb", "lb", "update").Override(lbUpdateBuilder)
	cmds.MustFind("lb", "lb", "delete").Override(lbDeleteBuilder)
	cmds.MustFind("lb", "lb", "get-stats").Override(lbGetStatsBuilder)

	cmds.MustFind("lb", "backend", "get").Override(backendGetBuilder)
	cmds.MustFind("lb", "backend", "create").Override(backendCreateBuilder)
	cmds.MustFind("lb", "backend", "update").Override(backendUpdateBuilder)
	cmds.MustFind("lb", "backend", "delete").Override(backendDeleteBuilder)
	cmds.MustFind("lb", "backend", "add-servers").Override(backendAddServersBuilder)
	cmds.MustFind("lb", "backend", "remove-servers").Override(backendRemoveServersBuilder)
	cmds.MustFind("lb", "backend", "set-servers").Override(backendSetServersBuilder)
	cmds.MustFind("lb", "backend", "update-healthcheck").Override(backendUpdateHealthcheckBuilder)

	cmds.MustFind("lb", "frontend", "get").Override(frontendGetBuilder)
	cmds.MustFind("lb", "frontend", "create").Override(frontendCreateBuilder)
	cmds.MustFind("lb", "frontend", "update").Override(frontendUpdateBuilder)
	cmds.MustFind("lb", "frontend", "delete").Override(frontendDeleteBuilder)

	cmds.MustFind("lb", "acl", "get").Override(ACLGetBuilder)
	cmds.MustFind("lb", "acl", "create").Override(ACLCreateBuilder)
	cmds.MustFind("lb", "acl", "update").Override(ACLUpdateBuilder)
	cmds.MustFind("lb", "acl", "delete").Override(ACLDeleteBuilder)

	cmds.MustFind("lb", "certificate", "get").Override(certificateGetBuilder)
	cmds.MustFind("lb", "certificate", "create").Override(certificateCreateBuilder)
	cmds.MustFind("lb", "certificate", "update").Override(certificateUpdateBuilder)
	cmds.MustFind("lb", "certificate", "delete").Override(certificateDeleteBuilder)

	cmds.Merge(core.NewCommands(loadBalancerCockpitCommand()))

	return cmds
}
