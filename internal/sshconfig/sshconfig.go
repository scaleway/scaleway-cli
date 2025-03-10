package sshconfig

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	// sshConfigFileName is the name of the file generated by this package
	sshConfigFileName = "scaleway.config"
	// sshDefaultConfigFileName is the name of the default ssh config
	sshDefaultConfigFileName = "config"
	sshConfigFolderHomePath  = ".ssh"

	sshConfigFileMode   = os.FileMode(0o600)
	sshConfigFolderMode = os.FileMode(0o700)

	ErrFileNotFound = errors.New("file not found")
)

type Host interface {
	Config() string
}

func Generate(hosts []Host) ([]byte, error) {
	configBuffer := bytes.NewBuffer(nil)

	for _, host := range hosts {
		configBuffer.WriteString(host.Config())
		configBuffer.WriteString("\n")
	}

	return configBuffer.Bytes(), nil
}

// ConfigFilePath returns the path of the generated file
// should be ~/.ssh/scaleway.config
func ConfigFilePath(homeDir string) string {
	configFolder := sshConfigFolder(homeDir)
	configFile := filepath.Join(configFolder, sshConfigFileName)

	return configFile
}

func Save(homeDir string, hosts []Host) error {
	cfg, err := Generate(hosts)
	if err != nil {
		return err
	}

	configFile := ConfigFilePath(homeDir)

	err = os.WriteFile(configFile, cfg, sshConfigFileMode)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(sshConfigFolder(homeDir), sshConfigFolderMode)
			if err != nil {
				return fmt.Errorf("failed to create ssh config folder: %w", err)
			}

			return os.WriteFile(configFile, cfg, sshConfigFileMode)
		}

		return err
	}

	return nil
}

func sshConfigFolder(homeDir string) string {
	return filepath.Join(homeDir, sshConfigFolderHomePath)
}

func includeLine() string {
	return "Include " + sshConfigFileName
}

// DefaultConfigFilePath returns the default ssh config file path
// should be ~/.ssh/config
func DefaultConfigFilePath(homeDir string) string {
	configFolder := sshConfigFolder(homeDir)
	configFilePath := filepath.Join(configFolder, sshDefaultConfigFileName)

	return configFilePath
}

func openDefaultConfigFile(homeDir string) (*os.File, error) {
	configFilePath := DefaultConfigFilePath(homeDir)

	configFile, err := os.Open(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotFound
		}

		return nil, fmt.Errorf("failed to open default ssh config file: %w", err)
	}

	return configFile, nil
}

// ConfigIsIncluded checks that ssh config file is included in user's .ssh/config
// Default config file ~/.ssh/config should start with "Include scaleway.config"
func ConfigIsIncluded(homeDir string) (bool, error) {
	configFile, err := openDefaultConfigFile(homeDir)
	if err != nil {
		return false, err
	}
	defer configFile.Close()

	expectedLine := includeLine()

	fileScanner := bufio.NewScanner(configFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		if strings.Contains(fileScanner.Text(), expectedLine) {
			return true, nil
		}
	}

	return false, nil
}

// IncludeConfigFile edit default ssh config to include this package generated file
// ~/.ssh/config will be prepended with "Include scaleway.config"
func IncludeConfigFile(homeDir string) error {
	configFileMode := sshConfigFileMode
	fileContent := []byte(nil)

	configFile, err := openDefaultConfigFile(homeDir)
	if err != nil && !errors.Is(err, ErrFileNotFound) {
		return err
	}

	if configFile != nil {
		// Keep file mode and permissions if it exists
		fi, err := configFile.Stat()
		if err != nil {
			_ = configFile.Close()

			return fmt.Errorf("failed to stat file: %w", err)
		}
		configFileMode = fi.Mode()

		fileContent, err = io.ReadAll(configFile)
		if err != nil {
			_ = configFile.Close()

			return fmt.Errorf("failed to read file: %w", err)
		}

		_ = configFile.Close()
	}

	// Prepend config file with Include line
	fileContent = append([]byte(includeLine()+"\n"), fileContent...)

	configFolder := sshConfigFolder(homeDir)
	configFilePath := filepath.Join(configFolder, sshDefaultConfigFileName)

	err = os.WriteFile(configFilePath, fileContent, configFileMode)
	if err != nil {
		return fmt.Errorf("failed to write config file %s: %w", configFilePath, err)
	}

	return nil
}
