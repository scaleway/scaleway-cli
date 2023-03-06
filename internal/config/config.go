package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/scaleway/scaleway-cli/v2/internal/alias"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	ScwConfigPathEnv = "SCW_CLI_CONFIG_PATH"

	DefaultConfigFileName   = "cli.yaml"
	defaultConfigPermission = 0644
)

type Config struct {
	Alias *alias.Config `json:"alias"`

	path string
}

// LoadConfig tries to load config file
// returns a new empty config if file doesn't exist
// return error if fail to load config file
func LoadConfig(configPath string) (*Config, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{
				Alias: alias.EmptyConfig(),
				path:  configPath,
			}, nil
		}
		return nil, fmt.Errorf("failed to read cli config file: %w", err)
	}
	config := &Config{
		Alias: alias.EmptyConfig(),
		path:  configPath,
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal cli config: %w", err)
	}

	return config, nil
}

// Save marshal config to config file
func (c *Config) Save() error {
	config, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(c.path), 0700)
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, config, defaultConfigPermission)
}

func FilePath() (string, error) {
	file := os.Getenv(ScwConfigPathEnv)
	if file != "" {
		return file, nil
	}
	configDir, err := scw.GetScwConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Clean(filepath.Join(configDir, DefaultConfigFileName)), nil
}
