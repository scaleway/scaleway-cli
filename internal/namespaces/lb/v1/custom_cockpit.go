package lb

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/cockpit/v1"
	lb "github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type loadBalancerCockpitArgs struct {
	LBID string
	Zone scw.Zone
}

func loadBalancerCockpitCommand() *core.Command {
	return &core.Command{
		Namespace: "lb",
		Resource:  "load-balancer",
		Verb:      "cockpit",
		Short:     `Open the Cockpit Grafana dashboard for this load-balancer`,
		ArgsType:  reflect.TypeFor[loadBalancerCockpitArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "lb-id", Short: "Load Balancer ID", Required: true, Positional: true},
			core.ZoneArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*loadBalancerCockpitArgs)

			api := lb.NewZonedAPI(core.ExtractClient(ctx))
			loadBalancer, err := api.GetLB(&lb.ZonedAPIGetLBRequest{
				Zone: args.Zone,
				LBID: args.LBID,
			})
			if err != nil {
				return nil, err
			}

			region, _ := loadBalancer.Zone.Region()

			return cockpit.OpenDashboard(ctx, "scw-lb-overview", map[string]string{
				"region":  string(region),
				"lb_name": loadBalancer.Name,
			})
		},
	}
}
