package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	iam "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
)

var (
	apiKeyExpireTime        = 24 * time.Hour
	lastChecksFileLocalName = "last-cli-checks"
)

type AfterCommandCheckFunc func(ctx context.Context)

// wasFileModifiedLast24h checks whether the file has been updated during last 24 hours.
func wasFileModifiedLast24h(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}

	yesterday := time.Now().AddDate(0, 0, -1)
	lastUpdate := stat.ModTime()
	return lastUpdate.After(yesterday)
}

func GetLatestVersionUpdateFilePath(cacheDir string) string {
	return filepath.Join(cacheDir, lastChecksFileLocalName)
}

// CreateAndCloseFile creates a file and closes it. It returns true on succeed, false on failure.
func CreateAndCloseFile(path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0o700)
	if err != nil {
		return fmt.Errorf("failed creating path %s: %s", path, err)
	}
	newFile, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o600)
	if err != nil {
		return fmt.Errorf("failed creating file %s: %s", path, err)
	}

	return newFile.Close()
}

// runAfterCommandChecks execute checks after a command has been executed
// skipped if command has disabled checks or was executed in the last 24 hours
func runAfterCommandChecks(ctx context.Context, checkFuncs ...AfterCommandCheckFunc) {
	cmd := extractMeta(ctx).command
	cmdDisableCheck := cmd != nil && cmd.DisableAfterChecks
	if cmdDisableCheck {
		ExtractLogger(ctx).Debug("skipping after command checks")
		return
	}

	lastChecksFilePath := GetLatestVersionUpdateFilePath(ExtractCacheDir(ctx))

	// do nothing if last refresh at during the last 24h
	if wasFileModifiedLast24h(lastChecksFilePath) {
		ExtractLogger(ctx).Debug("version was already checked during past 24 hours")
		return
	}

	// do nothing if we cannot create the file
	err := CreateAndCloseFile(lastChecksFilePath)
	if err != nil {
		ExtractLogger(ctx).Debug(err.Error())
		return
	}

	for _, checkFunc := range checkFuncs {
		checkFunc(ctx)
	}
}

// Check if API Key is about to expire
func checkAPIKey(ctx context.Context) {
	client := ExtractClient(ctx)
	if client == nil {
		return
	}
	accessKey, exists := client.GetAccessKey()
	if !exists {
		return
	}

	api := iam.NewAPI(client)
	apiKey, err := api.GetAPIKey(&iam.GetAPIKeyRequest{
		AccessKey: accessKey,
	})
	if err != nil || apiKey.ExpiresAt == nil {
		return
	}
	now := time.Now()
	if apiKey.ExpiresAt.Before(now.Add(apiKeyExpireTime)) {
		expiresIn := apiKey.ExpiresAt.Sub(now).Truncate(time.Second).String()
		ExtractLogger(ctx).Warningf("Current api key expires in %s\n", expiresIn)
	}
}
