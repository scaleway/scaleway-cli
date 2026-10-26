package connect

import (
	"fmt"
	"strings"

	vpcgw "github.com/scaleway/scaleway-sdk-go/api/vpcgw/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// PickGatewayForPrivateNetwork returns the VPC gateway whose SSH bastion can
// be used to reach the given private network: it must be attached to the
// network, running, have a public IPv4 address and the SSH bastion enabled.
//
// It is shared by the connect commands (redis and, soon, the other managed
// databases) that route a private endpoint through a VPC gateway. When
// several gateways match, the caller must pick one explicitly.
//
// ponytail: only gateways attached to the network are considered; a gateway
// that reaches it via default route advertisement is not detected. Upgrade
// path: match the advertised routes as well.
func PickGatewayForPrivateNetwork(
	client *scw.Client,
	zone scw.Zone,
	privateNetworkID string,
) (*vpcgw.Gateway, error) {
	api := vpcgw.NewAPI(client)

	req := &vpcgw.ListGatewaysRequest{Zone: zone}
	if projectID, ok := client.GetDefaultProjectID(); ok {
		req.ProjectID = &projectID
	}

	gateways, err := api.ListGateways(req, scw.WithAllPages())
	if err != nil {
		return nil, fmt.Errorf("failed to list VPC gateways: %w", err)
	}

	candidates := make([]*vpcgw.Gateway, 0, len(gateways.Gateways))
	for _, gw := range gateways.Gateways {
		if gatewayReachesPrivateNetwork(gw, privateNetworkID) {
			candidates = append(candidates, gw)
		}
	}

	switch len(candidates) {
	case 0:
		return nil, fmt.Errorf(
			"no VPC gateway with the SSH bastion enabled is attached to private network %s in %s",
			privateNetworkID,
			zone,
		)
	case 1:
		return candidates[0], nil
	default:
		ids := make([]string, 0, len(candidates))
		for _, gw := range candidates {
			ids = append(ids, gw.ID)
		}

		return nil, fmt.Errorf(
			"multiple VPC gateways with the SSH bastion enabled are attached to private network %s (%s): specify one with vpc-gateway-id",
			privateNetworkID,
			strings.Join(ids, ", "),
		)
	}
}

func gatewayReachesPrivateNetwork(gw *vpcgw.Gateway, privateNetworkID string) bool {
	if !gw.BastionEnabled || gw.IPv4 == nil || gw.Status != vpcgw.GatewayStatusRunning {
		return false
	}

	for _, network := range gw.GatewayNetworks {
		if network.PrivateNetworkID == privateNetworkID {
			return true
		}
	}

	return false
}
