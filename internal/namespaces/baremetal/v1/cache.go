package baremetal

import (
	"context"
	"strings"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type cacheServer struct {
	ID        string
	Name      string
	Zone      string
	ProjectID string
}

func (cacheServer) TableName() string {
	return "baremetal_server"
}

type cacheOs struct {
	UUID string
	Name string
	Zone string
}

func (cacheOs) TableName() string {
	return "baremetal_os"
}

func RefreshCache(ctx context.Context) (interface{}, error) {
	client := core.ExtractClient(ctx)
	api := baremetal.NewAPI(client)
	database := core.ExtractCacheDB(ctx)

	database.AutoMigrate(&cacheServer{})
	database.Unscoped().Delete(&cacheServer{})

	database.AutoMigrate(&cacheOs{})
	database.Unscoped().Delete(&cacheOs{})

	for _, zone := range []scw.Zone{scw.ZoneFrPar2} {
		listServers, err := api.ListServers(&baremetal.ListServersRequest{
			Zone: zone,
		}, scw.WithAllPages())
		if err != nil {
			return nil, err
		}
		for _, server := range listServers.Servers {
			database.Create(&cacheServer{
				ID:        server.ID,
				Name:      server.Name,
				Zone:      server.Zone.String(),
				ProjectID: server.OrganizationID,
			})
		}

		listOs, err := api.ListOS(&baremetal.ListOSRequest{
			Zone: zone,
		}, scw.WithAllPages())
		if err != nil {
			return nil, err
		}
		for _, os := range listOs.Os {
			database.Create(&cacheOs{
				UUID: os.ID,
				Name: strings.Join([]string{os.Name, os.Version}, " "),
				Zone: zone.String(),
			})
		}
	}

	return nil, nil
}
