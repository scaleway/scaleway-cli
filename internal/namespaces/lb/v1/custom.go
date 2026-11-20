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
	humanv2.RegisterEnumFormat(lbTypeStockFormat)
	humanv2.RegisterEnumFormat(lbStatusFormat)
	humanv2.RegisterEnumFormat(certificateStatusFormat)
	humanv2.RegisterEnumFormat(aclFormat)
	humanv2.RegisterEnumFormat(backendServerStatsHealthCheckStatusFormat)
	humanv2.RegisterEnumFormat(backendServerStatsServerStateFormat)
	humanv2.RegisterLayout(reflect.TypeFor[lb.LB](), lbLayout)
	humanv2.RegisterLayout(reflect.TypeFor[lb.Backend](), lbBackendLayout)
	humanv2.RegisterLayout(reflect.TypeFor[lb.Frontend](), lbFrontendLayout)
	humanv2.RegisterLayout(reflect.TypeFor[lb.Certificate](), lbCertificateLayout)
	humanv2.RegisterLayout(reflect.TypeFor[lb.ACL](), lbACLLayout)
	humanv2.RegisterLayout(reflect.TypeFor[[]*lb.PrivateNetwork](), lbPrivateNetworkListLayout)

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
