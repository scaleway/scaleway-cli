package k8s

import (
	"context"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type cacheCluster struct {
	ID        string
	Name      string
	Region    string
	ProjectID string
}

func (cacheCluster) TableName() string {
	return "k8s_cluster"
}

type cachePool struct {
	ID        string
	Name      string
	Region    string
	ProjectID string
	ClusterID string
}

func (cachePool) TableName() string {
	return "k8s_pool"
}

type cacheNode struct {
	ID        string
	Name      string
	Region    string
	ClusterID string
	ProjectID string
}

func (cacheNode) TableName() string {
	return "k8s_node"
}

func RefreshCache(ctx context.Context) (interface{}, error) {
	client := core.ExtractClient(ctx)
	api := k8s.NewAPI(client)
	database := core.ExtractCacheDB(ctx)

	database.AutoMigrate(&cacheCluster{})
	database.Unscoped().Delete(&cacheCluster{})

	database.AutoMigrate(&cacheNode{})
	database.Unscoped().Delete(&cacheNode{})

	database.AutoMigrate(&cachePool{})
	database.Unscoped().Delete(&cachePool{})

	for _, region := range []scw.Region{scw.RegionFrPar, scw.RegionNlAms} {
		listClusters, err := api.ListClusters(&k8s.ListClustersRequest{
			Region: region,
		})
		if err != nil {
			return nil, err
		}
		for _, cluster := range listClusters.Clusters {
			database.Create(&cacheCluster{
				ID:        cluster.ID,
				Name:      cluster.Name,
				Region:    cluster.Region.String(),
				ProjectID: cluster.OrganizationID,
			})

			listNodes, err := api.ListNodes(&k8s.ListNodesRequest{
				Region:    region,
				ClusterID: cluster.ID,
			}, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			for _, node := range listNodes.Nodes {
				database.Create(&cacheNode{
					ID:        node.ID,
					Name:      node.Name,
					Region:    node.Region.String(),
					ClusterID: node.ClusterID,
					ProjectID: cluster.OrganizationID,
				})
			}

			listPools, err := api.ListPools(&k8s.ListPoolsRequest{
				Region:    region,
				ClusterID: cluster.ID,
			}, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			for _, pool := range listPools.Pools {
				database.Create(&cachePool{
					ID:        pool.ID,
					Name:      pool.Name,
					Region:    pool.Region.String(),
					ClusterID: pool.ClusterID,
					ProjectID: cluster.OrganizationID,
				})
			}
		}
	}

	return nil, nil
}
