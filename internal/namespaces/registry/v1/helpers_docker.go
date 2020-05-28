package registry

import (
	"encoding/json"
	"errors"
	"os"
	"path"
)

const (
	dockerConfigDir      = ".docker"
	dockerConfigFilename = "config.json"
	dockerCredHelpersKey = "credHelpers"
)

func writeHelperScript(scriptPath string, scriptContent string) error {
	scriptDirArg := path.Dir(scriptPath)
	if _, err := os.Stat(scriptDirArg); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if err := os.MkdirAll(scriptDirArg, 0755); err != nil {
			return err
		}
	}

	f, err := os.OpenFile(scriptPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(scriptContent))
	if err != nil {
		return err
	}

	return nil
}

type dockerConfigFile struct {
	CredHelpers map[string]string `json:"credHelpers,omitempty"`
}

func setupDockerConfigFile(registries []string, binaryName string) error {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dockerConfigFilePath := path.Join(userHomeDir, dockerConfigDir, dockerConfigFilename)
	if _, err := os.Stat(dockerConfigFilePath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if err := os.MkdirAll(path.Dir(dockerConfigFilePath), 0755); err != nil {
			return err
		}
	}

	f, err := os.OpenFile(dockerConfigFilePath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	var data map[string]interface{}

	err = json.NewDecoder(f).Decode(&data)
	if err != nil {
		return err
	}
	credHelpers := make(map[string]interface{})

	credHelpersInterface, ok := data[dockerCredHelpersKey]
	if ok {
		credHelpers = credHelpersInterface.(map[string]interface{})
	}

	for _, reg := range registries {
		credHelpers[reg] = binaryName
	}

	data[dockerCredHelpersKey] = credHelpers

	jsonData, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		return err
	}

	err = f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = f.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}
