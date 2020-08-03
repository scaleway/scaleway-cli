package rdb

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type cacheInstance struct {
	ID        string `gorm:"type:uuid;primary_key" cli:"id"`
	Name      string
	Region    string
	ProjectID string
}

func (cacheInstance) TableName() string {
	return "rdb_instance"
}

type cachedResource interface {
	EnsureExists(ctx context.Context, concreteResource interface{})
	EnsureAbsent(ctx context.Context, concreteResource interface{})
}

func (*cacheInstance) EnsureExists(ctx context.Context, resI interface{}) {
	database := core.ExtractCacheDB(ctx)
	instance := resI.(*rdb.Instance)

	database.AutoMigrate(&cacheInstance{})
	var row cacheInstance
	database.Where("id = ?", instance.ID, &row)
	// If this ID is not present we create it
	if row.ID == "" {
		database.Create(&cacheInstance{
			ID:        instance.ID,
			Name:      instance.Name,
			Region:    instance.Region.String(),
			ProjectID: instance.OrganizationID,
		})
	} else {
		// otherwise we update the field name
		database.Model(&row).Where("id = ?", instance.ID).Update("name", instance.Name)
	}
}

func (*cacheInstance) EnsureAbsent(ctx context.Context, resI interface{}) {
	database := core.ExtractCacheDB(ctx)
	instance := resI.(*rdb.Instance)
	database.AutoMigrate(&cacheInstance{})

	var row cacheInstance
	if result := database.Where("id = ?", instance.ID).First(&row); result.Error == nil {
		database.Unscoped().Where("id = ?", instance.ID).Delete(&row)
	}
}

type cacheDatabase struct {
	Name      string `gorm:"primary_key"`
	Instance  string `gorm:"type:uuid;primary_key"`
	Region    string
	ProjectID string
}

func (d cacheDatabase) EnsureExists(ctx context.Context, argsI interface{}, resI interface{}) {
	database := core.ExtractCacheDB(ctx)

	instance := resI.(*rdb.Database)

	database.AutoMigrate(&cacheInstance{})
	var row cacheInstance
	database.Where("id = ?", instance.ID, &row)
	// If this ID is not present we create it
	if row.ID == "" {
		database.Create(&cacheInstance{
			ID:        instance.ID,
			Name:      instance.Name,
			Region:    instance.Region.String(),
			ProjectID: instance.OrganizationID,
		})
	} else {
		// otherwise we update the field name
		database.Model(&row).Where("id = ?", instance.ID).Update("name", instance.Name)
	}
}

func (d cacheDatabase) EnsureAbsent(ctx context.Context, argsI interface{}, resI interface{}) {
	database := core.ExtractCacheDB(ctx)
	args := argsI.(*rdb.DeleteDatabaseRequest)
	rdb_database := resI.(*rdb.Database)
	database.AutoMigrate(&cacheInstance{})

	var row cacheInstance
	if result := database.Where("name = ?", rdb_database.Name).Where("instance = ?", args.InstanceID).First(&row); result.Error == nil {
		database.Unscoped().Where("name = ?", rdb_database.Name).Where("instance = ?", args.InstanceID).Delete(&row)
	}
}

func (cacheDatabase) TableName() string {
	return "rdb_database"
}

type cacheUser struct {
	ID        string `gorm:"type:uuid;primary_key"`
	Name      string
	Instance  string
	Region    string
	ProjectID string
}

func (cacheUser) TableName() string {
	return "rdb_user"
}

var registerCache = map[interface{}]cachedResource{
	reflect.TypeOf(rdb.Instance{}): &cacheInstance{},
}

func RefreshCache(ctx context.Context) (interface{}, error) {
	client := core.ExtractClient(ctx)
	api := rdb.NewAPI(client)
	database := core.ExtractCacheDB(ctx)

	cInstance := &cacheInstance{}
	database.AutoMigrate(cInstance)
	database.Unscoped().Delete(cInstance)

	database.AutoMigrate(&cacheDatabase{})
	database.Unscoped().Delete(&cacheDatabase{})

	database.AutoMigrate(&cacheUser{})
	database.Unscoped().Delete(&cacheUser{})

	for _, region := range []scw.Region{scw.RegionFrPar, scw.RegionNlAms} {
		listInstances, err := api.ListInstances(&rdb.ListInstancesRequest{
			Region: region,
		})
		if err != nil {
			return nil, err
		}

		for _, dbInstance := range listInstances.Instances {
			cInstance.EnsureExists(ctx, dbInstance)

			listDatabases, err := api.ListDatabases(&rdb.ListDatabasesRequest{
				Region:     dbInstance.Region,
				InstanceID: dbInstance.ID,
			})
			if err != nil {
				return nil, err
			}
			for _, db := range listDatabases.Databases {
				database.Create(&cacheDatabase{
					Name:      db.Name,
					Instance:  dbInstance.ID,
					Region:    dbInstance.Region.String(),
					ProjectID: dbInstance.OrganizationID,
				})
			}

			listUsers, err := api.ListUsers(&rdb.ListUsersRequest{
				Region:     region,
				InstanceID: dbInstance.ID,
			})
			if err != nil {
				return nil, err
			}
			for _, user := range listUsers.Users {
				database.Create(&cacheUser{
					Name:      user.Name,
					Instance:  dbInstance.ID,
					Region:    dbInstance.Region.String(),
					ProjectID: dbInstance.OrganizationID,
				})
			}
		}
	}

	return nil, nil
}
