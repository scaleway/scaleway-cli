package account

import (
	"context"

	"github.com/scaleway/scaleway-cli/internal/core"
	account "github.com/scaleway/scaleway-sdk-go/api/account/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type cacheSSHKey struct {
	ID        string
	Name      string
	ProjectID string
}

func (cacheSSHKey) TableName() string {
	return "account_ssh_keys"
}

func RefreshCache(ctx context.Context) (interface{}, error) {
	client := core.ExtractClient(ctx)
	api := account.NewAPI(client)
	log := core.ExtractLogger(ctx)
	database := core.ExtractCacheDB(ctx)

	database.AutoMigrate(&cacheSSHKey{})

	// Fetching
	listSSHKeys, err := api.ListSSHKeys(&account.ListSSHKeysRequest{}, scw.WithAllPages())
	if err != nil {
		return nil, err
	}

	// Persisting
	database.AutoMigrate(&cacheSSHKey{})
	database.Unscoped().Delete(&cacheSSHKey{})
	for _, key := range listSSHKeys.SSHKeys {
		database.Create(&cacheSSHKey{
			ID:        key.ID,
			Name:      key.Name,
			ProjectID: key.OrganizationID,
		})
	}
	log.Info("Successfully build cache for account ssh keys")

	return nil, nil
}
