package registry

import (
	"context"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/registry/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type cacheNamespace struct {
	ID        string
	Name      string
	Region    string
	ProjectID string
}

func (cacheNamespace) TableName() string {
	return "registry_namespace"
}

func RefreshCache(ctx context.Context) (interface{}, error) {
	client := core.ExtractClient(ctx)
	api := registry.NewAPI(client)
	database := core.ExtractCacheDB(ctx)

	database.AutoMigrate(&cacheNamespace{})
	database.Unscoped().Delete(&cacheNamespace{})

	for _, region := range []scw.Region{scw.RegionFrPar, scw.RegionNlAms} {
		listNamespaces, err := api.ListNamespaces(&registry.ListNamespacesRequest{
			Region: region,
		}, scw.WithAllPages())
		if err != nil {
			return nil, err
		}
		for _, namespace := range listNamespaces.Namespaces {
			database.Create(&cacheNamespace{
				ID:        namespace.ID,
				Name:      namespace.Name,
				Region:    namespace.Region.String(),
				ProjectID: namespace.OrganizationID,
			})
		}
	}

	return nil, nil
}
