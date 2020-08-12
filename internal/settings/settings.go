package settings

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"gopkg.in/yaml.v2"
)

type Settings struct {
	Output *string `yaml:"output,omitempty" json:"output,omitempty"`
}

const (
	defaultPermission       = 0600
	defaultSettingsFileName = "cli.yaml"
	DefaultPrinter          = "human"

	settingsFileTemplate = `# Scaleway CLI settings file
# This settings file can be used only with Scaleway CLI (>2.0.0) (https://github.com/scaleway/scaleway-cli)
# Output sets the output format for all commands you run
{{ if .Output }}output: {{ .Output }}{{ else }}# output: human{{ end }}
`
)

// getSettingsFilePath returns the path to the settings file
func getSettingsFilePath() (string, bool) {
	configDir, err := scw.GetScwConfigDir()
	if err != nil {
		return "", false
	}
	return filepath.Clean(filepath.Join(configDir, defaultSettingsFileName)), true
}

func Default() *Settings {
	output := DefaultPrinter

	return &Settings{
		Output: &output,
	}
}

// clone deep copy settings object
func (s *Settings) clone() *Settings {
	s2 := &Settings{}
	settingsRaw, _ := yaml.Marshal(s)
	_ = yaml.Unmarshal(settingsRaw, s2)
	return s2
}

func (s *Settings) String() string {
	s2 := s.clone()
	settingsRaw, _ := yaml.Marshal(s2)
	return string(settingsRaw)
}

func (s *Settings) IsEmpty() bool {
	return s.String() == "{}\n"
}

func unmarshalSettings(content []byte) (*Settings, error) {
	var settings Settings

	err := yaml.Unmarshal(content, &settings)
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

// Load read the settings from the default path.
func Load() (*Settings, error) {
	path, exists := getSettingsFilePath()
	if !exists {
		return nil, fmt.Errorf("cli settings file does not exist")
	}
	return LoadSettingsFromPath(path)
}

// LoadSettingsFromPath read the settings from the given path.
func LoadSettingsFromPath(path string) (*Settings, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("cli settings not found")
	}
	if err != nil {
		return nil, err
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read cli settings file")
	}

	settings, err := unmarshalSettings(file)
	if err != nil {
		return nil, fmt.Errorf("content of cli settings file %s is invalid", path)
	}

	return settings, nil
}

// SaveTo will save the settings to the default settings path. This
// action will overwrite the previous file when it exists.
func (s *Settings) Save() error {
	path, _ := getSettingsFilePath()
	return s.SaveTo(path)
}

// HumanSettings will generate a CLI settings file with documented arguments.
func (s *Settings) HumanSettings() (string, error) {
	tmpl, err := template.New("configuration").Parse(settingsFileTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, s)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// SaveTo will save the settings to the given path. This action will
// overwrite the previous file when it exists.
func (s *Settings) SaveTo(path string) error {
	path = filepath.Clean(path)

	// STEP 1: Render the configuration file as a file
	file, err := s.HumanSettings()
	if err != nil {
		return err
	}

	// STEP 2: create config path dir in cases it didn't exist before
	err = os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		return err
	}

	// STEP 3: write new settings file
	err = ioutil.WriteFile(path, []byte(file), defaultPermission)
	if err != nil {
		return err
	}

	return nil
}
