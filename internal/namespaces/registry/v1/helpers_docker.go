package registry

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	dockerConfigDir      = ".docker"
	dockerConfigFilename = "config.json"
	dockerCredHelpersKey = "credHelpers"
)

func writeHelperScript(scriptPath string, scriptContent string) error {
	scriptDir := path.Dir(scriptPath)
	stats, err := os.Stat(scriptDir)
	if err != nil {
		return err
	}
	if !stats.IsDir() {
		return fmt.Errorf("%s is not a directory", scriptDir)
	}

	f, err := os.OpenFile(scriptPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(scriptContent)
	if err != nil {
		return err
	}

	return nil
}

func setupDockerConfigFile(ctx context.Context, registries []string, binaryName string) error {
	homeDir := core.ExtractUserHomeDir(ctx)

	dockerConfigFilePath := path.Join(homeDir, dockerConfigDir, dockerConfigFilename)
	if _, err := os.Stat(dockerConfigFilePath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if err := os.MkdirAll(path.Dir(dockerConfigFilePath), 0o755); err != nil {
			return err
		}
	}

	f, err := os.OpenFile(dockerConfigFilePath, os.O_RDWR|os.O_CREATE, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()

	dockerConfig := map[string]any{}

	dockerConfigRaw, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	if len(dockerConfigRaw) == 0 {
		dockerConfigRaw = []byte("{}")
	}

	err = json.Unmarshal(dockerConfigRaw, &dockerConfig)
	if err != nil {
		return err
	}

	credHelpers := map[string]any{}
	if ch, ok := dockerConfig[dockerCredHelpersKey]; ok {
		credHelpers = ch.(map[string]any)
	}

	for _, reg := range registries {
		credHelpers[reg] = binaryName
	}

	dockerConfig[dockerCredHelpersKey] = credHelpers

	err = f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(dockerConfig)
	if err != nil {
		return err
	}

	return nil
}

func getRegistryEndpoint(region scw.Region) string {
	return endpointPrefix + region.String() + endpointSuffix
}
