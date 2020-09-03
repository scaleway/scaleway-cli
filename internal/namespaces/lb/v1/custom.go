package lb

import (
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
)

func GetCommands() *core.Commands {
	human.RegisterMarshalerFunc(lb.LBTypeStock(""), human.EnumMarshalFunc(lbTypeStockMarshalSpecs))
	human.RegisterMarshalerFunc(lb.LBStatus(""), human.EnumMarshalFunc(lbStatusMarshalSpecs))
	human.RegisterMarshalerFunc(lb.CertificateStatus(""), human.EnumMarshalFunc(certificateStatusMarshalSpecs))
	human.RegisterMarshalerFunc(lb.ACLActionType(""), human.EnumMarshalFunc(aclMarshalSpecs))
	human.RegisterMarshalerFunc(lb.BackendServerStatsHealthCheckStatus(""), human.EnumMarshalFunc(backendServerStatsHealthCheckStatusMarshalSpecs))
	human.RegisterMarshalerFunc(lb.BackendServerStatsServerState(""), human.EnumMarshalFunc(backendServerStatsServerStateMarshalSpecs))

	cmds := GetGeneratedCommands()

	cmds.Add(
		lbWaitCommand(),
	)

	cmds.MustFind("lb", "lb", "create").Override(lbCreateBuilder)
	cmds.MustFind("lb", "lb", "get").Override(lbGetBuilder)
	cmds.MustFind("lb", "lb", "migrate").Override(lbMigrateBuilder)
	cmds.MustFind("lb", "lb", "get-stats").Override(lbGetStatsBuilder)

	cmds.MustFind("lb", "frontend", "get").Override(frontendGetBuilder)

	cmds.MustFind("lb", "certificate", "create").Override(certificateCreateBuilder)

	return cmds
}
