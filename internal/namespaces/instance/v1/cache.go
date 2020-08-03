package instance

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type cacheImage struct {
	UUID      string
	Name      string
	Zone      string
	ProjectID string
}

func (cacheImage) TableName() string {
	return "instance_images"
}

type cacheIP struct {
	UUID      string
	Address   string
	Zone      string
	ProjectID string
}

func (cacheIP) TableName() string {
	return "instance_ip"
}

type cachePlacementGroup struct {
	UUID      string
	Name      string
	Zone      string
	ProjectID string
}

func (cachePlacementGroup) TableName() string {
	return "instance_placement_group"
}

type cacheSecurityGroup struct {
	UUID      string
	Name      string
	Zone      string
	ProjectID string
}

func (cacheSecurityGroup) TableName() string {
	return "instance_security_group"
}

type cacheSnapshot struct {
	UUID      string
	Name      string
	Zone      string
	ProjectID string
}

func (cacheSnapshot) TableName() string {
	return "instance_snapshot"
}

type cacheVolume struct {
	UUID      string
	Name      string
	Zone      string
	ProjectID string
}

func (cacheVolume) TableName() string {
	return "instance_volume"
}

type cacheServer struct {
	UUID      string
	Name      string
	Zone      string
	ProjectID string
}

func (cacheServer) TableName() string {
	return "instance_server"
}

func RefreshCache(ctx context.Context) (interface{}, error) {
	client := core.ExtractClient(ctx)
	api := instance.NewAPI(client)
	database := core.ExtractCacheDB(ctx)

	database.AutoMigrate(&cachePlacementGroup{})
	database.Unscoped().Delete(&cachePlacementGroup{})

	database.AutoMigrate(&cacheSecurityGroup{})
	database.Unscoped().Delete(&cacheSecurityGroup{})

	database.AutoMigrate(&cacheServer{})
	database.Unscoped().Delete(&cacheServer{})

	database.AutoMigrate(&cacheSnapshot{})
	database.Unscoped().Delete(&cacheSnapshot{})

	database.AutoMigrate(&cacheImage{})
	database.Unscoped().Delete(&cacheImage{})

	database.AutoMigrate(&cacheIP{})
	database.Unscoped().Delete(&cacheIP{})

	database.AutoMigrate(&cacheVolume{})
	database.Unscoped().Delete(&cacheVolume{})

	for _, zone := range []scw.Zone{scw.ZoneFrPar1, scw.ZoneNlAms1} {
		err := refreshInstanceImage(api, database, zone)
		if err != nil {
			return nil, err
		}

		err = refreshInstanceIP(api, database, zone)
		if err != nil {
			return nil, err
		}

		err = refreshInstancePlacementGroup(api, database, zone)
		if err != nil {
			return nil, err
		}

		err = refreshInstanceSecurityGroup(api, database, zone)
		if err != nil {
			return nil, err
		}

		err = refreshInstanceServer(api, database, zone)
		if err != nil {
			return nil, err
		}

		err = refreshInstanceSnapshot(api, database, zone)
		if err != nil {
			return nil, err
		}

		err = refreshInstanceVolume(api, database, zone)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func refreshInstancePlacementGroup(api *instance.API, database *gorm.DB, zone scw.Zone) error {
	listPlacementGroups, err := api.ListPlacementGroups(&instance.ListPlacementGroupsRequest{
		Zone: zone,
	}, scw.WithAllPages())
	if err != nil {
		return err
	}
	for _, placementGroup := range listPlacementGroups.PlacementGroups {
		database.Create(&cachePlacementGroup{
			UUID:      placementGroup.ID,
			Name:      placementGroup.Name,
			Zone:      placementGroup.Zone.String(),
			ProjectID: placementGroup.Organization,
		})
	}
	return nil
}

func refreshInstanceSecurityGroup(api *instance.API, database *gorm.DB, zone scw.Zone) error {
	listSecurityGroups, err := api.ListSecurityGroups(&instance.ListSecurityGroupsRequest{
		Zone: zone,
	}, scw.WithAllPages())
	if err != nil {
		return err
	}
	for _, securityGroup := range listSecurityGroups.SecurityGroups {
		database.Create(&cacheSecurityGroup{
			UUID:      securityGroup.ID,
			Name:      securityGroup.Name,
			Zone:      securityGroup.Zone.String(),
			ProjectID: securityGroup.Organization,
		})
	}
	return nil
}

func refreshInstanceServer(api *instance.API, database *gorm.DB, zone scw.Zone) error {
	listServers, err := api.ListServers(&instance.ListServersRequest{
		Zone: zone,
	}, scw.WithAllPages())
	if err != nil {
		return err
	}
	for _, server := range listServers.Servers {
		database.Create(&cacheServer{
			UUID:      server.ID,
			Name:      server.Name,
			Zone:      server.Zone.String(),
			ProjectID: server.Organization,
		})
	}
	return nil
}

func refreshInstanceSnapshot(api *instance.API, database *gorm.DB, zone scw.Zone) error {
	listSnapshots, err := api.ListSnapshots(&instance.ListSnapshotsRequest{
		Zone: zone,
	})
	if err != nil {
		return err
	}
	for _, snapshots := range listSnapshots.Snapshots {
		database.Create(&cacheSnapshot{
			UUID:      snapshots.ID,
			Name:      snapshots.Name,
			Zone:      snapshots.Zone.String(),
			ProjectID: snapshots.Organization,
		})
	}
	return nil
}

func refreshInstanceVolume(api *instance.API, database *gorm.DB, zone scw.Zone) error {
	listVolumes, err := api.ListVolumes(&instance.ListVolumesRequest{
		Zone: zone,
	})
	if err != nil {
		return err
	}
	for _, volume := range listVolumes.Volumes {
		database.Create(&cacheVolume{
			UUID:      volume.ID,
			Name:      volume.Name,
			Zone:      volume.Zone.String(),
			ProjectID: volume.Organization,
		})
	}
	return nil
}

func refreshInstanceIP(api *instance.API, database *gorm.DB, zone scw.Zone) error {
	listIP, err := api.ListIPs(&instance.ListIPsRequest{
		Zone: zone,
	}, scw.WithAllPages())
	if err != nil {
		return err
	}
	for _, ip := range listIP.IPs {
		database.Create(&cacheIP{
			UUID:      ip.ID,
			Address:   ip.Address.String(),
			Zone:      ip.Zone.String(),
			ProjectID: ip.Organization,
		})
	}

	return nil
}

func refreshInstanceImage(api *instance.API, database *gorm.DB, zone scw.Zone) error {
	listImage, err := api.ListImages(&instance.ListImagesRequest{
		Zone: zone,
	}, scw.WithAllPages())
	if err != nil {
		return err
	}
	for _, image := range listImage.Images {
		database.Create(&cacheImage{
			UUID:      image.ID,
			Name:      image.Name,
			Zone:      image.Zone.String(),
			ProjectID: image.Organization,
		})
	}
	return nil
}
