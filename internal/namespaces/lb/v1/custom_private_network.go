package lb

import (
	"errors"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
)

func lbPrivateNetworksMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	privateNetworks, ok := i.([]*lb.PrivateNetwork)
	if !ok {
		return "", errors.New("invalid type: expected []*lb.PrivateNetwork")
	}

	type customPrivateNetwork struct {
		IpamIDs               []string                `json:"ipam_ids,omitempty"`
		DHCPConfigIPID        *string                 `json:"dhcp_config_ip_id,omitempty"`
		StaticConfigIPAddress *[]string               `json:"static_config_ip_address,omitempty"`
		PrivateNetworkID      string                  `json:"private_network_id"`
		Status                lb.PrivateNetworkStatus `json:"status"`
		CreatedAt             *time.Time              `json:"created_at"`
		UpdatedAt             *time.Time              `json:"updated_at"`
	}

	customPrivateNetworks := make([]customPrivateNetwork, 0, len(privateNetworks))
	for _, pn := range privateNetworks {
		if pn == nil {
			continue
		}

		customPN := customPrivateNetwork{
			IpamIDs:          pn.IpamIDs,
			PrivateNetworkID: pn.PrivateNetworkID,
			Status:           pn.Status,
			CreatedAt:        pn.CreatedAt,
			UpdatedAt:        pn.UpdatedAt,
		}

		//nolint: staticcheck
		if pn.DHCPConfig != nil {
			customPN.DHCPConfigIPID = pn.DHCPConfig.IPID
		}

		//nolint: staticcheck
		if pn.StaticConfig != nil {
			customPN.StaticConfigIPAddress = pn.StaticConfig.IPAddress
		}

		customPrivateNetworks = append(customPrivateNetworks, customPN)
	}

	return human.Marshal(customPrivateNetworks, opt)
}
