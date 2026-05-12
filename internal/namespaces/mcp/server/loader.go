package server

import (
	"slices"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
)

// SplitArg splits a comma-separated string into a slice of strings.
// Returns an empty slice if the input is empty.
func SplitArg(s string) []string {
	if s == "" {
		return []string{}
	}

	return strings.Split(s, ",")
}

var (
	// ExcludedNamespaces is used to filter out core.Command that should not be exposed as MCP tools based on their namespace.
	ExcludedNamespaces = []string{
		// Skip config namespace to avoid security risks
		"config",
		// Skip shell-centric namespaces
		"alias",
		"autocomplete",
		"shell",
		"login",
		"init",
	}

	// ExcludedVerbs is used to filter out core.Command that should not be exposed as MCP tools based on their verb.
	ExcludedVerbs = []string{
		"edit", // Shell-centric verb
	}
)

// CommandFilterConfig holds the configuration for filtering CLI commands
// when registering them as MCP tools/resources.
type CommandFilterConfig struct {
	ReadOnly          bool
	EnabledNamespaces []string
	EnabledResources  []string
	EnabledVerbs      []string
}

// FilterCommands filters CLI commands based on the given config.
// Returns a slice of CommandTool ready to be passed to NewMCPServer.
func FilterCommands(commands []*core.Command, config CommandFilterConfig) []*CommandTool {
	result := make([]*CommandTool, 0, len(commands))
	for _, cmd := range commands {
		if ShouldLoadCommand(cmd, config) {
			result = append(result, NewCommandTool(cmd))
		}
	}

	return result
}

// ShouldLoadCommand returns true if the command should be registered as an MCP tool.
// It filters out:
// - Hidden commands
// - Commands with ExcludeFromMCP flag set
// - Commands without a Run function (namespace/resource containers)
// - Commands in excluded namespaces
// - Commands with excluded verbs
// - When readOnly is true, only commands with get/list verbs are registered
// - When enabledNamespaces/ Resources/ Verbs are set, only matching commands are registered
func ShouldLoadCommand(cmd *core.Command, config CommandFilterConfig) bool {
	// Skip hidden commands
	if cmd.Hidden {
		return false
	}

	if cmd.ExcludeFromMCP {
		return false
	}

	// Skip commands without a Run function (namespace/resource containers)
	if cmd.Run == nil {
		return false
	}

	if slices.Contains(ExcludedNamespaces, cmd.Namespace) {
		return false
	}

	if slices.Contains(ExcludedVerbs, cmd.Verb) {
		return false
	}

	// If enabled namespaces are specified, only allow those namespaces
	if len(config.EnabledNamespaces) > 0 &&
		!slices.Contains(config.EnabledNamespaces, cmd.Namespace) {
		return false
	}

	// If enabled resources are specified, only allow those resources
	if len(config.EnabledResources) > 0 && !slices.Contains(config.EnabledResources, cmd.Resource) {
		return false
	}

	// If enabled verbs are specified, only allow those verbs
	if len(config.EnabledVerbs) > 0 && !slices.Contains(config.EnabledVerbs, cmd.Verb) {
		return false
	}

	// In read-only mode, only allow get/list operations
	if config.ReadOnly && !cmd.IsReadOnly() {
		return false
	}

	return true
}
