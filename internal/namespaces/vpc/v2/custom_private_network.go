package vpc

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	applesilicon "github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	inference "github.com/scaleway/scaleway-sdk-go/api/inference/v1"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/api/ipam/v1"
	"github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	mongodb "github.com/scaleway/scaleway-sdk-go/api/mongodb/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/api/redis/v1"
	"github.com/scaleway/scaleway-sdk-go/api/vpc/v2"
	"github.com/scaleway/scaleway-sdk-go/api/vpcgw/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"golang.org/x/sync/errgroup"
)

func privateNetworkMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type tmp vpc.PrivateNetwork
	pn := tmp(i.(vpc.PrivateNetwork))

	// Sections
	opt.Sections = []*human.MarshalSection{
		{
			FieldName:   "Subnets",
			Title:       "Subnets",
			HideIfEmpty: true,
		},
	}

	str, err := human.Marshal(pn, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

func privateNetworkGetBuilder(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		getPNResp, err := runner(ctx, argsI)
		if err != nil {
			return getPNResp, err
		}
		pn := getPNResp.(*vpc.PrivateNetwork)

		client := core.ExtractClient(ctx)

		var (
			instanceServers      []customInstanceServer
			baremetalServers     []customBaremetalServer
			k8sClusters          []customK8sCluster
			lbs                  []customLB
			rdbs                 []customRdb
			redisClusters        []customRedis
			gateways             []customGateway
			appleSiliconServers  []customAppleSiliconServer
			mongoDBs             []customMongoDB
			ipamIPs              []customIPAMIP
			inferenceDeployments []customInferenceDeployment
		)

		g, groupCtx := errgroup.WithContext(ctx)

		g.Go(func() (err error) {
			instanceServers, err = listCustomInstanceServers(groupCtx, client, pn)

			return
		})
		g.Go(func() (err error) {
			baremetalServers, err = listCustomBaremetalServers(groupCtx, client, pn)

			return
		})
		g.Go(func() (err error) {
			k8sClusters, err = listCustomK8sClusters(groupCtx, client, pn)

			return
		})
		g.Go(func() (err error) {
			lbs, err = listCustomLBs(groupCtx, client, pn)

			return
		})
		g.Go(func() (err error) {
			rdbs, err = listCustomRdbs(groupCtx, client, pn)

			return
		})
		g.Go(func() (err error) {
			redisClusters, err = listCustomRedisClusters(groupCtx, client, pn)

			return
		})
		g.Go(func() (err error) {
			gateways, err = listCustomGateways(groupCtx, client, pn)

			return
		})
		g.Go(func() (err error) {
			appleSiliconServers, err = listCustomAppleSiliconServers(groupCtx, client, pn)

			return
		})
		g.Go(func() (err error) {
			mongoDBs, err = listCustomMongoDBs(groupCtx, client, pn)

			return
		})
		g.Go(func() (err error) {
			ipamIPs, err = listCustomIPAMIPs(groupCtx, client, pn)

			return
		})
		g.Go(func() (err error) {
			inferenceDeployments, err = listCustomInferenceDeployments(groupCtx, client, pn)

			return
		})

		if err = g.Wait(); err != nil {
			return nil, err
		}

		return &struct {
			*vpc.PrivateNetwork
			InstanceServers      []customInstanceServer      `json:"instance_servers,omitempty"`
			BaremetalServers     []customBaremetalServer     `json:"baremetal_servers,omitempty"`
			K8sClusters          []customK8sCluster          `json:"K8s_clusters,omitempty"`
			LBs                  []customLB                  `json:"lbs,omitempty"`
			RdbInstances         []customRdb                 `json:"rdb_instances,omitempty"`
			RedisClusters        []customRedis               `json:"redis_clusters,omitempty"`
			Gateways             []customGateway             `json:"gateways,omitempty"`
			AppleSiliconServers  []customAppleSiliconServer  `json:"apple_silicon_servers,omitempty"`
			MongoDBInstances     []customMongoDB             `json:"mongodb_instances,omitempty"`
			IPAMIPs              []customIPAMIP              `json:"ipam_ips,omitempty"`
			InferenceDeployments []customInferenceDeployment `json:"inference_deployments,omitempty"`
		}{
			pn,
			instanceServers,
			baremetalServers,
			k8sClusters,
			lbs,
			rdbs,
			redisClusters,
			gateways,
			appleSiliconServers,
			mongoDBs,
			ipamIPs,
			inferenceDeployments,
		}, nil
	}

	c.View = &core.View{
		Sections: []*core.ViewSection{
			{
				FieldName:   "Subnets",
				Title:       "Subnets",
				HideIfEmpty: true,
			},
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
				FieldName:   "K8sClusters",
				Title:       "K8s Clusters",
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
			{
				FieldName:   "AppleSiliconServers",
				Title:       "AppleSilicon Servers",
				HideIfEmpty: true,
			},
			{
				FieldName:   "MongoDBInstances",
				Title:       "MongoDB Instances",
				HideIfEmpty: true,
			},
			{
				FieldName:   "IPAMIPs",
				Title:       "IPAM IPs",
				HideIfEmpty: true,
			},
			{
				FieldName:   "InferenceDeployments",
				Title:       "Inference Deployments",
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
	Name               string                 `json:"name"`
	State              baremetal.ServerStatus `json:"state"`
	BaremetalNetworkID string                 `json:"baremetal_network_id"`
	Vlan               *uint32                `json:"vlan"`
}
type customK8sCluster struct {
	ID    string            `json:"id"`
	Name  string            `json:"name"`
	State k8s.ClusterStatus `json:"state"`
}
type customLB struct {
	ID    string      `json:"id"`
	Name  string      `json:"name"`
	State lb.LBStatus `json:"state"`
}
type customRdb struct {
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	State      rdb.InstanceStatus `json:"state"`
	EndpointID string             `json:"endpoint_id"`
}
type customRedis struct {
	ID         string              `json:"id"`
	Name       string              `json:"name"`
	State      redis.ClusterStatus `json:"state"`
	EndpointID string              `json:"endpoint_id"`
}

type customMongoDB struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	State      mongodb.InstanceStatus `json:"state"`
	EndpointID string                 `json:"endpoint_id"`
}
type customGateway struct {
	ID               string              `json:"id"`
	Name             string              `json:"name"`
	State            vpcgw.GatewayStatus `json:"state"`
	GatewayNetworkID string              `json:"gateway_network_id"`
}

type customAppleSiliconServer struct {
	ID        string                    `json:"id"`
	Name      string                    `json:"name"`
	State     applesilicon.ServerStatus `json:"state"`
	Vlan      *uint32                   `json:"vlan"`
	MappingId string                    `json:"mapping_id"`
}

type customIPAMIP struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}

type customInferenceDeployment struct {
	ID         string                     `json:"id"`
	Name       string                     `json:"name"`
	State      inference.DeploymentStatus `json:"state"`
	EndpointID string                     `json:"endpoint_id"`
}

func listCustomInstanceServers(
	ctx context.Context,
	client *scw.Client,
	pn *vpc.PrivateNetwork,
) ([]customInstanceServer, error) {
	instanceAPI := instance.NewAPI(client)

	regionZones := pn.Region.GetZones()
	instanceZones := instanceAPI.Zones()
	zones := intersectZones(regionZones, instanceZones)

	var customInstanceServers []customInstanceServer
	for _, zone := range zones {
		listInstanceServers, err := instanceAPI.ListServers(&instance.ListServersRequest{
			PrivateNetwork: &pn.ID,
			Zone:           zone,
		}, scw.WithContext(ctx), scw.WithAllPages())
		if err != nil {
			return nil, err
		}
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
	}

	return customInstanceServers, nil
}

func listCustomBaremetalServers(
	ctx context.Context,
	client *scw.Client,
	pn *vpc.PrivateNetwork,
) ([]customBaremetalServer, error) {
	baremetalPNAPI := baremetal.NewPrivateNetworkAPI(client)
	baremetalAPI := baremetal.NewAPI(client)

	regionZones := pn.Region.GetZones()
	baremetalZones := baremetalAPI.Zones()
	zones := intersectZones(regionZones, baremetalZones)

	var customBaremetalServers []customBaremetalServer
	for _, zone := range zones {
		listBaremetalServers, err := baremetalPNAPI.ListServerPrivateNetworks(
			&baremetal.PrivateNetworkAPIListServerPrivateNetworksRequest{
				Zone:             zone,
				PrivateNetworkID: &pn.ID,
			}, scw.WithContext(ctx), scw.WithAllPages(),
		)
		if err != nil {
			return nil, err
		}

		for _, server := range listBaremetalServers.ServerPrivateNetworks {
			if server.PrivateNetworkID == pn.ID {
				getBaremetalServer, err := baremetalAPI.GetServer(&baremetal.GetServerRequest{
					Zone:     zone,
					ServerID: server.ServerID,
				}, scw.WithContext(ctx))
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
	}

	return customBaremetalServers, nil
}

func listCustomK8sClusters(
	ctx context.Context,
	client *scw.Client,
	pn *vpc.PrivateNetwork,
) ([]customK8sCluster, error) {
	k8sAPI := k8s.NewAPI(client)

	listK8sClusters, err := k8sAPI.ListClusters(&k8s.ListClustersRequest{
		Region: pn.Region,
	}, scw.WithContext(ctx), scw.WithAllPages())
	if err != nil {
		return nil, err
	}
	var customK8sClusters []customK8sCluster
	for _, cluster := range listK8sClusters.Clusters {
		if cluster.PrivateNetworkID != nil && *cluster.PrivateNetworkID == pn.ID {
			customK8sClusters = append(customK8sClusters, customK8sCluster{
				ID:    cluster.ID,
				Name:  cluster.Name,
				State: cluster.Status,
			})
		}
	}

	return customK8sClusters, nil
}

func listCustomLBs(
	ctx context.Context,
	client *scw.Client,
	pn *vpc.PrivateNetwork,
) ([]customLB, error) {
	LBAPI := lb.NewZonedAPI(client)

	regionZones := pn.Region.GetZones()
	lbZones := LBAPI.Zones()
	zones := intersectZones(regionZones, lbZones)

	var customLBs []customLB
	for _, zone := range zones {
		listLbs, err := LBAPI.ListLBs(&lb.ZonedAPIListLBsRequest{
			Zone: zone,
		}, scw.WithContext(ctx), scw.WithAllPages())
		if err != nil {
			return nil, err
		}

		var filteredLBs []*lb.LB
		for _, loadbalancer := range listLbs.LBs {
			if loadbalancer.PrivateNetworkCount >= 1 {
				filteredLBs = append(filteredLBs, loadbalancer)
			}
		}

		for _, loadbalancer := range filteredLBs {
			listLBpns, err := LBAPI.ListLBPrivateNetworks(&lb.ZonedAPIListLBPrivateNetworksRequest{
				Zone: zone,
				LBID: loadbalancer.ID,
			}, scw.WithContext(ctx), scw.WithAllPages())
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
	}

	return customLBs, nil
}

func listCustomRdbs(
	ctx context.Context,
	client *scw.Client,
	pn *vpc.PrivateNetwork,
) ([]customRdb, error) {
	rdbAPI := rdb.NewAPI(client)

	listDBs, err := rdbAPI.ListInstances(&rdb.ListInstancesRequest{
		Region: pn.Region,
	}, scw.WithContext(ctx), scw.WithAllPages())
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

func listCustomRedisClusters(
	ctx context.Context,
	client *scw.Client,
	pn *vpc.PrivateNetwork,
) ([]customRedis, error) {
	redisAPI := redis.NewAPI(client)

	regionZones := pn.Region.GetZones()
	redisZones := redisAPI.Zones()
	zones := intersectZones(regionZones, redisZones)

	var customClusters []customRedis
	for _, zone := range zones {
		listRedisClusters, err := redisAPI.ListClusters(&redis.ListClustersRequest{
			Zone: zone,
		}, scw.WithContext(ctx), scw.WithAllPages())
		if err != nil {
			return nil, err
		}
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
	}

	return customClusters, nil
}

func listCustomGateways(
	ctx context.Context,
	client *scw.Client,
	pn *vpc.PrivateNetwork,
) ([]customGateway, error) {
	vpcgwAPI := vpcgw.NewAPI(client)

	regionZones := pn.Region.GetZones()
	vpcgwZones := vpcgwAPI.Zones()
	zones := intersectZones(regionZones, vpcgwZones)

	var customGateways []customGateway
	for _, zone := range zones {
		listGateways, err := vpcgwAPI.ListGateways(&vpcgw.ListGatewaysRequest{
			Zone: zone,
		}, scw.WithContext(ctx), scw.WithAllPages())
		if err != nil {
			return nil, err
		}
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
	}

	return customGateways, nil
}

func listCustomAppleSiliconServers(
	ctx context.Context,
	client *scw.Client,
	pn *vpc.PrivateNetwork,
) ([]customAppleSiliconServer, error) {
	appleSiliconAPI := applesilicon.NewAPI(client)
	appleSiliconPrivateNetworkAPI := applesilicon.NewPrivateNetworkAPI(client)

	regionZones := pn.Region.GetZones()
	appleSiliconZones := appleSiliconAPI.Zones()
	zones := intersectZones(regionZones, appleSiliconZones)

	var customAppleSiliconServers []customAppleSiliconServer

	for _, zone := range zones {
		listAppleSiliconServers, err := appleSiliconPrivateNetworkAPI.ListServerPrivateNetworks(
			&applesilicon.PrivateNetworkAPIListServerPrivateNetworksRequest{
				Zone:             zone,
				PrivateNetworkID: &pn.ID,
			}, scw.WithContext(ctx), scw.WithAllPages(),
		)
		if err != nil {
			return nil, err
		}

		for _, server := range listAppleSiliconServers.ServerPrivateNetworks {
			if server.PrivateNetworkID == pn.ID {
				getAppleSiliconServer, err := appleSiliconAPI.GetServer(
					&applesilicon.GetServerRequest{
						Zone:     zone,
						ServerID: server.ServerID,
					}, scw.WithContext(ctx),
				)
				if err != nil {
					return nil, err
				}

				customAppleSiliconServers = append(
					customAppleSiliconServers,
					customAppleSiliconServer{
						ID:        getAppleSiliconServer.ID,
						Name:      getAppleSiliconServer.Name,
						State:     getAppleSiliconServer.Status,
						Vlan:      server.Vlan,
						MappingId: server.ID,
					},
				)
			}
		}
	}

	return customAppleSiliconServers, nil
}

func listCustomMongoDBs(
	ctx context.Context,
	client *scw.Client,
	pn *vpc.PrivateNetwork,
) ([]customMongoDB, error) {
	mongoAPI := mongodb.NewAPI(client)

	listDBs, err := mongoAPI.ListInstances(&mongodb.ListInstancesRequest{
		Region: pn.Region,
	}, scw.WithContext(ctx), scw.WithAllPages())
	if err != nil {
		return nil, err
	}
	var customDBs []customMongoDB
	for _, db := range listDBs.Instances {
		for _, endpoint := range db.Endpoints {
			if endpoint.PrivateNetwork != nil && endpoint.PrivateNetwork.PrivateNetworkID == pn.ID {
				customDBs = append(customDBs, customMongoDB{
					EndpointID: endpoint.ID,
					ID:         db.ID,
					Name:       db.Name,
					State:      db.Status,
				})
			}
		}
	}

	return customDBs, nil
}

func listCustomIPAMIPs(
	ctx context.Context,
	client *scw.Client,
	pn *vpc.PrivateNetwork,
) ([]customIPAMIP, error) {
	ipamAPI := ipam.NewAPI(client)

	listIPAMIPs, err := ipamAPI.ListIPs(&ipam.ListIPsRequest{
		Region:           pn.Region,
		PrivateNetworkID: &pn.ID,
	}, scw.WithContext(ctx), scw.WithAllPages())
	if err != nil {
		return nil, err
	}

	customIPAMIPs := make([]customIPAMIP, 0, len(listIPAMIPs.IPs))
	for _, ip := range listIPAMIPs.IPs {
		customIPAMIPs = append(customIPAMIPs, customIPAMIP{
			ID:      ip.ID,
			Address: ip.Address.String(),
		})
	}

	return customIPAMIPs, nil
}

func listCustomInferenceDeployments(
	ctx context.Context,
	client *scw.Client,
	pn *vpc.PrivateNetwork,
) ([]customInferenceDeployment, error) {
	inferenceAPI := inference.NewAPI(client)

	listDeployments, err := inferenceAPI.ListDeployments(&inference.ListDeploymentsRequest{
		Region: pn.Region,
	}, scw.WithContext(ctx), scw.WithAllPages())
	if err != nil {
		return nil, err
	}
	var customDeployments []customInferenceDeployment
	for _, deployment := range listDeployments.Deployments {
		for _, endpoint := range deployment.Endpoints {
			if endpoint.PrivateNetwork != nil && endpoint.PrivateNetwork.PrivateNetworkID == pn.ID {
				customDeployments = append(customDeployments, customInferenceDeployment{
					EndpointID: endpoint.ID,
					ID:         deployment.ID,
					Name:       deployment.Name,
					State:      deployment.Status,
				})
			}
		}
	}

	return customDeployments, nil
}

// intersectZones returns zones common to both provided slices
func intersectZones(regionZones, apiZones []scw.Zone) []scw.Zone {
	apiZonesMap := make(map[scw.Zone]bool)
	for _, zone := range apiZones {
		apiZonesMap[zone] = true
	}

	var intersect []scw.Zone
	for _, zone := range regionZones {
		if _, ok := apiZonesMap[zone]; ok {
			intersect = append(intersect, zone)
		}
	}

	return intersect
}

func privateNetworkDeleteBuilder(c *core.Command) *core.Command {
	c.Run = func(ctx context.Context, args any) (i any, e error) {
		request := args.(*vpc.DeletePrivateNetworkRequest)

		client := core.ExtractClient(ctx)
		api := vpc.NewAPI(client)

		return tryDeletingPrivateNetwork(ctx, api, request.Region, request.PrivateNetworkID, 5)
	}

	return c
}

func tryDeletingPrivateNetwork(
	ctx context.Context,
	api *vpc.API,
	region scw.Region,
	pnID string,
	retriesLeft int,
) (*vpc.PrivateNetwork, error) {
	err := api.DeletePrivateNetwork(&vpc.DeletePrivateNetworkRequest{
		PrivateNetworkID: pnID,
		Region:           region,
	}, scw.WithContext(ctx))

	var respErr *scw.ResponseError
	if errors.As(err, &respErr) && respErr.StatusCode == http.StatusInternalServerError {
		time.Sleep(time.Second * 5)
		if retriesLeft > 0 {
			return tryDeletingPrivateNetwork(ctx, api, region, pnID, retriesLeft-1)
		}
	}

	return nil, err
}
