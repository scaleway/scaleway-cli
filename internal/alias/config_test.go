package alias_test

import (
	"fmt"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/alias"
	"github.com/stretchr/testify/assert"
)

func TestConfig_ResolveAliases(t *testing.T) {
	tests := []struct {
		Aliases  map[string][]string
		Command  []string
		Expected []string
	}{
		{
			Aliases: map[string][]string{
				"isl": {"instance", "server", "list"},
			},
			Command:  []string{"isl"},
			Expected: []string{"instance", "server", "list"},
		},
		{
			Aliases: map[string][]string{
				"isl": {"instance", "server", "list"},
			},
			Command:  []string{"scw", "isl"},
			Expected: []string{"scw", "instance", "server", "list"},
		},
		{
			Aliases: map[string][]string{
				"sl": {"server", "list"},
			},
			Command:  []string{"instance", "sl", "zone=fr-par-1"},
			Expected: []string{"instance", "server", "list", "zone=fr-par-1"},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("Resolve_TestCase%d", i), func(t *testing.T) {
			config := &alias.Config{Aliases: tt.Aliases}
			actual := config.ResolveAliases(tt.Command)
			assert.Equal(t, tt.Expected, actual)
		})
	}
}
