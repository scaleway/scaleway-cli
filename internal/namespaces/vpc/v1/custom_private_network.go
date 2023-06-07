package vpc

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/api/redis/v1"
	"github.com/scaleway/scaleway-sdk-go/api/vpc/v1"
	"github.com/scaleway/scaleway-sdk-go/api/vpcgw/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func privateNetworkGetBuilder(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		getPNResp, err := runner(ctx, argsI)
		if err != nil {
			return getPNResp, err
		}
		pn := getPNResp.(*vpc.PrivateNetwork)

		client := core.ExtractClient(ctx)

		// Instance
		listInstanceServers, err := listCustomInstanceServers(client, pn)
		if err != nil {
			return nil, err
		}

		// Baremetal
		listBaremetalServers, err := listCustomBaremetalServers(client, pn)
		if err != nil {
			return nil, err
		}

		// LB
		listLBs, err := listCustomLBs(client, pn)
		if err != nil {
			return nil, err
		}

		// Rdb
		listRdbInstances, err := listCustomRdbs(client, pn)
		if err != nil {
			return nil, err
		}

		// Redis
		listRedisClusters, err := listCustomRedisClusters(client, pn)
		if err != nil {
			return nil, err
		}

		// Gateway
		listGateways, err := listCustomGateways(client, pn)
		if err != nil {
			return nil, err
		}

		return &struct {
			*vpc.PrivateNetwork
			InstanceServers  []customInstanceServer  `json:"instance_servers,omitempty"`
			BaremetalServers []customBaremetalServer `json:"baremetal_servers,omitempty"`
			LBs              []customLB              `json:"lbs,omitempty"`
			RdbInstances     []customRdb             `json:"rdb_instances,omitempty"`
			RedisClusters    []customRedis           `json:"redis_clusters,omitempty"`
			Gateways         []customGateway         `json:"gateways,omitempty"`
		}{
			pn,
			listInstanceServers,
			listBaremetalServers,
			listLBs,
			listRdbInstances,
			listRedisClusters,
			listGateways,
		}, nil
	}

	c.View = &core.View{
		Sections: []*core.ViewSection{
			{
				FieldName:   "InstanceServers",
				Title:       "Instance Servers",
				HideIfEmpty: true,
			},
			{
				FieldName:   "BaremetalServers",
				Title:       "Baremetal Servers",
				HideIfEmpty: true,
			},
			{
				FieldName:   "LBs",
				Title:       "Load-Balancers",
				HideIfEmpty: true,
			},
			{
				FieldName:   "RdbInstances",
				Title:       "Rdb Instances",
				HideIfEmpty: true,
			},
			{
				FieldName:   "RedisClusters",
				Title:       "Redis Clusters",
				HideIfEmpty: true,
			},
			{
				FieldName:   "Gateways",
				Title:       "Public Gateways",
				HideIfEmpty: true,
			},
		},
	}

	return c
}

type customInstanceServer struct {
	ID         string               `json:"id"`
	Name       string               `json:"name"`
	State      instance.ServerState `json:"state"`
	NicID      string               `json:"nic_id"`
	MacAddress string               `json:"mac"`
}
type customBaremetalServer struct {
	ID                 string                 `json:"id"`
	Name               string                 `json:"server_id"`
	State              baremetal.ServerStatus `json:"state"`
	BaremetalNetworkID string                 `json:"baremetal_network_id"`
	Vlan               *uint32                `json:"vlan"`
}
type customLB struct {
	ID    string      `json:"id"`
	Name  string      `json:"server_id"`
	State lb.LBStatus `json:"state"`
}
type customRdb struct {
	ID         string             `json:"id"`
	Name       string             `json:"server_id"`
	State      rdb.InstanceStatus `json:"state"`
	EndpointID string             `json:"endpoint_id"`
}
type customRedis struct {
	ID         string              `json:"id"`
	Name       string              `json:"server_id"`
	State      redis.ClusterStatus `json:"state"`
	EndpointID string              `json:"endpoint_id"`
}
type customGateway struct {
	ID               string              `json:"id"`
	Name             string              `json:"server_id"`
	State            vpcgw.GatewayStatus `json:"state"`
	GatewayNetworkID string              `json:"gateway_network_id"`
}

func listCustomInstanceServers(client *scw.Client, pn *vpc.PrivateNetwork) ([]customInstanceServer, error) {
	instanceAPI := instance.NewAPI(client)
	listInstanceServers, err := instanceAPI.ListServers(&instance.ListServersRequest{
		PrivateNetwork: &pn.ID,
	}, scw.WithAllPages())
	if err != nil {
		return nil, err
	}
	var customInstanceServers []customInstanceServer
	for _, server := range listInstanceServers.Servers {
		for _, nic := range server.PrivateNics {
			if nic.PrivateNetworkID == pn.ID {
				customInstanceServers = append(customInstanceServers, customInstanceServer{
					NicID:      nic.ID,
					ID:         nic.ServerID,
					MacAddress: nic.MacAddress,
					Name:       server.Name,
					State:      server.State,
				})
			}
		}
	}
	return customInstanceServers, nil
}

func listCustomBaremetalServers(client *scw.Client, pn *vpc.PrivateNetwork) ([]customBaremetalServer, error) {
	baremtalPNAPI := baremetal.NewPrivateNetworkAPI(client)
	baremetalAPI := baremetal.NewAPI(client)
	listBaremetalServers, err := baremtalPNAPI.ListServerPrivateNetworks(&baremetal.PrivateNetworkAPIListServerPrivateNetworksRequest{
		Zone:             pn.Zone,
		PrivateNetworkID: &pn.ID,
	}, scw.WithAllPages())
	if err != nil {
		return nil, err
	}
	var customBaremetalServers []customBaremetalServer
	for _, server := range listBaremetalServers.ServerPrivateNetworks {
		if server.PrivateNetworkID == pn.ID {
			getBaremetalServer, err := baremetalAPI.GetServer(&baremetal.GetServerRequest{
				Zone:     pn.Zone,
				ServerID: server.ServerID,
			})
			if err != nil {
				return nil, err
			}
			customBaremetalServers = append(customBaremetalServers, customBaremetalServer{
				ID:                 server.ServerID,
				State:              getBaremetalServer.Status,
				BaremetalNetworkID: server.ID,
				Name:               getBaremetalServer.Name,
				Vlan:               server.Vlan,
			})
		}
	}
	return customBaremetalServers, nil
}

func listCustomLBs(client *scw.Client, pn *vpc.PrivateNetwork) ([]customLB, error) {
	LBAPI := lb.NewZonedAPI(client)
	listLbs, err := LBAPI.ListLBs(&lb.ZonedAPIListLBsRequest{
		Zone: pn.Zone,
	})
	var filteredLBs []*lb.LB
	for _, loadbalancer := range listLbs.LBs {
		if loadbalancer.PrivateNetworkCount >= 1 {
			filteredLBs = append(filteredLBs, loadbalancer)
		}
	}
	if err != nil {
		return nil, err
	}

	var customLBs []customLB
	for _, loadbalancer := range filteredLBs {
		listLBpns, err := LBAPI.ListLBPrivateNetworks(&lb.ZonedAPIListLBPrivateNetworksRequest{
			Zone: loadbalancer.Zone,
			LBID: loadbalancer.ID,
		}, scw.WithAllPages())
		if err != nil {
			return nil, err
		}
		for _, res := range listLBpns.PrivateNetwork {
			if res.PrivateNetworkID == pn.ID {
				customLBs = append(customLBs, customLB{
					ID:    res.LB.ID,
					Name:  res.LB.Name,
					State: res.LB.Status,
				})
			}
		}
	}
	return customLBs, nil
}

func listCustomRdbs(client *scw.Client, pn *vpc.PrivateNetwork) ([]customRdb, error) {
	rdbAPI := rdb.NewAPI(client)
	region, err := scw.Zone.Region(pn.Zone)
	if err != nil {
		return nil, err
	}
	listDBs, err := rdbAPI.ListInstances(&rdb.ListInstancesRequest{
		Region: region,
	}, scw.WithAllPages())
	if err != nil {
		return nil, err
	}
	var customRdbs []customRdb
	for _, db := range listDBs.Instances {
		for _, endpoint := range db.Endpoints {
			if endpoint.PrivateNetwork != nil && endpoint.PrivateNetwork.PrivateNetworkID == pn.ID {
				customRdbs = append(customRdbs, customRdb{
					EndpointID: endpoint.ID,
					ID:         db.ID,
					Name:       db.Name,
					State:      db.Status,
				})
			}
		}
	}
	return customRdbs, nil
}

func listCustomRedisClusters(client *scw.Client, pn *vpc.PrivateNetwork) ([]customRedis, error) {
	redisAPI := redis.NewAPI(client)
	listRedisClusters, err := redisAPI.ListClusters(&redis.ListClustersRequest{
		Zone: pn.Zone,
	}, scw.WithAllPages())
	if err != nil {
		return nil, err
	}
	var customClusters []customRedis
	for _, cluster := range listRedisClusters.Clusters {
		for _, endpoint := range cluster.Endpoints {
			if endpoint.PrivateNetwork != nil && endpoint.PrivateNetwork.ID == pn.ID {
				customClusters = append(customClusters, customRedis{
					ID:         cluster.ID,
					Name:       cluster.Name,
					State:      cluster.Status,
					EndpointID: endpoint.ID,
				})
			}
		}
	}
	return customClusters, nil
}

func listCustomGateways(client *scw.Client, pn *vpc.PrivateNetwork) ([]customGateway, error) {
	vpcgwAPI := vpcgw.NewAPI(client)
	listGateways, err := vpcgwAPI.ListGateways(&vpcgw.ListGatewaysRequest{
		Zone: pn.Zone,
	}, scw.WithAllPages())
	if err != nil {
		return nil, err
	}
	var customGateways []customGateway
	for _, gateway := range listGateways.Gateways {
		for _, gatewayNetwork := range gateway.GatewayNetworks {
			if gatewayNetwork.PrivateNetworkID == pn.ID {
				customGateways = append(customGateways, customGateway{
					ID:               gateway.ID,
					Name:             gateway.Name,
					State:            gateway.Status,
					GatewayNetworkID: gatewayNetwork.GatewayID,
				})
			}
		}
	}
	return customGateways, nil
}
