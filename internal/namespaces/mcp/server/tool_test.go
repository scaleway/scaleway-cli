package server

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
)

func TestCommandNameToToolName(t *testing.T) {
	tests := []struct {
		name     string
		command  *core.Command
		expected string
	}{
		{
			name: "namespace only",
			command: &core.Command{
				Namespace: "config",
			},
			expected: "config",
		},
		{
			name: "namespace and resource",
			command: &core.Command{
				Namespace: "instance",
				Resource:  "server",
			},
			expected: "instance_server",
		},
		{
			name: "full command",
			command: &core.Command{
				Namespace: "instance",
				Resource:  "server",
				Verb:      "list",
			},
			expected: "instance_server_list",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CommandNameToToolName(tt.command)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
