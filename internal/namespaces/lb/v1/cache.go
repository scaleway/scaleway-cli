package lb

import (
	"context"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type cacheLB struct {
	UUID      string
	Name      string
	Region    string
	ProjectID string
}

func (cacheLB) TableName() string {
	return "lb_lb"
}

type cacheFrontend struct {
	UUID      string
	Name      string
	Region    string
	ProjectID string
}

func (cacheFrontend) TableName() string {
	return "lb_frontend"
}

type cacheBackend struct {
	UUID      string
	Name      string
	Region    string
	ProjectID string
}

func (cacheBackend) TableName() string {
	return "lb_backend"
}

func RefreshCache(ctx context.Context) (interface{}, error) {
	client := core.ExtractClient(ctx)
	api := lb.NewAPI(client)
	log := core.ExtractLogger(ctx)
	database := core.ExtractCacheDB(ctx)

	database.AutoMigrate(&cacheLB{})
	database.Unscoped().Delete(&cacheLB{})

	database.AutoMigrate(&cacheFrontend{})
	database.Unscoped().Delete(&cacheFrontend{})

	database.AutoMigrate(&cacheBackend{})
	database.Unscoped().Delete(&cacheBackend{})

	for _, region := range []scw.Region{scw.RegionFrPar, scw.RegionNlAms} {
		lbList, err := api.ListLBs(&lb.ListLBsRequest{
			Region: region,
		})
		if err != nil {
			return nil, err
		}

		for _, loadbalancer := range lbList.LBs {
			database.Create(&cacheLB{
				UUID:      loadbalancer.ID,
				Name:      loadbalancer.Name,
				Region:    loadbalancer.Region.String(),
				ProjectID: loadbalancer.OrganizationID,
			})

			listFrontends, err := api.ListFrontends(&lb.ListFrontendsRequest{
				Region: region,
				LBID:   loadbalancer.ID,
			}, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			for _, frontend := range listFrontends.Frontends {
				database.Create(&cacheFrontend{
					UUID:      frontend.ID,
					Name:      frontend.Name,
					Region:    loadbalancer.Region.String(),
					ProjectID: loadbalancer.OrganizationID,
				})
			}

			listBackends, err := api.ListBackends(&lb.ListBackendsRequest{
				Region: region,
				LBID:   loadbalancer.ID,
			}, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			for _, backend := range listBackends.Backends {
				database.Create(&cacheBackend{
					UUID:      backend.ID,
					Name:      backend.Name,
					Region:    loadbalancer.Region.String(),
					ProjectID: loadbalancer.OrganizationID,
				})
			}
			log.Infof("Successfully built cache for load-balancer %s\n", loadbalancer.Name)
		}
	}

	return nil, nil
}
