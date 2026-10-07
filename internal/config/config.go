package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/ghodss/yaml"
	"github.com/scaleway/scaleway-cli/v2/internal/alias"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	ScwConfigPathEnv = "SCW_CLI_CONFIG_PATH"

	DefaultConfigFileName   = "cli.yaml"
	defaultConfigPermission = 0o644

	DefaultOutput      = "human"
	configFileTemplate = `# Scaleway CLI config file
# This config file can be used only with Scaleway CLI (>2.0.0) (https://github.com/scaleway/scaleway-cli)
# Output sets the output format for all commands you run
{{ if .Output }}output: {{ .Output }}{{ else }}# output: human{{ end }}

# TUI (scw tui) configuration: keybindings overrides and : commands
{{- if .Tui }}
tui:
    {{- if .Tui.KeyBindings }}
    keybindings:
        {{- range $action, $key := .Tui.KeyBindings }}
        {{ $action }}: {{ $key }}
        {{- end }}
    {{- end }}
    {{- if .Tui.Commands }}
    commands:
        {{- range $name, $target := .Tui.Commands }}
        {{ $name }}: {{ $target }}
        {{- end }}
    {{- end }}
{{- else }}
# tui:
#     keybindings:
#         quit: q
#         menu: m
#         refresh: r
#     commands:
#         mydb: rdb instance
{{- end }}

# Whether autocomplete was proposed during init
{{- if .HasInstallAutocomplete }}
install_autocomplete: {{ .InstallAutocomplete }}
{{- else }}
# install_autocomplete: true
{{- end }}

# Alias creates custom aliases for your Scaleway CLI commands
{{- if .Alias }}
alias:
    aliases:
        {{- range $alias, $commands := .Alias.Aliases }}
        {{ $alias }}:
        {{- range $index, $command := $commands }}
            - {{ $command }}
        {{- end }}
        {{- end }}
{{- else }}
# alias:
#     aliases:
#         isl:
#             - instance
#             - server
#             - list
{{- end }}
`
)

type Config struct {
	Alias               *alias.Config `json:"alias"                          yaml:"alias"`
	Output              string        `json:"output"                         yaml:"output"`
	InstallAutocomplete *bool         `json:"install_autocomplete,omitempty" yaml:"install_autocomplete,omitempty"`
	Tui                 *TuiConfig    `json:"tui,omitempty"                  yaml:"tui"`

	path string
}

// TuiConfig configures the interactive TUI (scw tui).
type TuiConfig struct {
	// KeyBindings overrides the default TUI keybindings: one key per
	// action (e.g. refresh: R). The recognized actions are quit, menu,
	// command, filter, refresh, back, open, copy, next, prev, up,
	// down, page-up, page-down, home, end, zone (ctrl+z) and region
	// (ctrl+r); unset actions keep their defaults.
	KeyBindings map[string]string `json:"keybindings,omitempty" yaml:"keybindings"`
	// Commands adds custom : commands: a word mapped to a resource
	// name (e.g. mydb: rdb instance) or a built-in command (quit, help).
	Commands map[string]string `json:"commands,omitempty" yaml:"commands"`
}

// LoadConfig tries to load config file
// returns a new empty config if file doesn't exist
// return error if fail to load config file
func LoadConfig(configPath string) (*Config, error) {
	if runtime.GOARCH == "wasm" {
		return &Config{
			Alias:  alias.EmptyConfig(),
			Output: DefaultOutput,
			path:   configPath,
		}, nil
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{
				Alias:  alias.EmptyConfig(),
				Output: DefaultOutput,
				path:   configPath,
			}, nil
		}

		return nil, fmt.Errorf("failed to read cli config file: %w", err)
	}
	config := &Config{
		Alias:  alias.EmptyConfig(),
		Output: DefaultOutput,
		path:   configPath,
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal cli config: %w", err)
	}

	return config, nil
}

// Save marshal config to config file
func (c *Config) Save() error {
	if runtime.GOARCH == "wasm" {
		return nil
	}

	file, err := c.HumanConfig()
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(c.path), 0o700)
	if err != nil {
		return err
	}

	return os.WriteFile(c.path, []byte(file), defaultConfigPermission)
}

// HasInstallAutocomplete returns true if the install_autocomplete preference has been set
func (c *Config) HasInstallAutocomplete() bool {
	return c.InstallAutocomplete != nil
}

// HumanConfig will generate a config file with documented arguments
func (c *Config) HumanConfig() (string, error) {
	tmpl, err := template.New("configuration").Parse(configFileTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, c)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
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
