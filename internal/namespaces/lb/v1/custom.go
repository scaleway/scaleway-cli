package lb

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
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
	human.RegisterMarshalerFunc(lb.LBTypeStock(""), human.EnumMarshalFunc(lbTypeStockMarshalSpecs))
	human.RegisterMarshalerFunc(lb.LBStatus(""), human.EnumMarshalFunc(lbStatusMarshalSpecs))
	human.RegisterMarshalerFunc(
		lb.CertificateStatus(""),
		human.EnumMarshalFunc(certificateStatusMarshalSpecs),
	)
	human.RegisterMarshalerFunc(lb.ACLActionType(""), human.EnumMarshalFunc(aclMarshalSpecs))
	human.RegisterMarshalerFunc(
		lb.BackendServerStatsHealthCheckStatus(""),
		human.EnumMarshalFunc(backendServerStatsHealthCheckStatusMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		lb.BackendServerStatsServerState(""),
		human.EnumMarshalFunc(backendServerStatsServerStateMarshalSpecs),
	)
	human.RegisterMarshalerFunc(lb.LB{}, lbMarshalerFunc)
	human.RegisterMarshalerFunc(lb.Backend{}, lbBackendMarshalerFunc)
	human.RegisterMarshalerFunc(lb.Frontend{}, lbFrontendMarshalerFunc)
	human.RegisterMarshalerFunc(lb.Certificate{}, lbCertificateMarshalerFunc)
	human.RegisterMarshalerFunc(lb.ACL{}, lbACLMarshalerFunc)
	human.RegisterMarshalerFunc([]*lb.PrivateNetwork{}, lbPrivateNetworksMarshalerFunc)

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

	cmds.MustFind("lb", "ip", "list").
		Override(core.ListViewBuilder("ID", "IPAddress", "OrganizationID", "ProjectID", "LBID", "Reverse", "Tags", "Region", "Zone"))
	cmds.MustFind("lb", "lb", "list").
		Override(core.ListViewBuilder("ID", "Name", "Description", "Status", "Instances", "OrganizationID", "ProjectID", "IP", "Tags", "FrontendCount", "BackendCount", "Type", "SslCompatibilityLevel", "CreatedAt", "UpdatedAt", "PrivateNetworkCount", "RouteCount", "Region", "Zone"))
	cmds.MustFind("lb", "lb-types", "list").
		Override(core.ListViewBuilder("Name", "StockStatus", "Bandwidth", "Multicloud", "Description", "Region", "Zone"))
	cmds.MustFind("lb", "route", "list").
		Override(core.ListViewBuilder("ID", "FrontendID", "BackendID", "CreatedAt", "UpdatedAt", "Zone"))

	cmds.MustFind("lb", "acl", "list").
		Override(core.ListViewBuilder("ID", "Name", "Frontend", "Index", "CreatedAt", "UpdatedAt", "Description", "Zone"))
	cmds.MustFind("lb", "backend", "list").
		Override(core.ListViewBuilder("ID", "Name", "ForwardProtocol", "ForwardPort", "ForwardPortAlgorithm", "StickySessions", "StickySessionsCookieName", "Pool", "LB", "SendProxyV2", "TimeoutServer", "TimeoutConnect", "TimeoutTunnel", "OnMarkedDownAction", "ProxyProtocol", "CreatedAt", "UpdatedAt", "FailoverHost", "SslBridging", "IgnoreSslServerVerify", "RedispatchAttemptCount", "MaxRetries", "MaxConnections", "TimeoutQueue", "Host", "Zone"))
	cmds.MustFind("lb", "frontend", "list").
		Override(core.ListViewBuilder("ID", "Name", "InboundPort", "Backend", "LB", "TimeoutClient", "Certificate", "CertificateIDs", "CreatedAt", "UpdatedAt", "EnableHTTP3", "ConnectionRateLimit", "EnableAccessLogs", "Zone"))
	cmds.MustFind("lb", "subscriber", "list").Override(core.ListViewBuilder("ID", "Name", "Zone"))

	return cmds
}
