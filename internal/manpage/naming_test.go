package manpage_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/manpage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_PageName(t *testing.T) {
	tests := []struct {
		name        string
		commandPath []string
		expected    string
	}{
		{
			name:        "top level",
			commandPath: []string{},
			expected:    "scw",
		},
		{
			name:        "namespace only",
			commandPath: []string{"instance"},
			expected:    "scw-instance",
		},
		{
			name:        "full command path",
			commandPath: []string{"instance", "server", "create"},
			expected:    "scw-instance-server-create",
		},
		{
			name:        "dash containing namespace token",
			commandPath: []string{"edge-services", "origin", "list"},
			expected:    "scw-edge-services-origin-list",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, manpage.PageName(test.commandPath...))
			assert.Equal(
				t,
				test.expected+".1",
				manpage.FileName(manpage.PageName(test.commandPath...)),
			)
		})
	}
}

func Test_DetectNameCollisions(t *testing.T) {
	// Two distinct command paths that map to the same page name because a
	// token contains a dash: `edge-services origin list` and
	// `edge services-origin list` (hypothetical).
	conflicting := core.NewCommands(
		&core.Command{Namespace: "edge-services", Resource: "origin", Verb: "list", Short: "A"},
		&core.Command{Namespace: "edge", Resource: "services-origin", Verb: "list", Short: "B"},
	)

	_, err := manpage.DetectNameCollisions(conflicting)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "collision")
	assert.Contains(t, err.Error(), "scw-edge-services-origin-list.1")

	// Non conflicting set.
	ok := core.NewCommands(
		&core.Command{Namespace: "instance", Resource: "server", Verb: "list", Short: "A"},
		&core.Command{Namespace: "instance", Resource: "server", Verb: "create", Short: "B"},
	)
	_, err = manpage.DetectNameCollisions(ok)
	require.NoError(t, err)
}
